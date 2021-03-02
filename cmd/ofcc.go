package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/Qolzam/telar-cli/cmd/ui"
	"github.com/Qolzam/telar-cli/pkg/content"
	"github.com/Qolzam/telar-cli/pkg/log"
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

	log.Info("Starting download of pub-cert.pe %s, this could take a few moments.\n", downloadURL)
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

}

func echoProjectDir() {
	projectPath, err := getDefaultProjectDirectory()
	if err != nil {
		echoDialogInfo(err.Error(), "https://github.com/Qolzam/telar-cli/blob/master/docs/ofcc-setup/1.md")
		return
	}
	echoInput("projectDirectory", projectPath)
}

func CheckInitStep(input *ClientInputs) {
	_ = checkDirectory(input.ProjectDirectory)

	setupCache, err := getSetupYaml(input.ProjectDirectory)
	if err != nil {
		log.Info("\n[WARN] can not get setup cache. %s", err.Error())
	} else {
		setupCache.ClientInputs.ProjectDirectory = input.ProjectDirectory
		setupCache.ClientInputs.AppID = input.AppID
		setupCache.ClientInputs.AppName = input.AppName
		setupCache.ClientInputs.CompanyName = input.CompanyName
		setupCache.ClientInputs.SupportEmail = input.SupportEmail
		echoSetupDefaultValues(setupCache)
	}
	echoStep(1)

}

func OFCAccessSetting() {

	echoStep(2)
}
func CheckIngredient(projectPath string, githubUsername string) {
	helpURL := "https://github.com/Qolzam/telar-cli/blob/master/docs/ofcc-setup/3.md"
	echoInput("loadingCheckIngredients", true)

	// Check telar-web repository
	err := cloneTelarWeb(projectPath, githubUsername)
	if isError(err) && err.Error() != "repository already exists" {
		errMessage := "telar-web " + err.Error()

		echoDialogInfo(errMessage, helpURL)
		echoInput("loadingCheckIngredients", false)

		log.Error(errMessage)
		return
	}
	echoInput("cloneTelarWeb", true)

	// Check ts-serverless repository
	err = cloneTSServerless(projectPath, githubUsername)
	if isError(err) && err.Error() != "repository already exists" {
		errMessage := "ts-serverless " + err.Error()

		echoDialogInfo(errMessage, helpURL)
		echoInput("loadingCheckIngredients", false)

		log.Error(errMessage)
		return
	}
	echoInput("cloneTsServerless", true)

	// Check ts-ui repository
	err = cloneTSUI(projectPath, githubUsername)
	if isError(err) && err.Error() != "repository already exists" {
		errMessage := "ts-ui " + err.Error()

		echoDialogInfo(errMessage, helpURL)
		echoInput("loadingCheckIngredients", false)

		log.Error(errMessage)
		return
	}
	echoInput("cloneTsUi", true)

	echoStep(3)
}

func CheckStorage(projectPath string, bucketName string) {
	helpURL := "https://github.com/Qolzam/telar-cli/blob/master/docs/ofcc-setup/4.md"
	echoInput("loadingFirebaseStorage", true)

	// Check serviceAccount.json file for Firebase
	err := checkFirebaseServiceAccount(projectPath)
	if isError(err) {
		echoDialogInfo(err.Error(), helpURL)
		echoInput("loadingFirebaseStorage", false)
		log.Error(err.Error())
		return
	}
	echoInput("firebaseServiceAccount", true)

	// Check Firebase storage access
	err = checkFirebaseStorageBucket(projectPath, bucketName)
	if isError(err) {
		echoDialogInfo(err.Error(), helpURL)
		echoInput("loadingFirebaseStorage", false)
		log.Error(err.Error())
		return
	}
	echoInput("firebaseStorage", true)

	echoStep(4)

}

func CheckDatabase(mongoDBURI string) {
	helpURL := "https://github.com/Qolzam/telar-cli/blob/master/docs/ofcc-setup/5.md"

	echoInput("loadingMongoDB", true)

	err := checkDB(mongoDBURI)
	if isError(err) {
		echoDialogInfo(err.Error(), helpURL)
		echoInput("loadingMongoDB", false)
		log.Error(err.Error())
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
	helpURL := "https://github.com/Qolzam/telar-cli/blob/master/docs/ofcc-setup/8.md"

	payloadSecret, err := generatePayloadSecret()
	if isError(err) {
		echoDialogInfo(err.Error(), helpURL)
		log.Error(err.Error())
		return
	}
	echoInput("payloadSecret", payloadSecret)
	echoInput("gateway", ofGateway)
	echoStep(8)
}

func CheckWebsocket(clientInput ClientInputs) {
	// helpURL := "https://github.com/Qolzam/telar-cli/blob/master/docs/ofcc-setup/9.md"
	echoInput("loadingWebsocket", true)

	// err := pingWebsocket(clientInput.WebsocketURL)
	// if isError(err) {
	// 	echoDialogInfo(err.Error(), helpURL)
	// 	echoInput("loadingWebsocket", false)
	// 	log.Error(err.Error())
	// 	return
	// }
	echoInput("websocketConnection", true)
	starteploy(clientInput)
}

func starteploy(clientInput ClientInputs) {
	echoOpenDeploy(true)

	host := content.GetDomainFromURI(clientInput.SocialDomain)
	log.Info("Host: ", host)

	// Apply stack
	telarConfig := TelarConfig{
		AppID:            clientInput.AppID,
		SecretName:       clientInput.SecretName,
		GithubUsername:   TELAR_GITHUB_USER_NAME,
		PathWD:           clientInput.ProjectDirectory,
		CoockieDomain:    "." + host,
		Bucket:           clientInput.BucketName,
		ClientID:         clientInput.GithubOAuthClientID,
		Gateway:          GetGatewayURL(clientInput.OFGateway, "", "", os.Getenv(openFaaSURLEnvironment)),
		Origin:           clientInput.SocialDomain,
		WebsocketURL:     clientInput.WebsocketURL,
		MongoDBURI:       clientInput.MongoDBURI,
		MongoDatabase:    clientInput.MongoDBName,
		RecaptchaSiteKey: clientInput.SiteKeyRecaptcha,
		RefEmail:         clientInput.Gmail,
	}

	err := applyConfig(telarConfig)
	if isError(err) {
		echoDialogInfo(err.Error(), "")
		echoOpenDeploy(false)
		log.Error(err.Error())
		return
	}
	echoInput("loadingStackYaml", true)

	err = preparePublicPrivateKey(clientInput.ProjectDirectory)
	if isError(err) {
		echoDialogInfo(err.Error(), "")
		echoOpenDeploy(false)
		log.Error(err.Error())
		return
	}
	echoInput("loadingPublicPrivateKey", true)

	// Create secrets
	telarSecret := &TelarSecrets{
		MongoURI:       clientInput.MongoDBURI,
		MongoDB:        clientInput.MongoDBName,
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
		log.Error(err.Error())
		return
	}
	echoInput("loadingCreateSecret", true)

	// Deploy Telar Web
	openfaasPass, err := getOpenFaasPass()
	if isError(err) {
		echoDialogInfo(err.Error(), "")
		echoOpenDeploy(false)
		log.Error(err.Error())
		return
	}

	authConf, _, err := runFaaSLogin(clientInput.OFGateway, clientInput.OFUsername, openfaasPass, false)
	if isError(err) {
		echoDialogInfo(err.Error(), "")
		echoOpenDeploy(false)
		log.Error(err.Error())
		return
	}

	deploySignalName := make(map[string]string)
	deploySignalName["telar-web"] = "deployTelarWeb"
	deploySignalName["ts-serverless"] = "deployTsServerless"
	deploySignalName["ts-ui"] = "deploySocialUi"

	for _, microName := range []string{"telar-web", "ts-serverless", "ts-ui"} {

		microPath := path.Join(clientInput.ProjectDirectory, microName)
		log.Info("Deploying %s function on OpenFaaS using %s gateway", microName, clientInput.OFGateway)
		var env []string
		if microName == "ts-ui" {
			faasPullNodeTemplate(microPath)
			envDockerUser := "DOCKER_USER=" + clientInput.DockerUser
			uiStackVersion, err := ui.ReadStackVersion(microPath)
			if isError(err) {
				echoDialogInfo(err.Error(), "")
				echoOpenDeploy(false)
				log.Error(err.Error())
				return
			}
			envStackVer := "STACK_VER=" + uiStackVersion
			env = []string{envDockerUser, envStackVer}
			err = ui.UIUp(&ui.UIConfig{
				UIPath:       microPath,
				DockerUser:   clientInput.DockerUser,
				StackVersion: "",
				Gateway:      GetGatewayURL(clientInput.OFGateway, "", "", os.Getenv(openFaaSURLEnvironment)),
				BaseAPIRoute: clientInput.BaseAPIRoute,
				AppName:      clientInput.AppName,
				CompanyName:  clientInput.CompanyName,
				SupportEmail: clientInput.SupportEmail,
				WSURL:        clientInput.WebsocketURL,
				Env:          env,
			})
			if isError(err) {
				echoDialogInfo(err.Error(), "")
				echoOpenDeploy(false)
				log.Error(err.Error())
				return
			}
		} else {
			faasPullGoTemplate(microPath)
		}
		err = faasDeploy(microPath, clientInput.OFGateway, authConf.Token, &env)
		if isError(err) {
			echoDialogInfo(err.Error(), "")
			echoOpenDeploy(false)
			log.Error(err.Error())
			return
		}
		echoInput(deploySignalName[microName], true)
		log.Info(" %s deployed!", microName)
	}
	echoSetupState("done")

}

func getSetupYaml(projectPath string) (*OFCSetupCache, error) {
	filePath := path.Join(projectPath, SETUP_YAML_FILE_NAME)
	log.Info("Reading setup cache from ", filePath)

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
