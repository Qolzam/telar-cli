package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/websocket"
)

type Action struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

var conn *websocket.Conn

func ClientHandler(w http.ResponseWriter, r *http.Request) {
	var input []byte

	if r.Body != nil {
		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body)

		input = body
	}

	var action Action
	err := json.Unmarshal(input, &action)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Unmarshal body: %s", err.Error())))
	}
	switch action.Type {
	case START_STEP:
		go StartStep()

	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hello world, input was: %s", string(input))))
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

func Echo(message Action) {

	if err := conn.WriteJSON(message); err != nil {
		fmt.Println(err)
	}

}
