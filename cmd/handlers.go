package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"
	browser "github.com/pkg/browser"
)

type ServerAction struct {
	Type    string      `json:"type"`
	Payload ClientState `json:"payload"`
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
		go checkStep(action.Payload)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("input was: %s", string(input))))
}

func checkStep(payload ClientState) {

	switch payload.SetupStep {
	case 0:
		CheckInitStep(payload.Inputs.ProjectDirectory)
	case 1:
		CheckIngredient(payload.Inputs.ProjectDirectory, payload.Inputs.GithubUsername)
	case 2:
		CheckStorage(payload.Inputs.ProjectDirectory, payload.Inputs.BucketName)
	case 3:
		CheckDatabase(payload.Inputs.MongoDBHost, payload.Inputs.MongoDBPassword)
	case 4:
		CheckRecaptcha()
	case 5:
		CheckOAuth()
	case 6:
		CheckUserManagement(payload.Inputs.GithubUsername)
	case 7:
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
