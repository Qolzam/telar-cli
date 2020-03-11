package cmd

func createSecretFile(pathWD, githubUsername string, telarSecrets *TelarSecrets) error {
	name := githubUsername + "-secrets"
	args := make(map[string]string)
	args["mongo-pwd"] = telarSecrets.MongoPwd
	args["recaptcha-key"] = telarSecrets.RecaptchaKey
	args["ts-client-secret"] = telarSecrets.TsClientSecret
	args["redis-pwd"] = telarSecrets.RedisPwd
	args["admin-username"] = telarSecrets.AdminUsername
	args["admin-password"] = telarSecrets.AdminPwd
	args["payload-secret"] = telarSecrets.RecaptchaKey
	args["ref-email-pass"] = telarSecrets.RefEmailPwd
	args["phone-auth-token"] = telarSecrets.PhoneAuthToken
	args["phone-auth-id"] = telarSecrets.PhoneAuthId
	saPath := pathWD + "/serviceAccountKey.json"
	publicKeyPath := pathWD + "/key.pub"
	privateKeyPath := pathWD + "/key"
	files := []string{saPath, publicKeyPath, privateKeyPath}
	return runCloudSeal(name, pathWD, args, &files)
}
