package cmd

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"

	"gopkg.in/yaml.v2"
)

const SETUP_YAML_FILE_NAME = "setup.yml"

type OFCSetupCache struct {
	ClientInputs ClientInputs `json:"clientInputs" yaml:"clientInputs,omitempty"`
}
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

func echoSetupDefaultValues(payload *OFCSetupCache) {
	echoClient(SET_SETUP_DEFAULT_VALUES, *payload)
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
			return fmt.Errorf("The %s user is not registered in OpenFaaS CUSTOMERS https://raw.githubusercontent.com/openfaas/openfaas-cloud/master/CUSTOMERS", customers)
		}
	}
	return fmt.Errorf("Unkown error when checking %s for OpenFaaS CUSTOMERS ", customer)
}

func getOFCCGateway(ofGateway string) string {
	url := fmt.Sprintf("https://%s", ofGateway)
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

	setupCache, err := getSetupYaml(projectPath)
	if err != nil {
		fmt.Printf("\n[WARN] can not get setup cache. %s", err.Error())
	} else {
		echoSetupDefaultValues(setupCache)
	}
	echoStep(1)

}

func OFCAccessSetting() {

	echoStep(2)
}
func CheckIngredient(projectPath string, githubUsername string) {
	helpURL := "https://github.com/Qolzam/telar-cli/blob/master/docs/ofcc-setup/2.md"
	echoInput("loadingCheckIngredients", true)

	// Check telar-web repository
	err := cloneTelarWeb(projectPath, githubUsername)
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

	echoStep(3)
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

	echoStep(4)

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

	echoStep(5)

}

func CheckRecaptcha() {
	echoStep(6)
}

func CheckOAuth() {
	echoStep(7)
}

func CheckUserManagement(ofGateway string) {
	helpURL := "https://github.com/Qolzam/telar-cli/blob/master/docs/ofcc-setup/7.md"

	payloadSecret, err := generatePayloadSecret()
	if isError(err) {
		echoDialogInfo(err.Error(), helpURL)
		fmt.Println(err.Error())
		return
	}
	echoInput("payloadSecret", payloadSecret)
	echoInput("gateway", ofGateway)
	echoStep(8)
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

	parsedURL, err := url.Parse(clientInput.SocialDomain)
	host, _, _ := net.SplitHostPort(parsedURL.Host)
	fmt.Println("[INFO] Host: ", host)

	// Apply stack
	telarConfig := TelarConfig{
		AppID:            clientInput.AppID,
		SecretName:       clientInput.SecretName,
		GithubUsername:   TELAR_GITHUB_USER_NAME,
		PathWD:           clientInput.ProjectDirectory,
		CoockieDomain:    "." + host,
		Bucket:           clientInput.BucketName,
		ClientID:         clientInput.GithubOAuthClientID,
		URL:              clientInput.SocialDomain,
		WebsocketURL:     clientInput.WebsocketURL,
		MongoDBHost:      clientInput.MongoDBHost,
		MongoDatabase:    clientInput.MongoDBName,
		RecaptchaSiteKey: clientInput.SiteKeyRecaptcha,
		RefEmail:         clientInput.Gmail,
	}

	err = applyConfig(telarConfig)
	if isError(err) {
		echoDialogInfo(err.Error(), "")
		echoOpenDeploy(false)
		fmt.Println(err.Error())
		return
	}
	echoInput("loadingStackYaml", true)

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

	var kubeconfigPath *string
	if clientInput.KubeconfigPath != "" {
		kubeconfigPath = &clientInput.KubeconfigPath
	}

	err = prepareSecret(clientInput.ProjectDirectory, clientInput.SecretName, clientInput.Namespace, kubeconfigPath, telarSecret)
	if isError(err) {
		echoDialogInfo(err.Error(), "")
		echoOpenDeploy(false)
		fmt.Println(err.Error())
		return
	}
	echoInput("loadingCreateSecret", true)

	// Deploy Telar Web
	openfaasPass, err := getOpenFaasPass()
	if isError(err) {
		echoDialogInfo(err.Error(), "")
		echoOpenDeploy(false)
		fmt.Println(err.Error())
		return
	}

	authConf, _, err := runFaaSLogin(clientInput.OFGateway, clientInput.OFUsername, openfaasPass, false)
	if isError(err) {
		echoDialogInfo(err.Error(), "")
		echoOpenDeploy(false)
		fmt.Println(err.Error())
		return
	}

	deploySignalName := make(map[string]string)
	deploySignalName["telar-web"] = "deployTelarWeb"
	deploySignalName["ts-serverless"] = "deployTsServerless"
	deploySignalName["ts-ui"] = "deploySocialUi"

	for _, microName := range []string{"telar-web", "ts-serverless", "ts-ui"} {
		microPath := path.Join(clientInput.ProjectDirectory, microName)
		if microName != "ts-ui" {
			faasPullGoTemplate(microPath)
		} else {
			faasPullNodeTemplate(microPath)
		}
		fmt.Printf("\n[INFO] Deploying %s function on OpenFaaS using %s gateway", microName, clientInput.OFGateway)
		err = faasDeploy(microPath, clientInput.OFGateway, authConf.Token)
		if isError(err) {
			echoDialogInfo(err.Error(), "")
			echoOpenDeploy(false)
			fmt.Println(err.Error())
			return
		}
		echoInput(deploySignalName[microName], true)
		fmt.Printf("\n[INFO] %s deployed!", microName)
	}
	echoSetupState("done")

}

func getSetupYaml(projectPath string) (*OFCSetupCache, error) {
	filePath := path.Join(projectPath, SETUP_YAML_FILE_NAME)
	fmt.Println("[INFO] Reading setup cache from ", filePath)

	bytesOut, readErr := ioutil.ReadFile(filePath)
	if readErr != nil {
		return nil, readErr
	}

	cacheData := ClientInputs{}
	unmarshalErr := yaml.Unmarshal(bytesOut, &cacheData)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return &OFCSetupCache{ClientInputs: cacheData}, nil
}

func writeSetupCache(projectPath string, clientInputs ClientInputs) {
	filePath := path.Join(projectPath, SETUP_YAML_FILE_NAME)
	writeYamlFile(filePath, &clientInputs)
}
