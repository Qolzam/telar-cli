package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/openfaas/faas-cli/commands"
	"github.com/openfaas/faas-cli/proxy"
)

const (
	Source              = "garbage-collect"
	namespace           = ""
	FunctionLabelPrefix = "telar.dev."
)

var timeout = 3 * time.Second

// Handle function cleans up functions which were removed or renamed
// within the repo for the given user.
func RemoveOFFunctions(garbageReq GarbageRequest) (string, error) {

	appID := garbageReq.AppID
	if garbageReq.Repo == "*" {
		log.Printf("Removing all functions for %s", appID)
	}

	deployedFunctions, err := listFunctions(garbageReq.AppID, garbageReq.Namespace, garbageReq.Gateway, garbageReq.Token, garbageReq.TLSInsecure)

	if err != nil {
		log.Fatal(err)
	}

	deployedList := ""
	for _, fn := range deployedFunctions {
		deployedList += fn.GetRepo() + "/" + fn.Name + ", "
	}

	log.Printf("Functions with app ID %s:\n %s", appID, strings.Trim(deployedList, ", "))

	proxyClient, err := GetOFProxyClient(garbageReq.Gateway, garbageReq.Token, garbageReq.TLSInsecure)
	if err != nil {
		return "", err
	}

	deleted := 0
	for _, fn := range deployedFunctions {
		if garbageReq.Repo == "*" ||
			(fn.GetRepo() == garbageReq.Repo && included(&fn, appID, garbageReq.Functions)) {
			log.Printf("Delete: %s\n", fn.Name)
			err = proxyClient.DeleteFunction(context.Background(), fn.Name, namespace)
			if err != nil {
				echoDialogInfo(fmt.Sprintf("Unable to delete function: `%s`", fn.Name), "")
				log.Println(err)
			}
			deleted = deleted + 1
		}
	}

	echoDialogInfo(fmt.Sprintf("Garbage collection ran for %s/%s - %d functions deleted.", garbageReq.AppID, garbageReq.Repo, deleted), "")

	return fmt.Sprintf("Garbage collection ran for %s/%s - %d functions deleted.", garbageReq.AppID, garbageReq.Repo, deleted), nil
}

func included(fn *openFaaSFunction, appID string, functionStack []string) bool {

	for _, name := range functionStack {
		if strings.EqualFold(name, fn.Name) && fn.GetAppID() == appID {
			return true
		}
	}

	return false
}

func listFunctions(appID, namespace, gateway, token string, tlsInsecure bool) ([]openFaaSFunction, error) {
	proxyClient, err := GetOFProxyClient(gateway, token, tlsInsecure)
	if err != nil {
		return nil, err
	}

	functions, err := proxyClient.ListFunctions(context.Background(), namespace)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%-30s\t%-"+"s\t%-15s\t%-5s\n", "Function", "Image", "Invocations", "Replicas")
	var parsedFunctionList []openFaaSFunction
	for _, function := range functions {

		functionImage := function.Image
		// if len(function.Image) > 40 {
		// 	functionImage = functionImage[0:38] + ".."
		// }
		openFaasFn := openFaaSFunction{Name: function.Name, Image: function.Image, Labels: *function.Labels}
		parsedFunctionList = append(parsedFunctionList, openFaasFn)
		fmt.Printf("%-30s\t%-"+"s\t%-15d\t%-5d\n", function.Name, functionImage, int64(function.InvocationCount), function.Replicas)
	}

	return parsedFunctionList, nil
}

type GarbageRequest struct {
	Functions   []string `json:"functions"`
	Repo        string   `json:"repo"`
	Namespace   string   `json:"namespace"`
	AppID       string   `json:"appID"`
	Token       string   `json:"token"`
	Gateway     string   `json:"gateway"`
	TLSInsecure bool     `json:"tlsInsecure"`
}

type openFaaSFunction struct {
	Name   string            `json:"name"`
	Image  string            `json:"image"`
	Labels map[string]string `json:"labels"`
}

func (f *openFaaSFunction) GetAppID() string {
	return f.Labels[FunctionLabelPrefix+"appID"]
}

func (f *openFaaSFunction) GetRepo() string {
	return f.Labels[FunctionLabelPrefix+"repo"]
}

func GetOFProxyClient(gateway, token string, tlsInsecure bool) (*proxy.Client, error) {
	var gatewayAddress string
	gatewayAddress = GetGatewayURL(gateway, "", "", os.Getenv(openFaaSURLEnvironment))

	cliAuth, err := proxy.NewCLIAuth(token, gatewayAddress)
	if err != nil {
		return nil, err
	}
	transport := commands.GetDefaultCLITransport(tlsInsecure, &ofCommandTimeout)
	proxyClient, err := proxy.NewClient(cliAuth, gatewayAddress, transport, &ofCommandTimeout)
	if err != nil {
		return nil, err
	}

	return proxyClient, nil
}
