package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"road2ca/internal/entity"

	"github.com/redis/go-redis/v9"
)

type ItemRepo interface {
	SaveToCache(items []*entity.Item) error
	FindFromDB() ([]*entity.Item, error)
	FindFromCache() ([]*entity.Item, error)
}

type itemRepo struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewItemRepo(db *sql.DB, rdb *redis.Client) ItemRepo {
	return &itemRepo{
		db:  db,
		rdb: rdb,
	}
}

// SaveToCache アイテム情報をキャッシュする
func (r *itemRepo) SaveToCache(items []*entity.Item) error {
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

// FindFromDB DBからアイテムを取得する
func (r *itemRepo) FindFromDB() ([]*entity.Item, error) {
	query := "SELECT * FROM items"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query items: %w", err)
	}
	defer rows.Close()
	var items []*entity.Item
	for rows.Next() {
		var item entity.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Rarity, &item.Weight); err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, &item)
	}

	return items, nil
}

// FindFromCache キャッシュからアイテムを取得する
func (r *itemRepo) FindFromCache() ([]*entity.Item, error) {
	ctx := context.Background()
	keys, err := r.rdb.Keys(ctx, "item:*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get item keys: %w", err)
	}

	var items []*entity.Item
	for _, key := range keys {
		val, err := r.rdb.Get(ctx, key).Result()
		if err != nil {
			if err == redis.Nil {
				continue // アイテムが存在しない場合はスキップ
			}
			return nil, fmt.Errorf("failed to get item from cache: %w", err)
		}

		var item entity.Item
		if err := json.Unmarshal([]byte(val), &item); err != nil {
			return nil, fmt.Errorf("failed to unmarshal item: %w", err)
		}
		items = append(items, &item)
	}

	return items, nil
}
