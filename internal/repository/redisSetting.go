package repository

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"road2ca/internal/entity"
)

type RedisSettingRepo interface {
	Save(setting *entity.Setting) error
	FindLatest() (*entity.Setting, error)
}

type redisSettingRepo struct {
	rdb *redis.Client
}

func NewRedisSettingRepo(rdb *redis.Client) RedisSettingRepo {
	return &redisSettingRepo{
		rdb: rdb,
	}
}

// Save 設定情報をキャッシュする
func (r *redisSettingRepo) Save(setting *entity.Setting) error {
	ctx := context.Background()
	key := "setting:latest"
	jsonData, err := json.Marshal(setting)
	if err != nil {
		return err
	}
	err = r.rdb.Set(ctx, key, jsonData, 0).Err()
	return err
}

// FindLatest キャッシュから最新の設定情報を取得する
func (r *redisSettingRepo) FindLatest() (*entity.Setting, error) {
	ctx := context.Background()
	key := "setting:latest"

	// キャッシュから取得
	jsonData, err := r.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // キャッシュに存在しない場合はnilを返す
	} else if err != nil {
		return nil, err
	}

	var setting entity.Setting
	err = json.Unmarshal([]byte(jsonData), &setting)
	if err != nil {
		return nil, err
	}
	return &setting, nil
}
