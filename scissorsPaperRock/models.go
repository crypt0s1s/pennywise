package scissorspaperrock

import (
	"time"

	"terrible-ideas-2025/models"
)

// GameAction defines the possible actions in a game of scissors, paper, rock
type GameAction string

const (
	GameActionScissors GameAction = "scissors"
	GameActionPaper    GameAction = "paper"
	GameActionRock     GameAction = "rock"
)

// Move represents a single move in a game of scissors, paper, rock
type Move struct {
	Round  int        `json:"round"`
	Action GameAction `json:"action"`
}

// Game represents a single instance of a scissors, paper, rock game
type Game struct {
	GameID       string         `json:"game_id"`
	Product      models.Product `json:"product"`
	Player1ID    string         `json:"player1_id"`
	Player2ID    string         `json:"player2_id"`
	Player1Score int            `json:"player1_score"`
	Player2Score int            `json:"player2_score"`
	Player1Moves []Move         `json:"player1_moves"`
	Player2Moves []Move         `json:"player2_moves"`
	Status       string         `json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
}
