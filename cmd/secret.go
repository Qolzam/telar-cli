package cmd

import (
	"fmt"
	"io/ioutil"
)

func prepareSecret(pathWD, githubUsername string, telarSecrets *TelarSecrets) error {
	secretFileName := "secrets.yml"
	err := createSecretFile(pathWD, githubUsername+"-secrets", telarSecrets)
	if isError(err) {
		return err
	}

	secretFilePath := pathWD + "/" + secretFileName

	input, err := ioutil.ReadFile(secretFilePath)
	if isError(err) {
		return err
	}

	for _, repo := range []string{"telar-web", "ts-serverless"} {
		repoPath := pathWD + "/" + repo + "/" + secretFileName
		err = ioutil.WriteFile(repoPath, input, 0644)
		if isError(err) {
			fmt.Println("Error creating", repoPath)
			fmt.Println(err)
			return err
		}
	}

	return nil
}

func createSecretFile(pathWD, name string, telarSecrets *TelarSecrets) error {

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
