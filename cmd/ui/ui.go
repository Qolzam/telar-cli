package ui

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/Qolzam/telar-cli/pkg/env"
	"github.com/Qolzam/telar-cli/pkg/log"
)

type UIConfig struct {
	UIPath       string
	DockerUser   string
	StackVersion string
	Gateway      string
	BaseAPIRoute string
	AppName      string
	CompanyName  string
	SupportEmail string
	WSURL        string
	BaseHref     string
	Env          []string
}

// UIUp config/build/push UI
func UIUp(config *UIConfig) error {
	log.Info("Start building UI")
	err := SetUIEnv(config.UIPath, config.Gateway, config.BaseAPIRoute, config.AppName, config.CompanyName, config.SupportEmail, config.WSURL)
	if err != nil {
		return err
	}
	log.Info("[SetUIEnv]")
	buildErr := BuildUI(config.UIPath, &config.Env)
	if buildErr != nil {
		return buildErr
	}
	return nil

}

// BuildUI build UI
func BuildUI(uiPath string, env *[]string) error {
	cmd := exec.Command("make", "prepare-push")
	cmd.Dir = uiPath
	if env != nil {
		cmd.Env = append(os.Environ(), *env...)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("[BuildUI] %s", err.Error())
	}
	return err
}

// SetUIEnv Set UI environment
func SetUIEnv(uiPath, gateway, baseAPIRoute, appName, companyName, supportEmail, wsURL string) error {
	configPath := path.Join(uiPath, ".env.production")
	config, err := env.ReadEnvFile(configPath)
	if err != nil { // Handle errors reading the config file
		return fmt.Errorf("Fatal error read config file: %s \n", err)
	}

	config["REACT_APP_GATEWAY"] = gateway
	config["REACT_APP_BASE_ROUTE_API"] = baseAPIRoute
	config["REACT_APP_NAME"] = appName
	config["REACT_APP_COMPANY_NAME"] = companyName
	config["REACT_APP_EMAIL_SUPPORT"] = supportEmail
	config["REACT_APP_WEBSOCKET_URL"] = wsURL

	err = env.WriteEnvFile(configPath, &config)
	if err != nil { // Handle errors reading the config file
		return fmt.Errorf("Fatal error write config file: %s \n", err)
	}
	return nil
}

// ReadStackVersion read stack version
func ReadStackVersion(uiPath string) (string, error) {
	bytesOut, readErr := ioutil.ReadFile(path.Join(uiPath, "current_version"))
	if readErr != nil {
		return "", readErr
	}

	return string(bytesOut), nil
}
