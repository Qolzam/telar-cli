package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func checkCustomers(customer string) (string, error) {
	url := "https://raw.githubusercontent.com/openfaas/openfaas-cloud/master/CUSTOMERS"
	client := http.DefaultClient
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

	if res.Body != nil {
		defer res.Body.Close()
		res, _ := ioutil.ReadAll(res.Body)
		customers := string(res)
		fmt.Println(strings.Contains(customers, customer))
		return string(res), nil
	}
	return "", fmt.Errorf("error downloading %s", url)
}

func getOFCCGateway(githubUsername string) string {
	url := fmt.Sprintf("https://%s.o6s.io", githubUsername)
	return url
}

func StartStep() {
	payloadInfoDialog := struct {
		Message string `json:"message"`
		URL     string `json:"url"`
	}{
		Message: "",
		URL:     "",
	}
	action := Action{
		Type:    SHOW_INFO_DIALOG,
		Payload: payloadInfoDialog,
	}
	projectPath, err := getDefaultProjectDirectory()
	if err != nil {
		Echo(action)
		return
	}
	action.Type = SET_INPUT
	action.Payload = struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}{
		Key:   "projectDirectory",
		Value: projectPath,
	}
	Echo(action)
}
