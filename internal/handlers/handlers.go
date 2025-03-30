package handlers

import (
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sort"
)

var wsChan = make(chan WsPayload)
var clients = make(map[WebSocketConnection]string)

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

type WebSocketConnection struct {
	*websocket.Conn
}

// WsJsonResponse response send back from websocket
type WsJsonResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    int      `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}

type WsPayload struct {
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"`
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

	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println("Error writing JSON response:", err)
		return
	}

	go ListenForWs(&conn)

}

func ListenToWs() {
	var response WsJsonResponse

	for {
		e := <-wsChan

		switch e.Action {
		case "username":
			// get a list of all users and sen it back to broadcast
			clients[e.Conn] = e.Username
			users := getUserList()
			response.Action = "list_users"
			response.ConnectedUsers = users
			broadcastToAll(response)

		case "left":
			response.Action = "list_users"
			delete(clients, e.Conn)

			users := getUserList()
			response.ConnectedUsers = users
			broadcastToAll(response)

		case "broadcast":
			response.Action = "broadcast"
			response.Message = fmt.Sprintf("<b>%s</b>: %s", e.Username, e.Message)
			broadcastToAll(response)
		}

		//response.Action = "Got here"
		//response.Message = fmt.Sprintf("Some message, and action was %s", e.Action)
		//broadcastToAll(response)
	}
}

func getUserList() []string {
	var users []string
	for _, client := range clients {
		users = append(users, client)
	}
	sort.Strings(users)
	return users
}

func broadcastToAll(response WsJsonResponse) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("Error writing JSON to client:", err)
			_ = client.Close()
			delete(clients, client)
		} else {
			log.Println("Sent message to client")
		}
	}
}

func ListenForWs(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in ListenForWs:", r)
		}
	}()

	var payload WsPayload
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			//log.Println("Error reading JSON:", err)
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
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
