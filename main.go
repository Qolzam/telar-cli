package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Qolzam/telar-cli/cmd"
	"github.com/gobuffalo/packr"
	"github.com/zserge/webview"
)

// Message : struct for message
type Message struct {
	Text string `json:"text"`
}

type ActionModel struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

func main() {
	// Bind folder path for packaging with Packr
	folder := packr.NewBox("./ui/build")

	// Handle to ./static/build folder on root path
	http.Handle("/", http.FileServer(folder))

	// Handle to showMessage func on /hello path
	http.HandleFunc("/ws", cmd.wsHandler)
	http.HandleFunc("/", cmd.rootHandler)
	// Run server at port 8080 as goroutine
	// for non-block working
	go http.ListenAndServe(":8080", nil)

	// Let's open window app with:
	//  - name: Golang App
	//  - address: http://localhost:8000
	//  - sizes: 800x600 px
	//  - resizable: true
	webview.Open("Telar", "http://localhost:8080", 800, 600, true)
}

func ofccSetup(w http.ResponseWriter, r *http.Request) {

	var input []byte

	if r.Body != nil {
		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body)

		input = body
	}

	var model ActionModel
	json.Unmarshal(input, &model)
	message := ActionModel{
		Type: "UNKNOW_ACTION",
	}
	switch model.Type {
	case "INIT_STEP":
		message = initSetp(model.Payload)
	default:
		fmt.Println("Unknow action")
	}

	// Return JSON encoding to output
	output, err := json.Marshal(message)

	// Catch error, if it happens
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set header Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Write output
	w.Write(output)
}

// Initialize step
func initSetp(payload map[string]interface{}) ActionModel {
	newPayload := make(map[string]interface{})

	directory := payload["directory"].(string)
	dirStatus := checkDirectory(directory)
	newPayload["dirStatus"] = dirStatus

	if !dirStatus {
		newPayload["error"] = fmt.Sprintf("Can not access to the directory ( %s ). Try running app in admin/root mode.", directory)
	} else {
		err := createTelarConfig(directory)
		if err != nil {
			newPayload["error"] = fmt.Sprintf("Creating config file %s", err.Error())
		}
	}

	return ActionModel{
		Type:    "PROJECT_DIRECTORY_STATUS",
		Payload: newPayload,
	}
}

// Check directory
func checkDirectory(dir string) bool {

	if err := os.Mkdir(dir, 0755); os.IsExist(err) {
		return true
	}

	return false
}

// Create Telar config file
func createTelarConfig(dir string) error {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile("telar.yml", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = f.Write([]byte("Hello"))
	if err != nil {
		return err
	}

	f.Close()
	return nil
}

func installGit() {

}
