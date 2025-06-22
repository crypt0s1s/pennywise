package models

// BattleType defines the type of battle (e.g., scissors, paper, rock)

import "time"

type BattleType string

const (
	BattleTypeScissorsPaperRock BattleType = "scissors_paper_rock"
	BattleTypeArchery           BattleType = "archery"
)

// Battle represents a PvP battle instance
type Battle struct {
	ID        string     `json:"id"`
	Type      BattleType `json:"type"`
	Player1ID string     `json:"player1_id"`
	Player2ID string     `json:"player2_id"`
	Status    string     `json:"status"`  // e.g., "pending", "active", "completed"
	GameID    string     `json:"game_id"` // ID of the associated game (e.g., for SPR or Archery)
	CreatedAt time.Time  `json:"created_at"`
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
