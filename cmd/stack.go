package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"github.com/jinzhu/copier"
	stack "github.com/openfaas/faas-cli/stack"
	"gopkg.in/yaml.v2"
)

func applyConfig(telarConfig TelarConfig) error {
	for _, repo := range []string{"telar-web", "ts-serverless"} {
		err := applyAppConfig(telarConfig.PathWD+"/"+repo, telarConfig.MongoUser, telarConfig.MongoDatabase, telarConfig.RecaptchaSiteKey, telarConfig.RefEmail)
		if isError(err) {
			return err
		}
	}
	return nil
}

func applyAppConfig(repoPath, mongoUser, mongoDatabase, recaptchaSiteKey, refEmail string) error {
	filePath := "config/app_config.yml"
	envs, err := readConfigFile(repoPath, filePath)
	if isError(err) {
		return err
	}

	envs["mongo_user"] = mongoUser
	envs["mongo_database"] = mongoDatabase
	envs["recaptcha_site_key"] = recaptchaSiteKey
	envs["ref_email"] = refEmail
	envYaml := make(map[string]interface{})
	envYaml["environment"] = envs
	err = writeYamlFile(repoPath+"/"+filePath, &envYaml)
	if isError(err) {
		return err
	}
	return nil
}

func applyGatewayConfig(repoPath, coockieDomain, url, websocketURL string) error {
	filePath := "config/gateway_config.yml"
	envs, err := readConfigFile(repoPath, filePath)
	if isError(err) {
		return err
	}

	envs["gateway"] = url
	envs["origin"] = url
	envs["web_domain"] = url
	envs["external_domain"] = url
	envs["cookie_root_domain"] = coockieDomain
	envs["external_redirect_domain"] = url + "/auth"
	envs["websocket_server_url"] = websocketURL

	envYaml := make(map[string]interface{})
	envYaml["environment"] = envs
	err = writeYamlFile(repoPath+"/"+filePath, &envYaml)
	if isError(err) {
		return err
	}
	return nil
}

func applyStorageConfig(repoPath string, bucket string) error {
	filePath := "config/storage_config.yml"
	envs, err := readConfigFile(repoPath, filePath)
	if isError(err) {
		return err
	}

	bucketName := fmt.Sprintf("%s.appspot.com", bucket)
	envs["bucket_name"] = bucketName

	envYaml := make(map[string]interface{})
	envYaml["environment"] = envs
	err = writeYamlFile(repoPath+"/"+filePath, &envYaml)
	if isError(err) {
		return err
	}
	return nil
}

func applyAuthConfig(repoPath string, clientID string) error {
	filePath := "config/auth_config.yml"
	envs, err := readConfigFile(repoPath, filePath)
	if isError(err) {
		return err
	}

	envs["client_id"] = clientID

	envYaml := make(map[string]interface{})
	envYaml["environment"] = envs
	err = writeYamlFile(repoPath+"/"+filePath, &envYaml)
	if isError(err) {
		return err
	}
	return nil
}

func createAllStacks(pathWD string) {
	for _, repo := range []string{"ts-serverless", "telar-web"} {
		createStack(pathWD + "/" + repo)
	}
}

func createStack(pathWD string) {

	fmt.Printf("Current address: %s \n", pathWD)
	stackFile, _ := stack.ParseYAMLFile(path.Join(pathWD, "stack-init.yml"), "", "", false)

	for name, function := range stackFile.Functions {
		fmt.Println("map:", function.EnvironmentFile, name)
		//read environment variables from the file
		fileEnvironment, err := readFiles(function.EnvironmentFile, pathWD)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		//combine all environment variables
		allEnvironment, envErr := compileEnvironment([]string{}, function.Environment, fileEnvironment)
		if envErr != nil {
			log.Fatalf("error: %v", envErr)
		}

		// Set environments
		newFuncs := stack.Function{}
		copier.Copy(&newFuncs, stackFile.Functions[name])
		newFuncs.Environment = allEnvironment

		stackFile.Functions[name] = newFuncs
	}

	d, err := yaml.Marshal(&stackFile)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	errWrite := ioutil.WriteFile(path.Join(pathWD, "stack.yml"), d, 0644)
	if errWrite != nil {
		log.Fatalf("error: %v", err)
	}
}

func createTSUIStack(pathWD string, websocketURL string) error {
	filePath := path.Join(pathWD, "ts-ui/stack.yml")
	stackFile, err := stack.ParseYAMLFile(filePath, "", "", false)
	if err != nil {
		return err
	}
	for name, function := range stackFile.Functions {

		if name == "web" {
			envs := make(map[string]string)
			envs["websocket_url"] = websocketURL
			//combine all environment variables
			allEnvironment, envErr := compileEnvironment([]string{}, function.Environment, envs)
			if envErr != nil {
				return envErr
			}

			// Set environments
			newFuncs := stack.Function{}
			copier.Copy(&newFuncs, stackFile.Functions[name])
			newFuncs.Environment = allEnvironment
			stackFile.Functions[name] = newFuncs
			break
		}

	}

	d, err := yaml.Marshal(&stackFile)
	if err != nil {
		return err
	}

	errWrite := ioutil.WriteFile(filePath, d, 0644)
	if errWrite != nil {
		return err
	}
	return nil
}
