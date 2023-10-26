package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/plab0n/pigeon/server"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var hubs map[string]*server.Hub

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}
	//Note: The client is not assigned to any hub yet
	client := server.CreateClient(conn)
	hubId := r.URL.Query().Get("roomId")
	if len(hubId) > 0 {
		addClientToHub(client, hubId)
	}
	fmt.Println("Client connected")
	go client.ReadPump()
	go client.WritePump()
}
func addClientToHub(client *server.Client, hubId string) {
	log.Println("Adding client to hub. HubId: ", hubId)
	hub, ok := hubs[hubId]
	if !ok {
		hub = server.CreateHub()
		hubs[hubId] = hub
		go hub.Run()
	}
	hub.Register(client)
}
func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}
func sendMessage(w http.ResponseWriter, r *http.Request) {
	hubId := r.URL.Query().Get("roomId")
	message := r.URL.Query().Get("message")
	hub, ok := hubs[hubId]
	if ok {
		log.Println("Boradcasting: ", message)
		hub.Broadcast([]byte(message))
	}
}
func main() {
	hubs = make(map[string]*server.Hub)
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/ws/{roomId}", handleWebSocket)
	http.HandleFunc("/send", sendMessage)
	fmt.Println("WebSocket server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
