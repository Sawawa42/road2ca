package entity

import "github.com/google/uuid"

type Collection struct {
	ID     uuid.UUID
	UserID uuid.UUID
	// 他2つと異なり事前に人間が設定するため、UUIDではなくint型を採用
	ItemID int
}
