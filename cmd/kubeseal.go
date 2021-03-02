package cmd

// Copyright (c) OpenFaaS Author(s). All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
// Refer to code: https://github.com/openfaas/faas-cli/blob/master/commands/cloud.go
// Refer to License: https://github.com/openfaas/faas-cli/blob/master/LICENSE

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Qolzam/telar-cli/pkg/log"
	"github.com/openfaas/faas-cli/proxy"
	"github.com/openfaas/faas-cli/schema"
)

func checkKubeseal() error {

	exist := commandExists("kubeseal")
	if !exist {
		err := downloadKubeSeal()
		if isError(err) {
			return fmt.Errorf("Check kubeseal %s", err.Error())
		}
	}
	return nil
}

func runCloudSeal(name string, pathWD string, args map[string]string, fromFile *[]string) error {
	certFile := pathWD + "/pub-cert.pem"
	outputFile := pathWD + "/secrets.yml"
	namespace := "openfaas-fn"
	if len(name) == 0 {
		return fmt.Errorf("--name is required")
	}

	log.Info("Sealing secret: %s in namespace: %s\n", name, namespace)

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

	log.Info("%s written.\n", outputFile)

	return nil
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

	log.Info("Starting download of kubeseal %s, this could take a few moments.\n", releaseVersion)
	output, err := downloadFile(http.DefaultClient, downloadURL, "kubeseal", downloadTo)

	if err != nil {
		return err
	}

	log.Info(`Download completed, please run:

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

func downloadFile(client *http.Client, url, name, downloadTo string) (string, error) {
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
