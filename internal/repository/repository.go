package repository

import (
	"database/sql"
	"github.com/redis/go-redis/v9"
)

type Repositories struct {
	User         UserRepo
	MySQLItem    MySQLItemRepo
	RedisItem    RedisItemRepo
	MySQLSetting MySQLSettingRepo
	RedisSetting RedisSettingRepo
	Collection   CollectionRepo
	Ranking      RankingRepo
	DB           *sql.DB
}

func New(db *sql.DB, rdb *redis.Client) *Repositories {
	return &Repositories{
		User:         NewUserRepo(db),
		MySQLItem:    NewMySQLItemRepo(db),
		RedisItem:    NewRedisItemRepo(rdb),
		MySQLSetting: NewMySQLSettingRepo(db),
		RedisSetting: NewRedisSettingRepo(rdb),
		Collection:   NewCollectionRepo(db),
		Ranking:      NewRankingRepo(rdb),
		DB:           db,
	}
}
