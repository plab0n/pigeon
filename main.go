package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}
	//Note: The client is not assigned to any hub yet
	client := &Client{conn: conn, send: make(chan []byte, 256)}
	fmt.Println("Client connected")
	go client.readPump()
	go client.writePump()
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)

	fmt.Println("WebSocket server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
