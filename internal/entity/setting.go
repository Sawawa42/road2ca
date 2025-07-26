package entity

import "github.com/google/uuid"

type Setting struct {
	ID                   uuid.UUID
	Name                 string
	GachaCoinConsumption int
	DrawGachaMaxTimes    int
	GetRankingLimit      int
	RewardCoin           int
}
