package repository

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
)

type Repositories struct {
	User       UserRepo
	Item       ItemRepo
	Collection CollectionRepo
	Ranking    RankingRepo
	DB         *sql.DB
}

func New(db *sql.DB, rdb *redis.Client) *Repositories {
	return &Repositories{
		User:       NewUserRepo(db),
		Item:       NewItemRepo(db, rdb),
		Collection: NewCollectionRepo(db),
		Ranking:    NewRankingRepo(rdb),
		DB:         db,
	}
}
