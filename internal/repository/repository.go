package repository

import (
	"database/sql"
	"github.com/redis/go-redis/v9"
)

type Repositories struct {
	User       UserRepository
	Item       ItemRepository
	Collection CollectionRepository
	Ranking    RankingRepository
}

func New(db *sql.DB, rdb *redis.Client) *Repositories {
	return &Repositories{
		User:       NewUserRepository(db),
		Item:       NewItemRepository(db, rdb),
		Collection: NewCollectionRepository(db),
		Ranking:    NewRankingRepository(rdb),
	}
}
