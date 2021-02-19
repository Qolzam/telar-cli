package cmd

import (
	"fmt"
	"path/filepath"
)

func prepareSecret(pathWD, name, namespace string, kubeConfigPath *string, telarSecrets *TelarSecrets) error {

	err := createSecret(pathWD, name, namespace, kubeConfigPath, telarSecrets)
	if isError(err) {
		return err
	}

	return nil
}

func createSecret(pathWD, name, namespace string, kubeConfigPath *string, telarSecrets *TelarSecrets) error {

	args := make(map[string]string)
	args["mongo-pwd"] = telarSecrets.MongoPwd
	args["recaptcha-key"] = telarSecrets.RecaptchaKey
	args["ts-client-secret"] = telarSecrets.TsClientSecret
	args["redis-pwd"] = telarSecrets.RedisPwd
	args["admin-username"] = telarSecrets.AdminUsername
	args["admin-password"] = telarSecrets.AdminPwd
	args["payload-secret"] = telarSecrets.PayloadSecret
	args["ref-email-pass"] = telarSecrets.RefEmailPwd
	args["phone-auth-token"] = telarSecrets.PhoneAuthToken
	args["phone-auth-id"] = telarSecrets.PhoneAuthId
	saPath := pathWD + "/serviceAccountKey.json"
	publicKeyPath := pathWD + "/key.pub"
	privateKeyPath := pathWD + "/key"
	files := []string{saPath, publicKeyPath, privateKeyPath}
	err := kubectlCreateSecret(pathWD, name, namespace, kubeConfigPath, args, files)
	if isError(err) {
		return fmt.Errorf("Kubectl create secret, %s", err.Error())
	}
	secretFileName := fmt.Sprintf("%s-%s.yml", namespace, name)
	secretsYamlPath := filepath.Join(pathWD, secretFileName)
	err = kubectlApplyFile(secretsYamlPath, kubeConfigPath)
	if isError(err) {
		return fmt.Errorf("Kubectl apply file, %s", err.Error())
	}
	return nil
}

func createK8SSecret(pathWD, name string, telarSecrets *TelarSecrets) error {

	args := make(map[string]string)
	args["mongo-pwd"] = telarSecrets.MongoPwd
	args["recaptcha-key"] = telarSecrets.RecaptchaKey
	args["ts-client-secret"] = telarSecrets.TsClientSecret
	args["redis-pwd"] = telarSecrets.RedisPwd
	args["admin-username"] = telarSecrets.AdminUsername
	args["admin-password"] = telarSecrets.AdminPwd
	args["payload-secret"] = telarSecrets.PayloadSecret
	args["ref-email-pass"] = telarSecrets.RefEmailPwd
	args["phone-auth-token"] = telarSecrets.PhoneAuthToken
	args["phone-auth-id"] = telarSecrets.PhoneAuthId
	saPath := pathWD + "/serviceAccountKey.json"
	publicKeyPath := pathWD + "/key.pub"
	privateKeyPath := pathWD + "/key"
	files := []string{saPath, publicKeyPath, privateKeyPath}
	return runCloudSeal(name, pathWD, args, &files)
}
