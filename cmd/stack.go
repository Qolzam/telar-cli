package cmd

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/jinzhu/copier"
	stack "github.com/openfaas/faas-cli/stack"
	"gopkg.in/yaml.v2"
)

func applyConfig(telarConfig TelarConfig) error {
	for _, repo := range []string{"telar-web", "ts-serverless"} {
		repoPath := telarConfig.PathWD + "/" + repo

		fmt.Printf("\n[INFO] Applying app config %s", repoPath)
		err := applyAppConfig(repoPath, telarConfig.MongoDBHost, telarConfig.MongoDatabase, telarConfig.RecaptchaSiteKey, telarConfig.RefEmail)
		if isError(err) {
			return err
		}
		fmt.Printf("\n[INFO] Applied app config successfully %s", repoPath)

		fmt.Printf("\n[INFO] Applying gateway config %s", repoPath)
		err = applyGatewayConfig(repoPath, telarConfig.CoockieDomain, telarConfig.URL, telarConfig.WebsocketURL)
		if isError(err) {
			return err
		}
		fmt.Printf("\n[INFO] Applied gateway config successfully %s", repoPath)
	}

	telarWebRepoPath := telarConfig.PathWD + "/" + "telar-web"

	fmt.Printf("\n[INFO] Applying auth config %s", telarWebRepoPath)
	err := applyAuthConfig(telarWebRepoPath, telarConfig.ClientID)
	if isError(err) {
		return err
	}
	fmt.Printf("\n[INFO] Applied auth config successfully %s", telarWebRepoPath)

	fmt.Printf("\n[INFO] Applying storage config %s", telarWebRepoPath)
	err = applyStorageConfig(telarWebRepoPath, telarConfig.Bucket)
	if isError(err) {
		return err
	}
	fmt.Printf("\n[INFO] Applied storage config successfully %s", telarWebRepoPath)

	fmt.Printf("\n[INFO] Applying ts-ui config %s", telarConfig.PathWD)
	err = applyTSUIConfig(telarConfig.PathWD, telarConfig.WebsocketURL, telarConfig.URL)
	if isError(err) {
		return err
	}
	fmt.Printf("\n[INFO] Applied ts-ui config successfully %s", telarConfig.PathWD)

	fmt.Printf("\n[INFO] Create all stack config %s", telarConfig.PathWD)
	err = createAllStacks(telarConfig.PathWD, telarConfig.SecretName, telarConfig.AppID)
	if isError(err) {
		return err
	}
	fmt.Printf("\n[INFO] All stack config created successfully %s", telarConfig.PathWD)
	return nil
}

func applyAppConfig(repoPath, mongoDBHost, mongoDatabase, recaptchaSiteKey, refEmail string) error {
	filePath := "config/app_config.yml"
	envs, err := readConfigFile(repoPath, filePath)
	if isError(err) {
		return err
	}

	envs["mongo_host"] = strings.Replace(mongoDBHost, "<password>", "%s", -1)
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

func createAllStacks(pathWD, secretName, appID string) error {
	for _, repo := range []string{"ts-serverless", "telar-web"} {

		stackPath := path.Join(pathWD, repo)
		fmt.Printf("\n[INFO] Create stack config %s", stackPath)
		err := createStack(stackPath, secretName, repo, appID)
		if err != nil {
			return err
		}
		fmt.Printf("\n[INFO] Stack yaml file created successfully  %s", stackPath)
	}

	tsuiStackPath := path.Join(pathWD, "ts-ui")
	err := createTSUIStack(tsuiStackPath, "ts-ui", appID)
	if err != nil {
		return err
	}
	fmt.Printf("\n[INFO] Stack yaml file created successfully  %s", tsuiStackPath)
	return nil
}

func createStack(pathWD, secretName, repo, appID string) error {

	fmt.Printf("[INFO] Creating stack.yml file: %s \n", pathWD)
	stackFile, _ := stack.ParseYAMLFile(path.Join(pathWD, "stack-init.yml"), "", "", false)
	sha, err := gitShortSHA(pathWD)
	if err != nil {
		return err
	}
	for name, function := range stackFile.Functions {
		fmt.Printf("\n%v", name)
		fmt.Printf("\n%v", function)

		// Set environments
		newFuncs := stack.Function{}
		copier.Copy(&newFuncs, stackFile.Functions[name])
		newFuncs.Secrets = []string{secretName}
		newLabels := make(map[string]string)
		newLabels[FunctionLabelPrefix+"repo"] = repo
		newLabels[FunctionLabelPrefix+"appID"] = appID
		mergedLables := mergeMap(newLabels, *newFuncs.Labels)
		newFuncs.Labels = &mergedLables
		newFuncs.Image = formatImageShaTag(REGISTRY_URL, &newFuncs, sha, TELAR_GITHUB_USER_NAME, repo)
		stackFile.Functions[name] = newFuncs
	}

	d, err := yaml.Marshal(&stackFile)
	if err != nil {
		return err
	}

	errWrite := ioutil.WriteFile(path.Join(pathWD, "stack.yml"), d, 0644)
	if errWrite != nil {
		return err
	}
	return nil
}

func createTSUIStack(pathWD, repo, appID string) error {

	fmt.Printf("[INFO] Creating stack.yml file: %s \n", pathWD)
	var stackFile *stack.Services
	err, stackInitExist := directoryFileExist(path.Join(pathWD, "stack-init.yml"))
	if stackInitExist {
		stackFile, _ = stack.ParseYAMLFile(path.Join(pathWD, "stack-init.yml"), "", "", false)

	} else {
		stackFile, _ = stack.ParseYAMLFile(path.Join(pathWD, "stack.yml"), "", "", false)
		d, err := yaml.Marshal(&stackFile)
		if err != nil {
			return err
		}
		errWrite := ioutil.WriteFile(path.Join(pathWD, "stack-init.yml"), d, 0644)
		if errWrite != nil {
			return err
		}
	}
	sha, err := gitShortSHA(pathWD)
	if err != nil {
		return err
	}
	for name, function := range stackFile.Functions {
		fmt.Printf("\n%v", name)
		fmt.Printf("\n%v", function)
		// Set environments
		newFuncs := stack.Function{}
		copier.Copy(&newFuncs, stackFile.Functions[name])
		newLabels := make(map[string]string)
		newLabels[FunctionLabelPrefix+"repo"] = repo
		newLabels[FunctionLabelPrefix+"appID"] = appID
		mergedLables := mergeMap(newLabels, *newFuncs.Labels)
		newFuncs.Image = formatImageShaTag(REGISTRY_URL, &newFuncs, sha, TELAR_GITHUB_USER_NAME, repo)
		newFuncs.Labels = &mergedLables

		stackFile.Functions[name] = newFuncs
	}

	d, err := yaml.Marshal(&stackFile)
	if err != nil {
		return err
	}

	errWrite := ioutil.WriteFile(path.Join(pathWD, "stack.yml"), d, 0644)
	if errWrite != nil {
		return err
	}
	return nil
}

func applyTSUIConfig(pathWD, websocketURL, gateway string) error {
	filePath := path.Join(pathWD, "ts-ui/stack.yml")
	stackFile, err := stack.ParseYAMLFile(filePath, "", "", false)
	if err != nil {
		return err
	}
	for name, function := range stackFile.Functions {

		if name == "web" {
			envs := make(map[string]string)
			envs["websocket_url"] = websocketURL
			envs["gateway_url"] = gateway
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
