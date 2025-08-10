package frontend

import (
	"golang.org/x/net/websocket"
	"io"
	"log/slog"
	"sync"
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
		slog.Info("Client disconnected", slog.String("remote_addr", ws.RemoteAddr().String()))
	}
}

// HandleRegisterClient manages WebSocket connections.
// This is the starting point for the websocket connection
func HandleRegisterClient(ws *websocket.Conn, frontendPath string, hotReloadEnabled bool) {
	if hotReloadEnabled {
		if frontendPath == "" {
			slog.Error("Frontend path must be set when hot reload is enabled")
			return
		}

		addClientToClientList(ws) // add client to the global client list
		defer removeClient(ws)    // remove client when it is done listening
		listenForInput(ws)        // this is a blocking call
		return
	}
	return
}

func addClientToClientList(ws *websocket.Conn) {
	mutex.Lock()
	slog.Info("Client connected", slog.String("remote_addr", ws.RemoteAddr().String()))
	clients[ws] = true
	mutex.Unlock()
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
			slog.Error("WebSocket receive error", slog.Any("error", err))
			break // Exit loop on error, which will trigger deferred cleanup
		}
		slog.Info("Received message", slog.String("message", message))

		// Broadcast the received message
		Broadcast(message)
	}
}

func sendBroadcast() {
	slog.Info("File change detected, broadcasting to clients")
	mutex.Lock()
	defer mutex.Unlock()

	for client := range clients {
		if err := websocket.Message.Send(client, "build"); err != nil {
			slog.Error("WebSocket send error", slog.Any("error", err))
			// Use centralized removal for cleanup
			mutex.Unlock() // unlock before calling removeClient (which locks)
			removeClient(client)
			mutex.Lock() // re-lock for continuing the loop
		}
	}
}
