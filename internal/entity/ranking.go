package entity

import "github.com/google/uuid"

type Ranking struct {
	UserID uuid.UUID
	Score  int
	Rank   int
}
