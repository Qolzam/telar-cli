package cmd

import (
	"fmt"

	"github.com/Qolzam/telar-cli/pkg/log"
)

func removeFunctionFromCluster(payload RemoveFnPayload) {
	log.Info(" Start removing functions from cluster ...")
	setupCache, err := getSetupYaml(payload.ProjectDirectory)
	if err != nil {
		echoDialogInfo(fmt.Sprintf("Can not get `%s` file in [%s]. %s", SETUP_YAML_FILE_NAME, payload.ProjectDirectory, err.Error()), "")
		return
	}

	openfaasPass, err := getOpenFaasPass()
	if isError(err) {
		echoDialogInfo(err.Error(), "")
		log.Error(err.Error())
		return
	}

	authConf, _, err := runFaaSLogin(setupCache.ClientInputs.OFGateway, setupCache.ClientInputs.OFUsername, openfaasPass, false)
	if isError(err) {
		echoDialogInfo(err.Error(), "")
		log.Error(err.Error())
		return
	}
	garbageRequest := GarbageRequest{
		Namespace:   setupCache.ClientInputs.Namespace,
		AppID:       setupCache.ClientInputs.AppID,
		Token:       authConf.Token,
		Gateway:     setupCache.ClientInputs.OFGateway,
		TLSInsecure: false,
	}
	telarWebFns := []string{"admin", "auth", "profile", "setting", "storage", "actions", "notifications"}
	garbageRequest.Functions = telarWebFns
	garbageRequest.Repo = "telar-web"
	RemoveOFFunctions(garbageRequest)

	tsServerlessFns := []string{"posts", "media", "comments", "votes", "circles", "user-rels"}
	garbageRequest.Functions = tsServerlessFns
	garbageRequest.Repo = "ts-serverless"
	RemoveOFFunctions(garbageRequest)

	tsUIFns := []string{"web"}
	garbageRequest.Functions = tsUIFns
	garbageRequest.Repo = "ts-ui"
	RemoveOFFunctions(garbageRequest)

}
