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
    clients     = make(map[string]map[*websocket.Conn]bool) // Map of room ID to clients
    clientsLock sync.Mutex
)

// HandleWebSocket upgrades the connection and manages WebSocket communication
func HandleWebSocket(c *gin.Context) {
    roomID := c.Param("id")
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Printf("Failed to upgrade to WebSocket: %v", err)
        return
    }
    defer conn.Close()

    // Add the connection to the pool
    registerClient(roomID, conn)
    defer unregisterClient(roomID, conn)

    log.Printf("New WebSocket connection established in room %s", roomID)

    // Handle incoming messages
    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            log.Printf("Error reading message: %v", err)
            break
        }

        // Broadcast the message to all clients in the same room
        broadcastMessage(roomID, messageType, message)
    }
}

// registerClient adds a WebSocket connection to the pool for a specific room
func registerClient(roomID string, conn *websocket.Conn) {
    clientsLock.Lock()
    defer clientsLock.Unlock()
    if clients[roomID] == nil {
        clients[roomID] = make(map[*websocket.Conn]bool)
    }
    clients[roomID][conn] = true
}

// unregisterClient removes a WebSocket connection from the pool for a specific room
func unregisterClient(roomID string, conn *websocket.Conn) {
    clientsLock.Lock()
    defer clientsLock.Unlock()
    if clients[roomID] != nil {
        delete(clients[roomID], conn)
        if len(clients[roomID]) == 0 {
            delete(clients, roomID)
        }
    }
}

// broadcastMessage sends a message to all connected WebSocket clients in a specific room
func broadcastMessage(roomID string, messageType int, message []byte) {
    clientsLock.Lock()
    defer clientsLock.Unlock()

    for conn := range clients[roomID] {
        err := conn.WriteMessage(messageType, message)
        if err != nil {
            log.Printf("Error broadcasting message: %v", err)
            conn.Close()
            delete(clients[roomID], conn)
        }
    }
}
