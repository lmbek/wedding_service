package frontend

import (
	"golang.org/x/net/websocket"
	"io"
	"log"
	"sync"
	"wedding_service/flags"
)

var (
	clients   = make(map[*websocket.Conn]bool) // connected clients
	broadcast = make(chan string)              // broadcast channel
	mutex     sync.Mutex
)

// Broadcast sends a message to the broadcast channel.
func Broadcast(message string) {
	broadcast <- message
}

// removeClient safely removes a client from the clients map and closes its connection.
func removeClient(ws *websocket.Conn) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, ok := clients[ws]; ok {
		delete(clients, ws)
		ws.Close()
		log.Println("Client disconnected:", ws.RemoteAddr())
	}
}

// HandleRegisterClient manages WebSocket connections.
// This is the starting point for the websocket connection
func HandleRegisterClient(ws *websocket.Conn) {
	if flags.LoadFrontendFlag() == "" {
		// only allow websocket to register clients and handle websocket logic
		// if file operations is used instead of embedded files
		return
	}

	mutex.Lock()
	log.Println("Client connected:", ws.RemoteAddr())
	clients[ws] = true
	mutex.Unlock()

	defer removeClient(ws)

	listenForInput(ws) // this is a blocking call
}

// listenForInput is a blocking function that holds the connection to handle websocket functionality
func listenForInput(ws *websocket.Conn) {
	for {
		var message string
		err := websocket.Message.Receive(ws, &message)
		if err != nil {
			if err == io.EOF {
				// Client disconnected cleanly, ignore
				break
			}
			log.Println("WebSocket receive error:", err)
			break // Exit loop on error, which will trigger deferred cleanup
		}
		log.Println("Received:", message)

		// Broadcast the received message
		Broadcast(message)
	}
}

func sendBroadcast() {
	mutex.Lock()
	defer mutex.Unlock()

	for client := range clients {
		if err := websocket.Message.Send(client, "build"); err != nil {
			log.Println("WebSocket send error:", err)
			// Use centralized removal for cleanup
			mutex.Unlock() // unlock before calling removeClient (which locks)
			removeClient(client)
			mutex.Lock() // re-lock for continuing the loop
		}
	}
}
