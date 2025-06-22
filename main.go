package main

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	golf "terrible-ideas-2025/golf"
	models "terrible-ideas-2025/models"
	scissorspaperrock "terrible-ideas-2025/scissorsPaperRock"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, you should check the origin
	},
}

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex // Mutex to protect concurrent access to clients map
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		clients:    make(map[*websocket.Conn]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("websocket error: %s", err)
					client.Close()
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	hub.register <- conn

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go func() {
		// When this function returns, unregister the client and close the connection
		defer func() {
			hub.unregister <- conn
		}()
		for {
			// Read message from browser
			_, _, err := conn.ReadMessage()
			if err != nil {
				// If there's an error (e.g., client disconnected), break the loop
				break
			}
		}
	}()
}

var battles = make(map[string]*models.Battle)
var battlesMutex = &sync.Mutex{}

func main() {
	hub := newHub()
	go hub.run()

	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {

		// chrome-extension://fcakjiobnmmcgmhchmhedpokbemacee
		allowedOrigins := []string{"http://localhost:3000", "https://www.amazon.com.au", "chrome-extension://"}
		origin := c.Request.Header.Get("Origin")
		for _, allowedOrigin := range allowedOrigins {
			if strings.HasPrefix(origin, allowedOrigin) {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	scissorspaperrock.SetupScissorsPaperRockRoutes(router, hub.broadcast)
	golf.SetupGolfRoutes(router, hub.broadcast)

	router.GET("/ws", func(c *gin.Context) {
		serveWs(hub, c.Writer, c.Request)
	})

	router.GET("/games", func(c *gin.Context) {
		battlesMutex.Lock()
		defer battlesMutex.Unlock()

		var golfGames = golf.GetAvailableGames()
		var sprGames = scissorspaperrock.GetAvailableGames()

		// Combine both game types into a single slice
		allGames := make([]gin.H, 0)

		// Add golf games
		for _, game := range golfGames {
			allGames = append(allGames, gin.H{
				"game_id": game.GameID,
				"product": game.Product,
				"type":    "golf",
			})
		}

		// Add scissors paper rock games
		for _, game := range sprGames {
			allGames = append(allGames, gin.H{
				"game_id": game.GameID,
				"product": game.Product,
				"type":    "spr",
			})
		}

		c.JSON(http.StatusOK, allGames)
	})

	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
