package repository

import (
	// "context"
	// "database/sql"
	// "encoding/json"
	// "fmt"
	"log"
	"road2ca/internal/entity"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func FindItems(mysqlRepo MySQLItemRepo, redisRepo RedisItemRepo) ([]*entity.Item, error) {
	items, err := redisRepo.Find()
	if err != nil || len(items) == 0 {
		// キャッシュにアイテムがない場合はMySQLから取得
		if err == redis.Nil {
			items, err = mysqlRepo.Find()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	log.Printf("Found %d items in cache", len(items))
	return items, nil
}

func GetUUIDv7Bytes() ([]byte, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	return uuid.MarshalBinary()
}

// type ItemRepo interface {
// 	Save(items []*entity.Item) error
// 	Find() ([]*entity.Item, error)
// }

// type itemRepo struct {
// 	db  *sql.DB
// 	rdb *redis.Client
// }

// func NewItemRepo(db *sql.DB, rdb *redis.Client) ItemRepo {
// 	return &itemRepo{
// 		db:  db,
// 		rdb: rdb,
// 	}
// }

// // Save アイテム情報をキャッシュする
// func (r *itemRepo) Save(items []*entity.Item) error {
// 	pipe := r.rdb.Pipeline()
// 	ctx := context.Background()
// 	for _, item := range items {
// 		json, err := json.Marshal(item)
// 		if err != nil {
// 			continue
// 		}
// 		key := fmt.Sprintf("item:%s", item.ID)
// 		pipe.Set(ctx, key, json, 0)
// 	}

// 	if _, err := pipe.Exec(ctx); err != nil {
// 		return fmt.Errorf("failed to save items: %w", err)
// 	}
// 	return nil
// }

// // Find アイテム情報を取得する。キャッシュに存在しない場合はDBから取得する
// func (r *itemRepo) Find() ([]*entity.Item, error) {
// 	ctx := context.Background()
// 	keys, err := r.rdb.Keys(ctx, "item:*").Result()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get item keys: %w", err)
// 	}

// 	var items []*entity.Item
// 	if len(keys) == 0 {
// 		// キャッシュにアイテムがない場合DBから取得
// 		items, err = r.findFromDB()
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to find items from DB: %w", err)
// 		}
// 	} else {
// 		// キャッシュにアイテムがある場合はキャッシュから取得
// 		items, err = r.findFromCache()
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to find items from cache: %w", err)
// 		}
// 	}

// 	return items, nil
// }

// // findFromDB DBからアイテムを取得する
// func (r *itemRepo) findFromDB() ([]*entity.Item, error) {
// 	query := "SELECT * FROM items"
// 	rows, err := r.db.Query(query)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to query items: %w", err)
// 	}
// 	defer rows.Close()
// 	var items []*entity.Item
// 	for rows.Next() {
// 		var item entity.Item
// 		if err := rows.Scan(&item.ID, &item.Name, &item.Rarity, &item.Weight); err != nil {
// 			return nil, fmt.Errorf("failed to scan item: %w", err)
// 		}
// 		items = append(items, &item)
// 	}

// 	return items, nil
// }

// // findFromCache キャッシュからアイテムを取得する
// func (r *itemRepo) findFromCache() ([]*entity.Item, error) {
// 	ctx := context.Background()
// 	keys, err := r.rdb.Keys(ctx, "item:*").Result()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get item keys: %w", err)
// 	}

// 	var items []*entity.Item
// 	for _, key := range keys {
// 		val, err := r.rdb.Get(ctx, key).Result()
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to get item from cache: %w", err)
// 		}

// 		var item entity.Item
// 		if err := json.Unmarshal([]byte(val), &item); err != nil {
// 			return nil, fmt.Errorf("failed to unmarshal item: %w", err)
// 		}
// 		items = append(items, &item)
// 	}

// 	return items, nil
// }
