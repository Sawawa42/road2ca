package repository

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
)

type Repositories struct {
	User       UserRepository
	Item       ItemRepo
	Collection CollectionRepo
	Ranking    RankingRepository
	DB         *sql.DB
}

func New(db *sql.DB, rdb *redis.Client) *Repositories {
	return &Repositories{
		User:       NewUserRepository(db),
		Item:       NewItemRepo(db, rdb),
		Collection: NewCollectionRepo(db),
		Ranking:    NewRankingRepository(rdb),
		DB:         db,
	}
}
