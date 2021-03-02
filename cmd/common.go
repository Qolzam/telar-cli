package cmd

// Copyright (c) Alex Ellis 2017. All rights reserved.
// This script was adapted from https://github.com/openfaas/faas-cli/blob/master/commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Qolzam/telar-cli/pkg/log"
	execute "github.com/alexellis/go-execute/pkg/v1"
	stack "github.com/openfaas/faas-cli/stack"
	"gopkg.in/yaml.v2"
)

var (
	downloadVersion string
	downloadTo      string
)

const (
	TELAR_GITHUB_USER_NAME = "red-gold"
	IMAGE_OWNER            = "telar"
	REGISTRY_URL           = "docker.io/telar/"

	// Client Actions
	SET_SETUP_STATE          = "SET_SETUP_STATE"
	SET_SETUP_STEP           = "SET_SETUP_STEP"
	SET_INPUT                = "SET_INPUT"
	SET_DEPLOY_OPEN          = "SET_DEPLOY_OPEN"
	SET_SETUP_DEFAULT_VALUES = "SET_SETUP_DEFAULT_VALUES"
	POP_MESSAGE              = "POP_MESSAGE"
	SET_STEP_CONDITION       = "SET_STEP_CONDITION"
	SHOW_INFO_DIALOG         = "SHOW_INFO_DIALOG"

	// Server HTTP Actions
	START_STEP                 = "START_STEP"
	REMOVE_SOCIAL_FROM_CLUSTER = "REMOVE_SOCIAL_FROM_CLUSTER"
	ECHO_PROJECT_DIR           = "ECHO_PROJECT_DIR"
	CHECK_STEP                 = "CHECK_STEP"
)

type TelarSecrets struct {
	MongoURI,
	MongoDB,
	RecaptchaKey,
	TsClientSecret,
	RedisPwd,
	AdminUsername,
	AdminPwd,
	PayloadSecret,
	RefEmailPwd,
	PhoneAuthToken,
	PhoneAuthId string
}

type ClientState struct {
	SetupState string       `json:"setupState"`
	SetupStep  int          `json:"setupStep"`
	Inputs     ClientInputs `json:"inputs"`
}

type ClientInputs struct {
	BaseAPIRoute        string `json:"baseAPIRoute"  yaml:"baseAPIRoute,omitempty"`
	AppName             string `json:"appName"  yaml:"appName,omitempty"`
	CompanyName         string `json:"companyName"  yaml:"companyName,omitempty"`
	SupportEmail        string `json:"supportEmail"  yaml:"supportEmail,omitempty"`
	AppID               string `json:"appID"  yaml:"appID,omitempty"`
	OFUsername          string `json:"ofUsername"  yaml:"ofUsername,omitempty"`
	OFGateway           string `json:"ofGateway"  yaml:"ofGateway,omitempty"`
	SocialDomain        string `json:"socialDomain"  yaml:"socialDomain,omitempty"`
	SecretName          string `json:"secretName"  yaml:"secretName,omitempty"`
	Namespace           string `json:"namespace"  yaml:"namespace,omitempty"`
	DockerUser          string `json:"dockerUser"  yaml:"dockerUser,omitempty"`
	KubeconfigPath      string `json:"kubeconfigPath"  yaml:"kubeconfigPath,omitempty"`
	ProjectDirectory    string `json:"projectDirectory"  yaml:"projectDirectory,omitempty"`
	BucketName          string `json:"bucketName"  yaml:"bucketName,omitempty"`
	MongoDBURI          string `json:"mongoDBURI"  yaml:"mongoDBURI,omitempty"`
	MongoDBName         string `json:"mongoDBName"  yaml:"mongoDBName,omitempty"`
	SiteKeyRecaptcha    string `json:"siteKeyRecaptcha"  yaml:"siteKeyRecaptcha,omitempty"`
	RecaptchaKey        string `json:"recaptchaKey"  yaml:"recaptchaKey,omitempty"`
	GithubOAuthSecret   string `json:"githubOAuthSecret"  yaml:"githubOAuthSecret,omitempty"`
	GithubOAuthClientID string `json:"githubOAuthClientID"  yaml:"githubOAuthClientID,omitempty"`
	AdminUsername       string `json:"adminUsername"  yaml:"adminUsername,omitempty"`
	AdminPassword       string `json:"adminPassword"  yaml:"adminPassword,omitempty"`
	Gmail               string `json:"gmail"  yaml:"gmail,omitempty"`
	GmailPassword       string `json:"gmailPassword"  yaml:"gmailPassword,omitempty"`
	Gateway             string `json:"gateway"  yaml:"gateway,omitempty"`
	PayloadSecret       string `json:"payloadSecret"  yaml:"payloadSecret,omitempty"`
	WebsocketURL        string `json:"websocketURL"  yaml:"websocketURL,omitempty"`
}

type TelarConfig struct {
	AppID            string `json:"ppID"`
	SecretName       string `json:"secretName"`
	GithubUsername   string `json:"githubUsername"`
	PathWD           string `json:"pathWD"`
	CoockieDomain    string `json:"coockieDomain"`
	Bucket           string `json:"bucket"`
	ClientID         string `json:"clientID"`
	Gateway          string `json:"gateway"`
	Origin           string `json:"origin"`
	WebsocketURL     string `json:"websocketURL"`
	MongoDBURI       string `json:"mongoDBURI"`
	MongoDatabase    string `json:"mongoDatabase"`
	RecaptchaSiteKey string `json:"recaptchaSiteKey"`
	RefEmail         string `json:"refEmail"`
}

type UIGatewayConfig struct {
	Websocket UIWebsocketConfig `json:"websocket"`
}

type UIWebsocketConfig struct {
	URL string `json:"url"`
}

func readFiles(files []string, rootPath string) (map[string]string, error) {
	envs := make(map[string]string)

	for _, file := range files {
		bytesOut, readErr := ioutil.ReadFile(path.Join(rootPath, file))
		if readErr != nil {
			return nil, readErr
		}

		envFile := stack.EnvironmentFile{}
		unmarshalErr := yaml.Unmarshal(bytesOut, &envFile)
		if unmarshalErr != nil {
			return nil, unmarshalErr
		}
		for k, v := range envFile.Environment {
			envs[k] = v
		}

	}
	return envs, nil
}

func readConfigFile(rootPath string, file string) (map[string]string, error) {
	envs := make(map[string]string)
	bytesOut, readErr := ioutil.ReadFile(path.Join(rootPath, file))
	if readErr != nil {
		return nil, readErr
	}

	envFile := stack.EnvironmentFile{}
	unmarshalErr := yaml.Unmarshal(bytesOut, &envFile)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}
	for k, v := range envFile.Environment {
		envs[k] = v
	}
	return envs, nil
}

func compileEnvironment(envvarOpts []string, yamlEnvironment map[string]string, fileEnvironment map[string]string) (map[string]string, error) {
	envvarArguments, err := parseMap(envvarOpts, "env")
	if err != nil {
		return nil, fmt.Errorf("error parsing envvars: %v", err)
	}

	functionAndStack := mergeMap(yamlEnvironment, fileEnvironment)
	return mergeMap(functionAndStack, envvarArguments), nil
}

func parseMap(envvars []string, keyName string) (map[string]string, error) {
	result := make(map[string]string)
	for _, envvar := range envvars {
		s := strings.SplitN(strings.TrimSpace(envvar), "=", 2)
		if len(s) != 2 {
			return nil, fmt.Errorf("label format is not correct, needs key=value")
		}
		envvarName := s[0]
		envvarValue := s[1]

		if !(len(envvarName) > 0) {
			return nil, fmt.Errorf("empty %s name: [%s]", keyName, envvar)
		}
		if !(len(envvarValue) > 0) {
			return nil, fmt.Errorf("empty %s value: [%s]", keyName, envvar)
		}

		result[envvarName] = envvarValue
	}
	return result, nil
}

func mergeMap(i map[string]string, j map[string]string) map[string]string {
	merged := make(map[string]string)

	for k, v := range i {
		merged[k] = v
	}
	for k, v := range j {
		merged[k] = v
	}
	return merged
}

func isError(err error) bool {
	return err != nil
}

func chmod(pathWD string) error {
	task := execute.ExecTask{
		Command: fmt.Sprintf("chmod +x %s", pathWD),
	}

	_, taskErr := task.Execute()

	if taskErr != nil {
		return taskErr
	}
	return nil
}

func install(pathWD string) error {
	task := execute.ExecTask{
		Command: fmt.Sprintf("sudo install %s /usr/local/bin/", pathWD),
	}

	_, taskErr := task.Execute()

	if taskErr != nil {
		return taskErr
	}
	return nil
}

// Check running program in sudo
// https://www.socketloop.com/tutorials/golang-force-your-program-to-run-with-root-permissions
func checkSudo() bool {
	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()

	if err != nil {
		log.Error(err.Error())
	}

	// output has trailing \n
	// need to remove the \n
	// otherwise it will cause error for strconv.Atoi
	// log.Println(output[:len(output)-1])

	// 0 = root, 501 = non-root user
	i, err := strconv.Atoi(string(output[:len(output)-1]))

	if err != nil {
		log.Error(err.Error())
	}

	if i == 0 {
		return true
	} else {
		return false
	}
}

func getDefaultProjectDirectory() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/telar-social", usr.HomeDir), nil

}

func checkDirectory(dir string) bool {

	if err := os.MkdirAll(dir, 0755); os.IsExist(err) {
		return true
	}

	return false
}

func directoryFileExist(dir string) (error, bool) {
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return nil, false
		} else {
			return err, true
		}
	}
	return nil, true
}

func writeYamlFile(path string, yamlData interface{}) error {
	data, err := yaml.Marshal(yamlData)
	if err != nil {
		return &json.UnmarshalFieldError{}
	}
	err = ioutil.WriteFile(path, data, 0644)
	if isError(err) {
		return err
	}
	return nil
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func interfaceToMapString(in interface{}) map[string]string {
	stringMap := make(map[string]string)
	switch v := in.(type) {
	case map[string]string:
		for key, value := range v {
			stringMap[key] = value
		}
	}
	return stringMap
}

func generatePayloadSecret() (string, error) {
	task := execute.ExecTask{
		Command: `echo`,
		Args:    []string{`-n $(head -c 16 /dev/urandom | shasum | cut -d " " -f 1)`},
		Shell:   true,
	}
	taskExe, taskErr := task.Execute()
	if taskErr != nil {
		return "", taskErr
	}

	payloadSecret := fmt.Sprintf("%s", taskExe.Stdout)
	log.Info(payloadSecret)
	return payloadSecret, nil
}

func createPrivateKey(path string) error {
	arg := fmt.Sprintf(`-n $(openssl ecparam -genkey -name prime256v1 -noout -out %s/key)`, path)
	println(arg)
	task := execute.ExecTask{
		Command: `echo`,
		Args:    []string{arg},
		Shell:   true,
	}
	_, taskErr := task.Execute()
	if taskErr != nil {
		return taskErr
	}
	return nil
}

func createPublicKey(path string) error {
	arg := fmt.Sprintf(`-n $(openssl ec -in %s/key -pubout -out %s/key.pub)`, path, path)
	println(arg)
	task := execute.ExecTask{
		Command: `echo`,
		Args:    []string{arg},
		Shell:   true,
	}
	_, taskErr := task.Execute()
	if taskErr != nil {
		return taskErr
	}
	return nil
}

func preparePublicPrivateKey(path string) error {
	err := createPrivateKey(path)
	if isError(err) {
		return err
	}
	err = createPublicKey(path)
	if isError(err) {
		return err
	}
	return nil
}

// To provide input to the pipeline, assign an io.Reader to the first's Stdin.
func Pipeline(cmds ...*exec.Cmd) (pipeLineOutput, collectedStandardError []byte, pipeLineError error) {
	// Require at least one command
	if len(cmds) < 1 {
		return nil, nil, nil
	}

	// Collect the output from the command(s)
	var output bytes.Buffer
	var stderr bytes.Buffer

	last := len(cmds) - 1
	for i, cmd := range cmds[:last] {
		var err error
		// Connect each command's stdin to the previous command's stdout
		if cmds[i+1].Stdin, err = cmd.StdoutPipe(); err != nil {
			return nil, nil, err
		}
		// Connect each command's stderr to a buffer
		cmd.Stderr = &stderr
	}

	// Connect the output and error for the last command
	cmds[last].Stdout, cmds[last].Stderr = &output, &stderr

	// Start each command
	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	// Wait for each command to complete
	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	// Return the pipeline output and the collected standard error
	return output.Bytes(), stderr.Bytes(), nil
}

func kubectlCreateSecret(path, name, namespace string, kubeConfigPath *string, args map[string]string, files []string) error {

	cmdArgs := []string{"kubectl", "-n", namespace, "create", "secret", "generic", name,
		"--save-config", "--dry-run=client"}
	for k, v := range args {
		newArg := fmt.Sprintf("--from-literal='%s=%s'", k, v)
		cmdArgs = append(cmdArgs, newArg)
	}
	for _, v := range files {
		newArg := fmt.Sprintf("--from-file='%s'", v)
		cmdArgs = append(cmdArgs, newArg)
	}
	cmdArgs = append(cmdArgs, "-o yaml")
	cmdBash := fmt.Sprintf("cd %s; %s;", path, strings.Join(cmdArgs[:], " "))
	if kubeConfigPath != nil {
		cmdExportKubeConfig := fmt.Sprintf("export KUBECONFIG=%s", *kubeConfigPath)
		cmdBash = fmt.Sprintf("cd %s; %s; %s;", path, cmdExportKubeConfig, strings.Join(cmdArgs[:], " "))

	}
	log.Info("Create secret final command ", cmdBash)
	cmd := exec.Command("/bin/sh", "-c", cmdBash)
	secretFileName := fmt.Sprintf("%s-%s.yml", namespace, name)
	secretsYamlPath := filepath.Join(path, secretFileName)
	outfile, err := os.Create(secretsYamlPath)
	if err != nil {
		return fmt.Errorf("Can not create file in path %s, error: %s", secretsYamlPath, err.Error())
	}
	defer outfile.Close()
	cmd.Stdout = outfile

	err = cmd.Start()
	if err != nil {
		return err
	}
	cmd.Wait()
	return nil
}

func kubectlApplyFile(path string, kubeConfigPath *string) error {
	cmdApply := "kubectl apply -f " + path
	cmdBash := fmt.Sprintf("%s;", cmdApply)
	if kubeConfigPath != nil {
		cmdExportKubeConfig := fmt.Sprintf("export KUBECONFIG=%s", *kubeConfigPath)
		cmdBash = fmt.Sprintf("%s; %s;", cmdExportKubeConfig, cmdApply)

	}
	log.Info("Kubectl apply final command ", cmdBash)
	out, err := exec.Command("/bin/sh", "-c", cmdBash).CombinedOutput()

	if isError(err) {
		return fmt.Errorf("%s - %s", err.Error(), string(out))
	}
	return nil
}
