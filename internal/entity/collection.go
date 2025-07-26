package entity

import "github.com/google/uuid"

type Collection struct {
	ID     uuid.UUID
	UserID uuid.UUID
	ItemID uuid.UUID
}
