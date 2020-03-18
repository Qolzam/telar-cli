package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/openfaas/faas-cli/proxy"
)

func pingWebsocket(url string) error {
	timeout := time.Second * 10
	client := proxy.MakeHTTPClient(&timeout, false)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	req, err := http.NewRequest(http.MethodHead, url+"/ping", nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	if res.StatusCode == 200 {
		return nil
	}
	return fmt.Errorf("Ping Websocket unknown error happend! Status Code: %d - URL: %s", res.StatusCode, url)
}
