package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/Qolzam/telar-cli/pkg/log"
)

var ofCommandTimeout = 60 * time.Second

func checkKubesealPublicKey(pathWD string) error {
	var file, err = os.OpenFile(pathWD+"/pub-cert.pem", os.O_RDONLY, 0444)
	if isError(err) {
		return err
	}
	defer file.Close()
	return nil
}

func faasPullTemplates(path string) error {

	cmd := exec.Command("faas-cli", "template", "pull")
	cmd.Dir = path
	out, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("Faas template: %s - %s", err.Error(), string(out))
	}
	return nil
}

func faasPullNodeTemplate(path string) error {

	cmd := exec.Command("faas-cli", "template", "pull", "https://github.com/openfaas-incubator/node10-express-template")
	cmd.Dir = path
	out, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("Faas template: %s - %s", err.Error(), string(out))
	}
	return nil
}

func faasPullGoTemplate(path string) error {

	cmd := exec.Command("faas-cli", "template", "store", "pull", "golang-middleware")
	cmd.Dir = path
	out, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("Faas template: %s - %s", err.Error(), string(out))
	}
	return nil
}

func getOpenFaasPass() (string, error) {
	cmdFaasLogin := fmt.Sprintf("echo $(kubectl get secret -n openfaas basic-auth -o jsonpath=\"{.data.basic-auth-password}\" | base64 --decode; echo)")
	cmdBash := fmt.Sprintf("%s;", cmdFaasLogin)
	cmd := exec.Command("/bin/sh", "-c", cmdBash)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("OpenFaaS password: %s - %s", err.Error(), string(out))
	} else {
		log.Info("OpenFaaS password %s", string(out))
	}
	cmd.Wait()
	return string(out), nil
}

func faasDeploy(path, gateway, token string, env *[]string) error {
	cmd := exec.Command("faas-cli", "deploy", "--gateway", gateway)
	cmd.Dir = path
	if env != nil {
		cmd.Env = append(os.Environ(), *env...)
	}
	out, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("Faas deploy: %s - %s", err.Error(), string(out))
	}
	return nil
}
