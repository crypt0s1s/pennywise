package models

// BattleType defines the type of battle (e.g., scissors, paper, rock)
type BattleType string

const (
	BattleTypeScissorsPaperRock BattleType = "scissors_paper_rock"
	BattleTypeArchery           BattleType = "archery"
)

// Battle represents a PvP battle instance
type Battle struct {
	ID           string     `json:"id"`
	ItemName     string     `json:"item_name"`
	ItemCost     float64    `json:"item_cost"`
	BattleType   BattleType `json:"battle_type"`
	OwnerID      string     `json:"owner_id"`
	ChallengerID string     `json:"challenger_id,omitempty"`
	Status       string     `json:"status"` // e.g., "pending", "accepted", "completed"
	GameID       string     `json:"game_id,omitempty"`
}

// BattleRequest represents the request body for initiating a battle
type BattleRequest struct {
	ItemName   string     `json:"item_name" binding:"required"`
	ItemCost   float64    `json:"item_cost" binding:"required"`
	BattleType BattleType `json:"battle_type" binding:"required"`
	OwnerID    string     `json:"owner_id" binding:"required"`
}

// AcceptBattleRequest represents the request body for accepting a battle
type AcceptBattleRequest struct {
	BattleID     string `json:"battle_id" binding:"required"`
	ChallengerID string `json:"challenger_id" binding:"required"`
}

type ArcheryGame struct {
	GameID       string `json:"game_id"`
	BattleID     string `json:"battle_id"`
	Player1ID    string `json:"player1_id"`
	Player2ID    string `json:"player2_id"`
	Player1Score int    `json:"player1_score"`
	Player2Score int    `json:"player2_score"`
	Status       string `json:"status"`
}
