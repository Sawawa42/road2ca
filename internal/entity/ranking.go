package entity

import "github.com/gofrs/uuid"

type Ranking struct {
	UserID uuid.UUID
	Score  int
	Rank   int
}
