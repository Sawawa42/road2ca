package main

import (
	"context"
	"database/sql"
	"flag"
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
	"road2ca/internal/server"
	"road2ca/internal/service"
	"time"
)

var (
	// Listenするアドレス+ポート
	addr string
)

func init() {
	flag.StringVar(&addr, "addr", ":8080", "tcp host:port to connect")
	flag.Parse()
}

func main() {
	// MySQL接続の初期化
	db := initMySQL()
	defer db.Close()

	// Redis接続の初期化
	rdb := initRedis()
	defer rdb.Close()

	h, m, err := initServer(db, rdb)
	if err != nil {
		log.Fatalf("Failed to initialize server: %+v", err)
	}

	server.Serve(addr, h, m)
}

// initMySQL MySQLデータベースに接続する
func initMySQL() *sql.DB {
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

// initRedis Redis接続の初期化
func initRedis() *redis.Client {
	addr := "localhost:6379"
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %+v", err)
	}
	return rdb
}

func initServer(db *sql.DB, rdb *redis.Client) (*handler.Handler, *middleware.Middleware, error) {
	r := repository.New(db, rdb)

	gachaProps, err := loadGachaServiceProps(r.Item)
	if err != nil {
		return nil, nil, err
	}

	s := service.New(r, gachaProps)
	h := handler.New(s)
	m := middleware.New(s)

	// マスターデータをキャッシュに設定
	if err := setMasterDataToCache(s); err != nil {
		return nil, nil, err
	}

	return h, m, nil
}

func setMasterDataToCache(s *service.Services) error {
	// 設定をキャッシュ
	if err := s.Setting.SetSettingToCache(); err != nil {
		return err
	}

	// itemをキャッシュ
	if err := s.Item.SetItemToCache(); err != nil {
		return err
	}

	return nil
}

func loadGachaServiceProps(itemRepo repository.ItemRepo) (*service.GachaServiceProps, error) {
	items, err := itemRepo.Find()
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("no items found in cache")
	}

	totalWeight := 0
	for _, item := range items {
		if item.Weight == 0 {
			continue // 重みが0のアイテムは無視する
		} else if item.Weight < 0 {
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
