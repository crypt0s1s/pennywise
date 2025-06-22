package scissorspaperrock

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"terrible-ideas-2025/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var games = make(map[string]*Game)
var gamesMutex = &sync.Mutex{}

func GetAvailableGames() []Game {
	gamesMutex.Lock()
	defer gamesMutex.Unlock()

	var availableGames []Game
	now := time.Now()
	for _, game := range games {
		// Only include games that are less than 2 minutes old and don't have a second player
		if game.Player2ID == "" && now.Sub(game.CreatedAt) < 2*time.Minute {
			availableGames = append(availableGames, *game)
		}
	}
	return availableGames
}

func createGame(ownerID string, product models.Product) string {
	var gameID = uuid.New().String() // Assign gameID
	newGame := &Game{
		GameID:       gameID,
		Product:      product,
		Player1ID:    ownerID,
		Player2ID:    "", // Will be set when battle is accepted
		Player1Score: 0,
		Player2Score: 0,
		Player1Moves: []Move{},
		Player2Moves: []Move{},
		Status:       "in_progress",
		CreatedAt:    time.Now(),
	}
	gamesMutex.Lock()
	games[gameID] = newGame
	gamesMutex.Unlock()

	return gameID
}

func SetupScissorsPaperRockRoutes(router *gin.Engine, broadcast chan []byte) {
	log.Println("Setting up scissor paper rock routes")

	router.POST("/spr/initiate", func(c *gin.Context) {
		var req struct {
			UserID  string         `json:"user_id" binding:"required"`
			Product models.Product `json:"product" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate that product fields are not empty
		if req.Product.Name == "" || req.Product.Price == 0.0 || req.Product.URL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product details are required"})
			return
		}

		var gameID = createGame(req.UserID, req.Product)

		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Game created", "game_id": gameID})
	})

	router.POST("/spr/accept", func(c *gin.Context) {
		var req struct {
			UserID string `json:"user_id" binding:"required"`
			GameID string `json:"game_id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		gamesMutex.Lock()
		game, ok := games[req.GameID]
		if !ok {
			gamesMutex.Unlock()
			c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
			return
		}

		// Check if game already has a second player
		if game.Player2ID != "" {
			gamesMutex.Unlock()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Game already has two players"})
			return
		}

		// Check if player trying to join is the same as player 1
		if game.Player1ID == req.UserID {
			gamesMutex.Unlock()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot join your own game"})
			return
		}

		// Add player 2 to the game
		game.Player2ID = req.UserID

		// Notify all clients about the game update
		notification := gin.H{
			"type": "game_update",
			"game": game,
		}
		notificationBytes, _ := json.Marshal(notification)
		broadcast <- notificationBytes

		gamesMutex.Unlock()

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Successfully joined game",
			"game":    game,
		})

	})

	router.POST("/spr/:game_id", func(c *gin.Context) {
		gameID := c.Param("game_id")

		var req struct {
			UserID string     `json:"user_id" binding:"required"`
			Action GameAction `json:"action" binding:"required"`
			Round  int        `json:"round" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		gamesMutex.Lock()
		game, ok := games[gameID]
		if !ok {
			gamesMutex.Unlock()
			c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
			return
		}

		// Check if the user is part of this game
		if game.Player1ID != req.UserID && game.Player2ID != req.UserID {
			gamesMutex.Unlock()
			c.JSON(http.StatusForbidden, gin.H{"error": "User not part of this game"})
			return
		}

		// Log the move
		newMove := Move{
			Action: req.Action,
			Round:  req.Round,
		}

		// Check if move for this round already exists
		if req.UserID == game.Player1ID {
			for _, move := range game.Player1Moves {
				if move.Round == req.Round {
					gamesMutex.Unlock()
					c.JSON(http.StatusBadRequest, gin.H{"error": "Move already exists for this round"})
					return
				}
			}
			game.Player1Moves = append(game.Player1Moves, newMove)
		} else {
			for _, move := range game.Player2Moves {
				if move.Round == req.Round {
					gamesMutex.Unlock()
					c.JSON(http.StatusBadRequest, gin.H{"error": "Move already exists for this round"})
					return
				}
			}
			game.Player2Moves = append(game.Player2Moves, newMove)
		}

		// For simplicity, let's just broadcast the updated game state for now
		notification := gin.H{
			"type": "game_update",
			"game": game,
		}
		notificationBytes, _ := json.Marshal(notification)
		broadcast <- notificationBytes

		gamesMutex.Unlock()

		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Move logged", "game": game})
	})

	router.GET("/spr/:game_id", func(c *gin.Context) {
		gameID := c.Param("game_id")

		gamesMutex.Lock()
		game, ok := games[gameID]
		if !ok {
			gamesMutex.Unlock()
			c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
			return
		}
		gamesMutex.Unlock()

		c.JSON(http.StatusOK, gin.H{"status": "success", "game": game})
	})
}
