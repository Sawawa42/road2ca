package entity

import "github.com/google/uuid"

type Item struct {
	ID     uuid.UUID
	Name   string
	Rarity int
	Weight int
}
