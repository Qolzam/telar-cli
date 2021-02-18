package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	browser "github.com/pkg/browser"
)

type ServerAction struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type RemoveFnPayload struct {
	ProjectDirectory string `json:"projectDirectory"`
}

type OpenURLModel struct {
	URL string `json:"url"`
}

var conn *websocket.Conn

func OpenURLHandler(w http.ResponseWriter, r *http.Request) {
	var input []byte

	if r.Body != nil {
		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body)

		input = body
	}

	var model OpenURLModel
	err := json.Unmarshal(input, &model)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Unmarshal body: %s", err.Error())))
	}
	browser.OpenURL(model.URL)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("input was: %s", string(input))))
}

func ClientHandler(w http.ResponseWriter, r *http.Request) {
	var input []byte

	if r.Body != nil {
		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body)

		input = body
	}

	var action ServerAction
	err := json.Unmarshal(input, &action)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Unmarshal body: %s", err.Error())))
	}

	switch action.Type {
	case START_STEP:
		go StartStep()
	case CHECK_STEP:
		var output ClientState
		cfg := &mapstructure.DecoderConfig{
			Metadata: nil,
			Result:   &output,
			TagName:  "json",
		}
		decoder, _ := mapstructure.NewDecoder(cfg)
		decoder.Decode(action.Payload)
		go checkStep(output)

	case REMOVE_SOCIAL_FROM_CLUSTER:
		var output RemoveFnPayload
		cfg := &mapstructure.DecoderConfig{
			Metadata: nil,
			Result:   &output,
			TagName:  "json",
		}
		decoder, _ := mapstructure.NewDecoder(cfg)
		decoder.Decode(action.Payload)
		go removeFunctionFromCluster(output)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("input was: %s", string(input))))
}

func checkStep(payload ClientState) {

	if payload.SetupStep != 0 {
		writeSetupCache(payload.Inputs.ProjectDirectory, payload.Inputs)
		fmt.Printf("Writing user inputs in setup cache...")
	}
	switch payload.SetupStep {
	case 0:
		CheckInitStep(payload.Inputs.ProjectDirectory)
	case 1:
		OFCAccessSetting()
	case 2:
		CheckIngredient(payload.Inputs.ProjectDirectory, TELAR_GITHUB_USER_NAME)
	case 3:
		CheckStorage(payload.Inputs.ProjectDirectory, payload.Inputs.BucketName)
	case 4:
		CheckDatabase(payload.Inputs.MongoDBHost, payload.Inputs.MongoDBPassword)
	case 5:
		CheckRecaptcha()
	case 6:
		CheckOAuth()
	case 7:
		CheckUserManagement(payload.Inputs.SocialDomain)
	case 8:
		CheckWebsocket(payload.Inputs)

	}
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not allowed", 403)
		return
	}
	var err error
	conn, err = websocket.Upgrade(w, r, w.Header(), 1024, 1024)

	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

}

func PingHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Pong")))

}
func Echo(message ClientAction) {

	if err := conn.WriteJSON(message); err != nil {
		fmt.Println(err)
	}

}
