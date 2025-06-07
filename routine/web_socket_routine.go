package routine

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"trading_platform_backend/model"
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

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	clients    map[*websocket.Conn]bool
	Broadcast  chan model.ApiResponse
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mutex      sync.Mutex
}

// Run starts the Hub's event loop. It must be run in its own goroutine.
func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			// Register a new client connection.
			h.mutex.Lock()
			h.clients[conn] = true
			h.mutex.Unlock()
			log.Println("Client registered")

			//TODO Send the current list of stocks to the newly connected client immediately.
			//stockMutex.RLock()
			//currentStocks := make([]Stock, 0, len(stocks))
			//for _, stock := range stocks {
			//	currentStocks = append(currentStocks, *stock)
			//}
			//stockMutex.RUnlock()

			data, err := json.Marshal("Client registered")
			if err != nil {
				log.Printf("Error marshalling initial stock data: %v", err)
				continue
			}
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("Error sending initial stock data to new client: %v", err)
			}

		case conn := <-h.unregister:
			// Unregister a client connection.
			h.mutex.Lock()
			if _, ok := h.clients[conn]; ok {
				delete(h.clients, conn)
				conn.Close()
				log.Println("Client unregistered")
			}
			h.mutex.Unlock()

		case message := <-h.Broadcast:
			// Broadcast a message to all registered clients.
			fmt.Println("New message: " + message.ErrorMessage)
			h.mutex.Lock()
			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshalling stock data: %v", err)
				continue
			}
			for conn := range h.clients {
				if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
					log.Printf("Write error: %v. Unregistering client.", err)
					// If there's an error writing (e.g., connection closed), unregister them.
					go func(c *websocket.Conn) { h.unregister <- c }(conn)
				}
			}
			h.mutex.Unlock()
		}
	}
}

// ServeWs handles WebSocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Called WebSocket ServeWs")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}

	// Register the new client.
	WsHub.register <- conn

	// This function will run as long as the client is connected.
	// It's mainly to detect when the client closes the connection.
	go func() {
		defer func() {
			WsHub.unregister <- conn
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
		Broadcast:  make(chan model.ApiResponse),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		clients:    make(map[*websocket.Conn]bool),
	}
	go WsHub.Run()
}
