package main

import (
	"fmt"
	"net/http"

	"github.com/Qolzam/telar-cli/cmd"
	"github.com/gobuffalo/packr/v2"
	browser "github.com/pkg/browser"
)

const (
	port = 31115
)

type Message struct {
	Text string `json:"text"`
}

type ActionModel struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

func main() {
	// Bind folder path for packaging with Packr
	folder := packr.New("telar-pack", "./ui/build")
	fmt.Println("[INFO] UI packed.", folder.Path)

	// Handle to ./static/build folder on root path
	http.Handle("/", http.FileServer(folder))

	// Other handlers
	http.HandleFunc("/ws", cmd.WsHandler)
	http.HandleFunc("/dispatch", cmd.ClientHandler)
	http.HandleFunc("/open-url", cmd.OpenURLHandler)
	http.HandleFunc("/ping", cmd.PingHandler)

	// Run server as goroutine
	// for non-block working
	go browser.OpenURL(fmt.Sprintf("http://localhost:%d", port))
	fmt.Printf("\n[INFO] Server started on http://localhost:%d", port)

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

}
