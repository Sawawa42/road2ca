package entity

import "github.com/gofrs/uuid"

type User struct {
	ID        uuid.UUID
	Name      string
	HighScore int
	Coin      int
	Token     string
}
