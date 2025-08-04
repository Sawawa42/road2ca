package server

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"log"
	"math/rand"
	"os"
	"road2ca/internal/handler"
	"road2ca/internal/middleware"
	"road2ca/internal/repository"
	"road2ca/internal/service"
	"time"
	"road2ca/internal/entity"
	"math"
)

// InitMySQL MySQL接続の初期化
func InitMySQL() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	dsn := os.Getenv("DSN")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %+v", err)
	}

	// DB接続の確認
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %+v", err)
	}
	return db
}

// InitRedis Redis接続の初期化
func InitRedis() *redis.Client {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	addr := os.Getenv("REDIS_ADDR")
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %+v", err)
	}
	return rdb
}

// SetupServer サーバーの初期設定
func SetupServer(db *sql.DB, rdb *redis.Client) (*handler.Handler, *middleware.Middleware, middleware.Logger, error) {
	// 依存関係の注入
	r, s, h, m := injectDependencies(db, rdb)
	l, err := middleware.NewLogger()
	if err != nil {
		return nil, nil, nil, err
	}

	// MySQLから設定を取得
	setting, err := r.MySQLSetting.FindLatest()
	if err != nil {
		return nil, nil, nil, err
	}

	// MySQLからアイテムを取得
	items, err := r.MySQLItem.Find()
	if err != nil || len(items) == 0 {
		if err != nil {
			return nil, nil, nil, err
		}
		return nil, nil, nil, fmt.Errorf("no items found")
	}

	// Redisに設定をキャッシュ
	if err := r.RedisSetting.Save(setting); err != nil {
		return nil, nil, nil, err
	}

	// アイテムに重みを設定
	if err := setWeightToItems(items, setting); err != nil {
		return nil, nil, nil, err
	}

	// Redisにアイテムをキャッシュ
	if err := r.RedisItem.Save(items); err != nil {
		return nil, nil, nil, err
	}

	// MySQLに重みを更新したアイテムを保存
	if err := r.MySQLItem.Save(items); err != nil {
		return nil, nil, nil, err
	}

	props, err := loadGachaServiceProps(items)
	if err != nil {
		return nil, nil, nil, err
	}
	s.Gacha.SetGachaProps(props)

	return h, m, l, nil
}

func loadGachaServiceProps(items []*entity.Item) (*service.GachaServiceProps, error) {
	totalWeight := 0
	for _, item := range items {
		if item.Weight < 1 {
			return nil, fmt.Errorf("item has invalid weight: %s", item.Name)
		}
		totalWeight += item.Weight
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return &service.GachaServiceProps{
		TotalWeight: totalWeight,
		RandGen:     r,
	}, nil
}

// injectDependencies 依存関係の注入
func injectDependencies(db *sql.DB, rdb *redis.Client) (
	*repository.Repositories,
	*service.Services,
	*handler.Handler,
	*middleware.Middleware) {
	r := repository.New(db, rdb)
	s := service.New(r)
	h := handler.New(s)
	m := middleware.New(s)
	return r, s, h, m
}

// setWeightToItems アイテムの重みを設定する
func setWeightToItems(items []*entity.Item, setting *entity.Setting) error {
	sum := setting.Rarity3Ratio + setting.Rarity2Ratio + setting.Rarity1Ratio
	const epsilon = 1e-6
	if math.Abs(sum-100.0) > epsilon {
		return fmt.Errorf("total rarity ratio must be 100, got %f", sum)
	}

	const rarityRatioScale = 100000
	totalRarityWeights := map[int]int{
		3: int(math.Round(setting.Rarity3Ratio * float64(rarityRatioScale))),
		2: int(math.Round(setting.Rarity2Ratio * float64(rarityRatioScale))),
		1: int(math.Round(setting.Rarity1Ratio * float64(rarityRatioScale))),
	}

	itemsByRarity := make(map[int][]*entity.Item)
	for _, item := range items {
		itemsByRarity[item.Rarity] = append(itemsByRarity[item.Rarity], item)
	}

	for rarity, itemsInGroup := range itemsByRarity {
		totalWeight := totalRarityWeights[rarity]
		count := len(itemsInGroup)

		if count == 0 {
			continue
		}

		baseWeight := totalWeight / count
		remainder := totalWeight % count

		for i, item := range itemsInGroup {
			item.Weight = baseWeight
			if i < remainder {
				item.Weight++
			}
		}
	}

	return nil
}
