package cmd

import (
	"net/http"
	"time"

	"github.com/openfaas/faas-cli/proxy"
)

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
