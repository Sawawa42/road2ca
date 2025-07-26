package repository

import (
	"database/sql"
	"github.com/redis/go-redis/v9"
)

type Repositories struct {
	User       UserRepo
	// Item       ItemRepo
	MySQLItem  MySQLItemRepo
	RedisItem  RedisItemRepo
	Collection CollectionRepo
	Ranking    RankingRepo
	Setting    SettingRepo
	DB         *sql.DB
}

func New(db *sql.DB, rdb *redis.Client) *Repositories {
	return &Repositories{
		User:       NewUserRepo(db),
		// Item:       NewItemRepo(db, rdb),
		MySQLItem:  NewMySQLItemRepo(db),
		RedisItem:  NewRedisItemRepo(rdb),
		Collection: NewCollectionRepo(db),
		Ranking:    NewRankingRepo(rdb),
		Setting:    NewSettingRepo(db, rdb),
		DB:         db,
	}
}
