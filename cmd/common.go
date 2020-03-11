package cmd

// Copyright (c) Alex Ellis 2017. All rights reserved.
// This script was adapted from https://github.com/openfaas/faas-cli/blob/master/commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strconv"
	"strings"

	execute "github.com/alexellis/go-execute/pkg/v1"
	stack "github.com/openfaas/faas-cli/stack"
	"gopkg.in/yaml.v2"
)

var (
	downloadVersion string
	downloadTo      string
)

const (
	OFCC_DOMAIN = ".o6s.io"

	// Client Actions
	SET_SETUP_STATE    = "SET_SETUP_STATE"
	SET_SETUP_STEP     = "SET_SETUP_STEP"
	SET_INPUT          = "SET_INPUT"
	SET_DEPLOY_OPEN    = "SET_DEPLOY_OPEN"
	POP_MESSAGE        = "POP_MESSAGE"
	SET_STEP_CONDITION = "SET_STEP_CONDITION"
	SHOW_INFO_DIALOG   = "SHOW_INFO_DIALOG"

	// Server HTTP Actions
	START_STEP = "START_STEP"
)

type TelarSecrets struct {
	MongoPwd,
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

type TelarConfig struct {
	GithubUsername   string `json:"githubUsername"`
	PathWD           string `json:"pathWD"`
	CoockieDomain    string `json:"coockieDomain"`
	Bucket           string `json:"bucket"`
	ClientID         string `json:"clientID"`
	URL              string `json:"url"`
	WebsocketURL     string `json:"websocketURL"`
	MongoUser        string `json:"mongoUser"`
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
		log.Fatal(err)
	}

	// output has trailing \n
	// need to remove the \n
	// otherwise it will cause error for strconv.Atoi
	// log.Println(output[:len(output)-1])

	// 0 = root, 501 = non-root user
	i, err := strconv.Atoi(string(output[:len(output)-1]))

	if err != nil {
		log.Fatal(err)
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

	if err := os.Mkdir(dir, 0755); os.IsExist(err) {
		return true
	}

	return false
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
