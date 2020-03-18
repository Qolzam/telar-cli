package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"context"

	firebase "firebase.google.com/go"
	execute "github.com/alexellis/go-execute/pkg/v1"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/copier"
	"github.com/openfaas/faas-cli/proxy"
	"github.com/openfaas/faas-cli/schema"
	stack "github.com/openfaas/faas-cli/stack"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/api/option"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
	"gopkg.in/yaml.v2"
)

var (
	downloadVersion string
	downloadTo      string
)

const (
	OFCC_DOMAIN = ".o6s.io"
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
	MongoDBHost      string `json:"mongoDBHost"`
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

type msg struct {
	Num int
}

func main() {
	// Bind folder path for packaging with Packr
	// checkSudo()
	path, err := getDefaultProjectDirectory()
	if err != nil {
		log.Fatal(err)
	}
	// checkKubeseal()
	// args := make(map[string]string)
	// args["mongo-pwd"] = "$MONGO_PWD"
	// args["recaptcha-key"] = "$RECAPTCHA_KEY"
	// saPath := prDir + "/serviceAccountKey.json"
	// files := []string{saPath}
	// err = runCloudSeal("red-gold", prDir, args, &files)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// checkCustomers("Qolzam")
	// err = cloneTSUI(prDir)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = cloneTSServerless(prDir)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = cloneTelarWeb(prDir)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = checkFirebaseServiceAccountExist(prDir)
	// if isError(err) {
	// 	log.Fatal(err)
	// }
	// err = checkFirebaseStorageBucket(prDir, "resume-web-app")
	// if isError(err) {
	// 	log.Fatal(err)
	// }
	// ok, err := checkWebsocket("https://red-gold-socket.herokuapp.com")
	// if isError(err) {
	// 	log.Fatal(err)
	// }
	// if ok {
	// 	log.Println("Webdocket is ok")
	// }
	// err := checkDB()
	// if isError(err) {
	// 	log.Fatal(err)
	// }
	// err = applyAppConfig(prDir)
	// if isError(err) {
	// 	log.Fatal(err)
	// }
	// err = applyConfigTSUI(prDir, "wss://some.heroku-app.com")
	// if isError(err) {
	// 	log.Fatal(err)
	// }
	// telarConfig := TelarConfig{
	// 	GithubUsername:   "Qolzam",
	// 	PathWD:           prDir,
	// 	CoockieDomain:    ".io4.com",
	// 	Bucket:           "resume-test",
	// 	ClientID:         "34545",
	// 	URL:              "https://Qolzam.io4.com",
	// 	WebsocketURL:     "https://Qolzam.heroku.com",
	// 	MongoDBHost:        "mongoDBHost",
	// 	MongoDatabase:    "mongoDB",
	// 	RecaptchaSiteKey: "re-si-key",
	// 	RefEmail:         "ref@email.com",
	// }
	// err = applyConfig(telarConfig)
	// if isError(err) {
	// 	log.Fatal(err)
	// }
	// err = createTSUIStack(prDir, "https://websocket.heroku.com")
	// if isError(err) {
	// 	log.Fatal(err)
	// }

	// err = gitCommit(prDir + "/ts-ui")
	// if isError(err) {
	// 	log.Fatal(err)
	// }

	// err = gitPush(prDir + "/ts-ui")
	// if isError(err) {
	// 	log.Fatal(err)
	// }

	// http.HandleFunc("/ws", wsHandler)
	// http.HandleFunc("/", rootHandler)

	// panic(http.ListenAndServe(":8080", nil))
	// _, err = generatePayloadSecret()
	// if isError(err) {
	// 	log.Fatal(err)
	// }
	err = createPrivateKey(path)
	if isError(err) {
		log.Fatal(err)
	}
	err = createPublicKey(path)
	if isError(err) {
		log.Fatal(err)
	}
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
	fmt.Println(payloadSecret)
	return payloadSecret, nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadFile("index.html")
	if err != nil {
		fmt.Println("Could not open file.", err)
	}
	fmt.Fprintf(w, "%s", content)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not allowed", 403)
		return
	}
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	go echo(conn)
}
func echo(conn *websocket.Conn) {
	for {
		m := msg{}

		err := conn.ReadJSON(&m)
		if err != nil {
			fmt.Println("Error reading json.", err)
		}

		fmt.Printf("Got message: %#v\n", m)

		if err = conn.WriteJSON(m); err != nil {
			fmt.Println(err)
		}
	}
}

// as util
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func checkKubeseal() {

	exist := commandExists("kubeseal")
	if exist {
		fmt.Println("exist")
	} else {
		downloadKubeSeal()
	}
}

func applyConfig(telarConfig TelarConfig) error {
	for _, repo := range []string{"telar-web", "ts-serverless"} {
		err := applyAppConfig(telarConfig.PathWD+"/"+repo, telarConfig.MongoDBHost, telarConfig.MongoDatabase, telarConfig.RecaptchaSiteKey, telarConfig.RefEmail)
		if isError(err) {
			return err
		}
	}
	return nil
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

func applyAppConfig(repoPath, mongoDBHost, mongoDatabase, recaptchaSiteKey, refEmail string) error {
	filePath := "config/app_config.yml"
	envs, err := readConfigFile(repoPath, filePath)
	if isError(err) {
		return err
	}

	envs["mongo_user"] = mongoDBHost
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

func getOFCCGateway(githubUsername string) string {
	url := fmt.Sprintf("https://%s.o6s.io", githubUsername)
	return url
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

func runCloudSeal(name string, pathWD string, args map[string]string, fromFile *[]string) error {
	certFile := pathWD + "/pub-cert.pem"
	outputFile := pathWD + "/secrets.yml"
	namespace := "openfaas-fn"
	if len(name) == 0 {
		return fmt.Errorf("--name is required")
	}

	fmt.Printf("Sealing secret: %s in namespace: %s\n", name, namespace)

	fmt.Println("")

	enc := base64.StdEncoding

	secret := schema.KubernetesSecret{
		ApiVersion: "v1",
		Kind:       "Secret",
		Metadata: schema.KubernetesSecretMetadata{
			Name:      name,
			Namespace: namespace,
		},
		Data: make(map[string]string),
	}

	for k, v := range args {
		secret.Data[k] = enc.EncodeToString([]byte(v))
	}

	if fromFile != nil {
		for _, file := range *fromFile {
			bytesOut, err := ioutil.ReadFile(file)
			if err != nil {
				return err
			}

			key := filepath.Base(file)
			secret.Data[key] = enc.EncodeToString(bytesOut)
		}
	}

	sec, err := json.Marshal(secret)
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(certFile); err != nil {
		return fmt.Errorf("unable to load public certificate %s", certFile)
	}

	kubeseal := exec.Command("kubeseal", "--format=yaml", "--cert="+certFile)

	stdin, stdinErr := kubeseal.StdinPipe()
	if stdinErr != nil {
		panic(stdinErr)
	}

	stdin.Write(sec)
	stdin.Close()

	out, err := kubeseal.CombinedOutput()
	if err != nil {
		return fmt.Errorf("unable to start \"kubeseal\", check it is installed, error: %s", err.Error())
	}

	writeErr := ioutil.WriteFile(outputFile, out, 0755)

	if writeErr != nil {
		return fmt.Errorf("unable to write secret: %s to %s", name, outputFile)
	}

	fmt.Printf("%s written.\n", outputFile)

	return nil
}

func getDefaultProjectDirectory() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/telar-social", usr.HomeDir), nil

}

// Check directory
func checkDirectory(dir string) bool {

	if err := os.Mkdir(dir, 0755); os.IsExist(err) {
		return true
	}

	return false
}

func downloadKubeSeal() error {
	releases := "https://github.com/bitnami-labs/sealed-secrets/releases/latest"

	releaseVersion := downloadVersion
	if len(downloadVersion) == 0 {
		version, err := findRelease(releases)
		if err != nil {
			return err
		}
		releaseVersion = version
	}

	osVal := runtime.GOOS
	arch := runtime.GOARCH

	if arch == "x86_64" {
		arch = "amd64"
	}

	downloadURL := "https://github.com/bitnami/sealed-secrets/releases/download/" + releaseVersion + "/kubeseal-" + osVal + "-" + arch

	fmt.Printf("Starting download of kubeseal %s, this could take a few moments.\n", releaseVersion)
	output, err := downloadBinary(http.DefaultClient, downloadURL, "kubeseal", downloadTo)

	if err != nil {
		return err
	}

	fmt.Printf(`Download completed, please run:

  chmod +x %s
  %s --version
  sudo install %s /usr/local/bin/

  `, output, output, output)

	err = chmod(output)

	if err != nil {
		return err
	}

	err = install(output)
	if err != nil {
		return err
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

// Copyright (c) OpenFaaS Author(s). All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
// Refer to code: https://github.com/openfaas/faas-cli/blob/master/commands/cloud.go
// Refer to License: https://github.com/openfaas/faas-cli/blob/master/LICENSE
//
func findRelease(url string) (string, error) {
	timeout := time.Second * 5
	client := proxy.MakeHTTPClient(&timeout, false)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	if res.StatusCode != 302 {
		return "", fmt.Errorf("incorrect status code: %d", res.StatusCode)
	}

	loc := res.Header.Get("Location")
	if len(loc) == 0 {
		return "", fmt.Errorf("unable to determine release of kubeseal")
	}
	version := loc[strings.LastIndex(loc, "/")+1:]
	return version, nil
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
		return fmt.Errorf("HTTP error status code %s", res.StatusCode)
	}

	if res.Body != nil {
		defer res.Body.Close()
		res, _ := ioutil.ReadAll(res.Body)
		customers := string(res)
		if strings.Contains(customers, customer) {
			return nil
		} else {
			return fmt.Errorf("The %s user is not registered in OpenFaaS Cloud CUSTOMERS https://raw.githubusercontent.com/openfaas/openfaas-cloud/master/CUSTOMERS", customers)
		}
	}
	return fmt.Errorf("Unkown error when checking %s for OpenFaaS Cloud CUSTOMERS ", customer)
}

func downloadBinary(client *http.Client, url, name, downloadTo string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	res, err := client.Do(req)

	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("could not find release, http status code was %d, release may not exist for this architecture", res.StatusCode)
	}

	var tempDir string
	if len(downloadTo) == 0 {
		tempDir = os.TempDir()
	} else {
		tempDir = downloadTo
	}

	outputPath := path.Join(tempDir, name)
	if res.Body != nil {
		defer res.Body.Close()
		res, _ := ioutil.ReadAll(res.Body)

		err := ioutil.WriteFile(outputPath, res, 0600)
		if err != nil {
			return "", err
		}
		return outputPath, nil
	}
	return "", fmt.Errorf("error downloading %s", url)
}

func cloneTSUI(rootDir string) error {
	return gitClone(rootDir+"/ts-ui", "git@github.com:red-gold/ts-ui.git")
}

func cloneTSServerless(rootDir string) error {
	return gitClone(rootDir+"/ts-serverless", "git@github.com:red-gold/ts-serverless.git")
}

func cloneTelarWeb(rootDir string) error {
	return gitClone(rootDir+"/telar-web", "git@github.com:red-gold/telar-web.git")
}

func checkFirebaseServiceAccountExist(pathWD string) error {
	var file, err = os.OpenFile(pathWD+"/serviceAccountKey.json", os.O_RDONLY, 0444)
	if isError(err) {
		return err
	}
	defer file.Close()
	return nil
}

func gitClone(projectDirectory string, repoURL string) error {
	currentUser, err := user.Current()
	if err != nil {
		return err
	}
	sshAuth, err := ssh.NewPublicKeysFromFile("git", currentUser.HomeDir+"/.ssh/id_rsa", "")
	if err != nil {
		return err
	}
	_, err = git.PlainClone(projectDirectory, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
		Auth:     sshAuth,
	})

	if err != nil {
		return err
	}
	return nil
}

func isError(err error) bool {
	return err != nil
}

func checkFirebaseStorageBucket(pathWD string, bucketName string) error {
	ctx := context.Background()
	config := &firebase.Config{
		StorageBucket: bucketName + ".appspot.com",
	}
	opt := option.WithCredentialsFile(pathWD + "/serviceAccountKey.json")
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return err
	}

	client, err := app.Storage(ctx)
	if err != nil {
		return err
	}
	bucket, err := client.DefaultBucket()
	if err != nil {
		return err
	}
	r := bytes.NewReader([]byte("Test firebase."))
	wc := bucket.Object("telar.test").NewWriter(ctx)
	if _, err = io.Copy(wc, r); err != nil {
		fmt.Println(err.Error())
	}
	if err := wc.Close(); err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func checkWebsocket(url string) (bool, error) {
	timeout := time.Second * 5
	client := proxy.MakeHTTPClient(&timeout, false)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	req, err := http.NewRequest(http.MethodHead, url+"/ping", nil)
	if err != nil {
		return false, err
	}

	res, err := client.Do(req)
	if err != nil {
		return false, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	if res.StatusCode == 200 {
		return true, nil
	}
	return false, nil
}

func checkDB() error {
	dbURL := "mongodb+srv://telar_user:pass@cluster0-l6ojz.mongodb.net/test?retryWrites=true&w=majority"
	client, err := mongo.NewClient(options.Client().ApplyURI(fmt.Sprintf(dbURL)))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if isError(err) {
		return err
	}
	return nil
}

func gitCommit(path string) error {
	repo, err := git.PlainOpen(path)
	if isError(err) {
		return err
	}
	w, err := repo.Worktree()
	if isError(err) {
		return err
	}

	_, err = w.Add(".")
	if isError(err) {
		return err
	}

	fmt.Println(path)
	status, err := w.Status()

	fmt.Println(status)

	commit, err := w.Commit("[TELAR] Initialize Telar Social.", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Telar",
			Email: "support@red-gold.tech",
			When:  time.Now(),
		},
	})
	if isError(err) {
		return err
	}

	obj, err := repo.CommitObject(commit)
	if isError(err) {
		return err
	}

	fmt.Println(obj)
	return nil
}

func gitPush(path string) error {
	currentUser, err := user.Current()
	if err != nil {
		return err
	}
	repo, err := git.PlainOpen(path)
	if isError(err) {
		return err
	}

	sshAuth, err := ssh.NewPublicKeysFromFile("git", currentUser.HomeDir+"/.ssh/id_rsa", "")
	if err != nil {
		return err
	}
	err = repo.Push(&git.PushOptions{
		Auth:     sshAuth,
		Progress: os.Stdout,
	})
	if isError(err) {
		return err
	}
	return nil
}
