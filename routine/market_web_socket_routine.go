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

// Market upgrader takes a normal HTTP connection and upgrades it to a persistent WebSocket connection.
var marketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin allows us to configure which domains are allowed to connect.
	// For development, we can allow all origins.
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type MarketClient struct {
	Conn    *websocket.Conn
	StockID int64
}

// MarketHub maintains the set of active market clients and broadcasts messages to them.
type MarketHub struct {
	clients    map[*websocket.Conn]*MarketClient
	Broadcast  chan string
	register   chan *MarketClient
	unregister chan *MarketClient
	mutex      sync.Mutex
}

// Run starts the MarketHub's event loop. It must be run in its own goroutine.
func (h *MarketHub) Run() {
	for {
		select {
		case client := <-h.register:
			// Register a new client connection.
			h.mutex.Lock()
			h.clients[client.Conn] = client
			h.mutex.Unlock()
			log.Println("Market client registered")

			marketData := service.GetMarketData(client.StockID)

			data, err := json.Marshal(marketData)
			if err != nil {
				log.Printf("Error marshalling market data: %v", err)
				continue
			}

			if err := client.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("Error sending initial market data to new client: %v", err)
			}

		case client := <-h.unregister:
			// Unregister a client connection.
			h.mutex.Lock()
			if _, ok := h.clients[client.Conn]; ok {
				client.Conn.Close()
				delete(h.clients, client.Conn)
				log.Println("Market client unregistered")
			}
			h.mutex.Unlock()

		case _ = <-h.Broadcast:
			// Broadcast a message to all registered clients.
			h.mutex.Lock()

			for conn, client := range h.clients {
				marketData := service.GetMarketData(client.StockID)

				data, err := json.Marshal(marketData)
				if err != nil {
					log.Printf("Error marshalling market data: %v", err)
					continue
				}

				if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
					log.Printf("Write error: %v. Unregistering market client.", err)
					// If there's an error writing (e.g., connection closed), unregister them.
					go func(c *MarketClient) { h.unregister <- c }(client)
				}
			}
			h.mutex.Unlock()
		}
	}
}

// ServeMarketWs handles WebSocket requests from the peer for market data.
func ServeMarketWs(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Called Market WebSocket ServeMarketWs")

	// Parse stock ID from query parameters
	stockIDStr := r.URL.Query().Get("stockId")
	var stockID int64
	if stockIDStr == "" {
		fmt.Println("Stock ID is required")
		return
	} else {
		stockIDL, err := strconv.ParseInt(stockIDStr, 10, 64)
		if err != nil {
			fmt.Println("Invalid stock ID")
			return
		}
		stockID = stockIDL
	}

	conn, err := marketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Market upgrade error:", err)
		return
	}

	// Register the new client with stock ID.
	MarketWsHub.register <- &MarketClient{
		Conn:    conn,
		StockID: stockID,
	}

	// This function will run as long as the client is connected.
	// It's mainly to detect when the client closes the connection.
	go func() {
		defer func() {
			MarketWsHub.unregister <- &MarketClient{Conn: conn, StockID: stockID}
		}()
		for {
			// We must read from the connection to detect a close.
			// If we don't care about incoming messages, we can just discard them.
			if _, _, err := conn.ReadMessage(); err != nil {
				// This error will trigger when the client disconnects.
				log.Printf("Read error, market client disconnecting: %v", err)
				break
			}
		}
	}()
}

var MarketWsHub MarketHub

func initMarketWebSocket() {

	// Create and run the Market WebSocket hub.
	MarketWsHub = MarketHub{
		Broadcast:  make(chan string),
		register:   make(chan *MarketClient),
		unregister: make(chan *MarketClient),
		clients:    make(map[*websocket.Conn]*MarketClient),
	}
	go MarketWsHub.Run()
}
