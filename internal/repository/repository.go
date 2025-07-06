package repository

import (
	"database/sql"
	"github.com/redis/go-redis/v9"
)

type Repositories struct {
	User UserRepository
	Item ItemRepository
}

func New(db *sql.DB, rdb *redis.Client) *Repositories {
	return &Repositories{
		User: NewUserRepository(db),
		Item: NewItemRepository(db, rdb),
	}
}
