package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/openfaas/faas-cli/proxy"

	"github.com/openfaas/faas-cli/config"
)

var (
	username      string
	password      string
	passwordStdin bool
)

const openFaaSURLEnvironment = "OPENFAAS_URL"

func runFaaSLogin(gateway, username, password string, tlsInsecure bool) (*config.AuthConfig, *string, error) {
	timeout := time.Second * 5
	gateway = GetGatewayURL(gateway, "", "", os.Getenv(openFaaSURLEnvironment))
	if len(username) == 0 {
		return nil, nil, fmt.Errorf("must provide --username or -u")
	}

	if len(password) > 0 {
		fmt.Println("WARNING! Using --password is insecure, consider using: cat ~/faas_pass.txt | faas-cli login -u user --password-stdin")
		if passwordStdin {
			return nil, nil, fmt.Errorf("--password and --password-stdin are mutually exclusive")
		}

		if len(username) == 0 {
			return nil, nil, fmt.Errorf("must provide --username with --password")
		}
	}

	if passwordStdin {
		if len(username) == 0 {
			return nil, nil, fmt.Errorf("must provide --username with --password-stdin")
		}

		passwordStdin, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return nil, nil, err
		}

		password = strings.TrimSpace(string(passwordStdin))
	}

	password = strings.TrimSpace(password)
	if len(password) == 0 {
		return nil, nil, fmt.Errorf("must provide a non-empty password via --password or --password-stdin")
	}

	fmt.Println("Calling the OpenFaaS server to validate the credentials...")

	if err := validateFaasLogin(gateway, username, password, timeout, tlsInsecure); err != nil {
		return nil, nil, err
	}

	token := config.EncodeAuth(username, password)
	if err := config.UpdateAuthConfig(gateway, token, config.BasicAuthType); err != nil {
		return nil, nil, err
	}

	authConfig, err := config.LookupAuthConfig(gateway)
	if err != nil {
		return nil, nil, err
	}

	user, _, err := config.DecodeAuth(authConfig.Token)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("credentials saved for", user, gateway)

	return &authConfig, &user, nil
}

func validateFaasLogin(gatewayURL string, user string, pass string, timeout time.Duration, insecureTLS bool) error {

	if len(checkTLSInsecure(gatewayURL, insecureTLS)) > 0 {
		fmt.Printf(NoTLSWarn)
	}

	client := proxy.MakeHTTPClient(&timeout, insecureTLS)
	req, err := http.NewRequest("GET", gatewayURL+"/system/functions", nil)
	if err != nil {
		return fmt.Errorf("invalid URL: %s", gatewayURL)
	}

	req.SetBasicAuth(user, pass)
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("cannot connect to OpenFaaS on URL: %s. %v", gatewayURL, err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	switch res.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized:
		return fmt.Errorf("unable to login, either username or password is incorrect")
	default:
		bytesOut, err := ioutil.ReadAll(res.Body)
		if err == nil {
			return fmt.Errorf("server returned unexpected status code: %d - %s", res.StatusCode, string(bytesOut))
		}
	}

	return nil
}

const (
	// NoTLSWarn Warning thrown when no SSL/TLS is used
	NoTLSWarn = "WARNING! You are not using an encrypted connection to the gateway, consider using HTTPS."
)

// checkTLSInsecure returns a warning message if the given gateway does not have https.
// Use tsInsecure to skip validations
func checkTLSInsecure(gateway string, tlsInsecure bool) string {
	if !tlsInsecure {
		if strings.HasPrefix(gateway, "https") == false &&
			strings.HasPrefix(gateway, "http://127.0.0.1") == false &&
			strings.HasPrefix(gateway, "http://localhost") == false {
			return NoTLSWarn
		}
	}
	return ""
}

func GetGatewayURL(argumentURL, defaultURL, yamlURL, environmentURL string) string {
	var gatewayURL string

	if len(argumentURL) > 0 && argumentURL != defaultURL {
		gatewayURL = argumentURL
	} else if len(yamlURL) > 0 && yamlURL != defaultURL {
		gatewayURL = yamlURL
	} else if len(environmentURL) > 0 {
		gatewayURL = environmentURL
	} else {
		gatewayURL = defaultURL
	}

	gatewayURL = strings.ToLower(strings.TrimRight(gatewayURL, "/"))
	if !strings.HasPrefix(gatewayURL, "http") {
		gatewayURL = fmt.Sprintf("http://%s", gatewayURL)
	}

	return gatewayURL
}
