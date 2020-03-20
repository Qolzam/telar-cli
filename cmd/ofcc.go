package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type ClientAction struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type DialogInfoPayload struct {
	Message string `json:"message"`
	URL     string `json:"url"`
}

func echoDialogInfo(message string, url string) {
	action := ClientAction{
		Type: SHOW_INFO_DIALOG,
		Payload: DialogInfoPayload{
			Message: message,
			URL:     url,
		},
	}
	Echo(action)
}

func echoClient(actionType string, payload interface{}) {
	action := ClientAction{
		Type:    actionType,
		Payload: payload,
	}
	Echo(action)
}

func echoInput(key string, value interface{}) {
	echoClient(SET_INPUT,
		struct {
			Key   string      `json:"key"`
			Value interface{} `json:"value"`
		}{
			Key:   key,
			Value: value,
		},
	)
}

func echoStep(step int) {
	echoClient(SET_SETUP_STEP,
		struct {
			Step int `json:"step"`
		}{
			Step: step,
		},
	)
}

func echoOpenDeploy(open bool) {
	echoClient(SET_DEPLOY_OPEN,
		struct {
			Open bool `json:"open"`
		}{
			Open: open,
		},
	)
}

func echoSetupState(state string) {
	echoClient(SET_SETUP_STATE,
		struct {
			State string `json:"state"`
		}{
			State: state,
		},
	)
}

func downloadOFCCPublicKey(downloadTo string) error {
	downloadURL := "https://raw.githubusercontent.com/openfaas/cloud-functions/master/pub-cert.pem"

	fmt.Printf("Starting download of pub-cert.pe %s, this could take a few moments.\n", downloadURL)
	_, err := downloadFile(http.DefaultClient, downloadURL, "pub-cert.pem", downloadTo)

	if err != nil {
		return err
	}

	return nil
}

func checkCustomers(customer string) error {
	url := "https://raw.githubusercontent.com/openfaas/openfaas-cloud/master/CUSTOMERS"
	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)

	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error status code %d", res.StatusCode)
	}

	if res.Body != nil {
		defer res.Body.Close()
		res, _ := ioutil.ReadAll(res.Body)
		customers := string(res)
		if strings.Contains(customers, customer) {
			return nil
		} else {
			return fmt.Errorf("The %s user is not registered in OpenFaaS Cloud CUSTOMERS https://raw.githubusercontent.com/openfaas/openfaas-cloud/master/CUSTOMERS", customers)
		}
	}
	return fmt.Errorf("Unkown error when checking %s for OpenFaaS Cloud CUSTOMERS ", customer)
}

func getOFCCGateway(githubUsername string) string {
	url := fmt.Sprintf("https://%s.o6s.io", githubUsername)
	return url
}

func StartStep() {

	projectPath, err := getDefaultProjectDirectory()
	if err != nil {
		echoDialogInfo(err.Error(), "https://github.com/Qolzam/telar-cli/blob/master/docs/ofcc-setup/1.md")
		return
	}
	echoInput("projectDirectory", projectPath)

}

func CheckInitStep(projectPath string) {
	_ = checkDirectory(projectPath)
	echoStep(1)

}

func CheckIngredient(projectPath string, githubUsername string) {
	helpURL := "https://github.com/Qolzam/telar-cli/blob/master/docs/ofcc-setup/1.md"
	echoInput("loadingCheckIngredients", true)

	// Check kubeseal
	err := checkKubeseal()
	if isError(err) {
		echoDialogInfo(err.Error(), helpURL)
		echoInput("loadingCheckIngredients", false)
		fmt.Println(err.Error())
		return
	}
	echoInput("installKubeseal", true)

	// Check github username is registered in OpenFaaS Cloud CUSTOMERS
	err = checkCustomers(githubUsername)
	if isError(err) {
		echoDialogInfo(fmt.Sprintf("Github user name [%s] is not registered in OpenFaaS Cloud Community Cluster. Please check if you have typo.", githubUsername), helpURL)
		echoInput("loadingCheckIngredients", false)
		fmt.Println(err.Error())
		return
	}
	echoInput("githubUsernameRegisterd", true)

	// Check telar-web repository
	err = cloneTelarWeb(projectPath, githubUsername)
	if isError(err) && err.Error() != "repository already exists" {
		errMessage := "telar-web " + err.Error()

		echoDialogInfo(errMessage, helpURL)
		echoInput("loadingCheckIngredients", false)

		fmt.Println(errMessage)
		return
	}
	echoInput("cloneTelarWeb", true)

	// Check ts-serverless repository
	err = cloneTSServerless(projectPath, githubUsername)
	if isError(err) && err.Error() != "repository already exists" {
		errMessage := "ts-serverless " + err.Error()

		echoDialogInfo(errMessage, helpURL)
		echoInput("loadingCheckIngredients", false)

		fmt.Println(errMessage)
		return
	}
	echoInput("cloneTsServerless", true)

	// Check ts-ui repository
	err = cloneTSUI(projectPath, githubUsername)
	if isError(err) && err.Error() != "repository already exists" {
		errMessage := "ts-ui " + err.Error()

		echoDialogInfo(errMessage, helpURL)
		echoInput("loadingCheckIngredients", false)

		fmt.Println(errMessage)
		return
	}
	echoInput("cloneTsUi", true)

	echoStep(2)
}

func CheckStorage(projectPath string, bucketName string) {
	helpURL := "https://github.com/Qolzam/telar-cli/blob/master/docs/ofcc-setup/3.md"
	echoInput("loadingFirebaseStorage", true)

	// Check serviceAccount.json file for Firebase
	err := checkFirebaseServiceAccount(projectPath)
	if isError(err) {
		echoDialogInfo(err.Error(), helpURL)
		echoInput("loadingFirebaseStorage", false)
		fmt.Println(err.Error())
		return
	}
	echoInput("firebaseServiceAccount", true)

	// Check Firebase storage access
	err = checkFirebaseStorageBucket(projectPath, bucketName)
	if isError(err) {
		echoDialogInfo(err.Error(), helpURL)
		echoInput("loadingFirebaseStorage", false)
		fmt.Println(err.Error())
		return
	}
	echoInput("firebaseStorage", true)

	echoStep(3)

}

func CheckDatabase(mongoDBHost, MongoDBPassword string) {
	helpURL := "https://github.com/Qolzam/telar-cli/blob/master/docs/ofcc-setup/4.md"

	echoInput("loadingMongoDB", true)

	err := checkDB(mongoDBHost, MongoDBPassword)
	if isError(err) {
		echoDialogInfo(err.Error(), helpURL)
		echoInput("loadingMongoDB", false)
		fmt.Println(err.Error())
		return
	}
	echoInput("mongoDBConnection", true)

	echoStep(4)

}

func CheckRecaptcha() {
	echoStep(5)
}

func CheckOAuth() {
	echoStep(6)
}

func CheckUserManagement(githubUsername string) {
	helpURL := "https://github.com/Qolzam/telar-cli/blob/master/docs/ofcc-setup/7.md"

	payloadSecret, err := generatePayloadSecret()
	if isError(err) {
		echoDialogInfo(err.Error(), helpURL)
		fmt.Println(err.Error())
		return
	}
	echoInput("payloadSecret", payloadSecret)
	echoInput("gateway", fmt.Sprintf("https://%s.o6s.io", githubUsername))
	echoStep(7)
}

func CheckWebsocket(clientInput ClientInputs) {
	helpURL := "https://github.com/Qolzam/telar-cli/blob/master/docs/ofcc-setup/8.md"
	echoInput("loadingWebsocket", true)

	err := pingWebsocket(clientInput.WebsocketURL)
	if isError(err) {
		echoDialogInfo(err.Error(), helpURL)
		echoInput("loadingWebsocket", false)
		fmt.Println(err.Error())
		return
	}
	echoInput("websocketConnection", true)
	starteploy(clientInput)
}

func starteploy(clientInput ClientInputs) {
	echoOpenDeploy(true)

	// Apply stack
	telarConfig := TelarConfig{
		GithubUsername:   clientInput.GithubUsername,
		PathWD:           clientInput.ProjectDirectory,
		CoockieDomain:    ".o6s.io",
		Bucket:           clientInput.BucketName,
		ClientID:         clientInput.GithubOAuthClientID,
		URL:              fmt.Sprintf("https://%s.o6s.io", clientInput.GithubUsername),
		WebsocketURL:     clientInput.WebsocketURL,
		MongoDBHost:      clientInput.MongoDBHost,
		MongoDatabase:    clientInput.MongoDBHost,
		RecaptchaSiteKey: clientInput.SiteKeyRecaptcha,
		RefEmail:         clientInput.Gmail,
	}

	err := applyConfig(telarConfig)
	if isError(err) {
		echoDialogInfo(err.Error(), "")
		echoOpenDeploy(false)
		fmt.Println(err.Error())
		return
	}
	echoInput("loadingStackYaml", true)

	err = downloadOFCCPublicKey(telarConfig.PathWD)
	if isError(err) {
		echoDialogInfo(err.Error(), "")
		echoOpenDeploy(false)
		fmt.Println(err.Error())
		return
	}

	err = preparePublicPrivateKey(clientInput.ProjectDirectory)
	if isError(err) {
		echoDialogInfo(err.Error(), "")
		echoOpenDeploy(false)
		fmt.Println(err.Error())
		return
	}
	echoInput("loadingPublicPrivateKey", true)

	// Create secrets
	telarSecret := &TelarSecrets{
		MongoPwd:       clientInput.MongoDBPassword,
		RecaptchaKey:   clientInput.RecaptchaKey,
		TsClientSecret: clientInput.GithubOAuthSecret,
		RedisPwd:       "",
		AdminPwd:       clientInput.AdminPassword,
		AdminUsername:  clientInput.AdminUsername,
		PayloadSecret:  clientInput.PayloadSecret,
		RefEmailPwd:    clientInput.GmailPassword,
		PhoneAuthId:    "nil",
		PhoneAuthToken: "nil",
	}
	err = prepareSecret(clientInput.ProjectDirectory, clientInput.GithubUsername, telarSecret)
	if isError(err) {
		echoDialogInfo(err.Error(), "")
		echoOpenDeploy(false)
		fmt.Println(err.Error())
		return
	}
	echoInput("loadingCreateSecret", true)

	// Deploy Telar Web
	err = gitDeploy(clientInput.ProjectDirectory + "/telar-web")
	if isError(err) {
		echoDialogInfo(err.Error(), "")
		echoOpenDeploy(false)
		fmt.Println(err.Error())
		return
	}
	echoInput("deployTelarWeb", true)

	// Deploy Telar Social Serverless
	err = gitDeploy(clientInput.ProjectDirectory + "/ts-serverless")
	if isError(err) {
		echoDialogInfo(err.Error(), "")
		echoOpenDeploy(false)
		fmt.Println(err.Error())
		return
	}
	echoInput("deployTsServerless", true)

	// Deploy Telar Social UI
	err = gitDeploy(clientInput.ProjectDirectory + "/ts-ui")
	if isError(err) {
		echoDialogInfo(err.Error(), "")
		echoOpenDeploy(false)
		fmt.Println(err.Error())
		return
	}
	echoInput("deploySocialUi", true)
	echoSetupState("done")

}
