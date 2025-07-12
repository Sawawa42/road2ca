package repository

import (
	"github.com/redis/go-redis/v9"
	"road2ca/internal/entity"
	"context"
	"encoding/json"
	"fmt"
	"database/sql"
)

type ItemRepository interface {
	SaveItemsToCache(items []*entity.Item) error
	FindItemByIdFromCache(id int) (*entity.Item, error)
	FindAllItemsFromDB() ([]*entity.Item, error)
}

type itemRepository struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewItemRepository(db *sql.DB, rdb *redis.Client) ItemRepository {
	return &itemRepository{
		db:  db,
		rdb: rdb,
	}
}

// SaveItemsToCache アイテム情報をキャッシュする
func (r *itemRepository) SaveItemsToCache(items []*entity.Item) error {
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

// FindItemByIdFromCache キャッシュからitemIDに対応するアイテムを取得する
func (r *itemRepository) FindItemByIdFromCache(id int) (*entity.Item, error) {
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

// FindAllItemsFromDB DBから全てのアイテムを取得する
func (r *itemRepository) FindAllItemsFromDB() ([]*entity.Item, error) {
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
