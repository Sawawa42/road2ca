package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"road2ca/internal/entity"
	"github.com/redis/go-redis/v9"
)

type RedisItemRepo interface {
	Save(items []*entity.Item) error
	Find() ([]*entity.Item, error)
}

type redisItemRepo struct {
	rdb *redis.Client
}

func NewRedisItemRepo(rdb *redis.Client) RedisItemRepo {
	return &redisItemRepo{
		rdb: rdb,
	}
}

// Save アイテム情報をキャッシュする
func (r *redisItemRepo) Save(items []*entity.Item) error {
	pipe := r.rdb.Pipeline()
	ctx := context.Background()
	for _, item := range items {
		jsonData, err := json.Marshal(item)
		if err != nil {
			return err
		}
		key := fmt.Sprintf("item:%d", item.ID)
		pipe.Set(ctx, key, jsonData, 0)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}
	return nil
}

// Find アイテム情報を取得する。
func (r *redisItemRepo) Find() ([]*entity.Item, error) {
	ctx := context.Background()
	keys, err := r.rdb.Keys(ctx, "item:*").Result()
	if err != nil {
		return nil, err
	}

	var items []*entity.Item
	for _, key := range keys {
		jsonData, err := r.rdb.Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		var item entity.Item
		if err := json.Unmarshal([]byte(jsonData), &item); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	return items, nil
}
