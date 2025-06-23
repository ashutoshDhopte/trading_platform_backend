package routine

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"trading_platform_backend/service"

	"github.com/gorilla/websocket"
)

// upgrader takes a normal HTTP connection and upgrades it to a persistent WebSocket connection.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin allows us to configure which domains are allowed to connect.
	// For development, we can allow all origins.
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	UserID int64
	Conn   *websocket.Conn
}

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	clients    map[int64]*websocket.Conn
	Broadcast  chan string
	register   chan Client
	unregister chan int64
	mutex      sync.Mutex
}

// Run starts the Hub's event loop. It must be run in its own goroutine.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// Register a new client connection.
			h.mutex.Lock()
			h.clients[client.UserID] = client.Conn
			h.mutex.Unlock()
			log.Println("Client registered")

			dashboard := service.GetDashboardData(client.UserID)

			data, err := json.Marshal(dashboard)
			if err != nil {
				log.Printf("Error marshalling stock data: %v", err)
				continue
			}

			if err := client.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("Error sending initial stock data to new client: %v", err)
			}

		case userId := <-h.unregister:
			// Unregister a client connection.
			h.mutex.Lock()
			if _, ok := h.clients[userId]; ok {
				h.clients[userId].Close()
				delete(h.clients, userId)
				log.Println("Client unregistered")
			}
			h.mutex.Unlock()

		case _ = <-h.Broadcast:
			// Broadcast a message to all registered clients.
			h.mutex.Lock()

			for userId, conn := range h.clients {

				dashboard := service.GetDashboardData(userId)

				data, err := json.Marshal(dashboard)
				if err != nil {
					log.Printf("Error marshalling stock data: %v", err)
					continue
				}

				if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
					log.Printf("Write error: %v. Unregistering client.", err)
					// If there's an error writing (e.g., connection closed), unregister them.
					go func(u int64) { h.unregister <- u }(userId)
				}
			}
			h.mutex.Unlock()
		}
	}
}

// ServeWs handles WebSocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Called WebSocket ServeWs")

	userIdStr := r.URL.Query().Get("userId")
	var userId int64
	if userIdStr == "" {
		fmt.Println("User id error")
		return
	} else {
		userIdL, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			fmt.Println("User id")
			return
		}
		userId = userIdL
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}

	// Register the new client.
	WsHub.register <- Client{
		UserID: userId,
		Conn:   conn,
	}

	// This function will run as long as the client is connected.
	// It's mainly to detect when the client closes the connection.
	go func() {
		defer func() {
			WsHub.unregister <- userId
		}()
		for {
			// We must read from the connection to detect a close.
			// If we don't care about incoming messages, we can just discard them.
			if _, _, err := conn.ReadMessage(); err != nil {
				// This error will trigger when the client disconnects.
				log.Printf("Read error, client disconnecting: %v", err)
				break
			}
		}
	}()
}

var WsHub Hub

func initWebSocket() {

	// Create and run the WebSocket hub.
	WsHub = Hub{
		Broadcast:  make(chan string),
		register:   make(chan Client),
		unregister: make(chan int64),
		clients:    make(map[int64]*websocket.Conn),
	}
	go WsHub.Run()
}
