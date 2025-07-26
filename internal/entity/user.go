package entity

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID
	Name      string
	HighScore int
	Coin      int
	Token     string
}
