package controllers

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Define a connection pool to manage all WebSocket clients
var (
	upgrader    = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	clients     = make(map[*websocket.Conn]bool)
	clientsLock sync.Mutex
)

// HandleWebSocket upgrades the connection and manages WebSocket communication
func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// Add the connection to the pool
	registerClient(conn)
	defer unregisterClient(conn)

	log.Println("New WebSocket connection established")

	// Handle incoming messages
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		// Broadcast the message to all clients
		broadcastMessage(messageType, message)
	}
}

// registerClient adds a WebSocket connection to the pool
func registerClient(conn *websocket.Conn) {
	clientsLock.Lock()
	defer clientsLock.Unlock()
	clients[conn] = true
}

// unregisterClient removes a WebSocket connection from the pool
func unregisterClient(conn *websocket.Conn) {
	clientsLock.Lock()
	defer clientsLock.Unlock()
	delete(clients, conn)
}

// broadcastMessage sends a message to all connected WebSocket clients
func broadcastMessage(messageType int, message []byte) {
	clientsLock.Lock()
	defer clientsLock.Unlock()

	for conn := range clients {
		err := conn.WriteMessage(messageType, message)
		if err != nil {
			log.Printf("Error broadcasting message: %v", err)
			conn.Close()
			delete(clients, conn)
		}
	}
}
