package repository

import (
	"github.com/redis/go-redis/v9"
	"road2ca/internal/entity"
	"context"
	"encoding/json"
	"fmt"
)

type ItemRepository interface {
	SaveAll(items []*entity.Item) error
	FindByID(id int) (*entity.Item, error)
}

type itemRepository struct {
	rdb *redis.Client
}

func NewItemRepository(rdb *redis.Client) ItemRepository {
	return &itemRepository{rdb: rdb}
}

// SaveAll Redisにアイテム情報を保存する
func (r *itemRepository) SaveAll(items []*entity.Item) error {
	pipe := r.rdb.Pipeline()
	ctx := context.Background()
	for _, item := range items {
		json, err := json.Marshal(item)
		if err != nil {
			continue
		}
		key := fmt.Sprintf("item:%d", item.ID)
		pipe.Set(ctx, key, json, 0)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to save items: %w", err)
	}
	return nil
}

// FindByID Redisからアイテム情報を取得する
func (r *itemRepository) FindByID(id int) (*entity.Item, error) {
	ctx := context.Background()
	key := fmt.Sprintf("item:%d", id)
	val, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("item with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	var item entity.Item
	if err := json.Unmarshal([]byte(val), &item); err != nil {
		return nil, fmt.Errorf("failed to unmarshal item: %w", err)
	}
	return &item, nil
}
