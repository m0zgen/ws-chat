package handlers

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode())

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WsJsonResponse response send back from websocket
type WsJsonResponse struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType int    `json:"message_type"`
}

func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Error upgrading connection", http.StatusInternalServerError)
		return
	}

	log.Println("Client connected")

	var response WsJsonResponse
	response.Message = `<em><small>Welcome to server</small></em>`

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println("Error writing JSON response:", err)
		return
	}

}

// Home - render Home page
func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}

}

func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		http.Error(w, "Template not found", http.StatusNotFound)
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return err
	}

	return nil

}
