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
	"os/signal"
	"syscall"
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
	db := initMySQL()
	defer db.Close()

	rdb := initRedis()
	defer rdb.Close()

	// Ctrl+C(SIGINT)で終了した際の処理
	sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-sigs
        // Redisキャッシュをクリア
        if err := rdb.FlushAll(context.Background()).Err(); err != nil {
            log.Printf("Failed to clear Redis cache: %v", err)
        }
		db.Close()
		rdb.Close()
        os.Exit(0)
    }()

	h, m, err := initServer(db, rdb)
	if err != nil {
		log.Fatalf("Failed to initialize server: %+v", err)
	}

	server.Serve(addr, h, m)
}

// initMySQL MySQL接続の初期化
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

	s := service.New(r)
	h := handler.New(s)
	m := middleware.New(s)

	if err := setDataToCache(s); err != nil {
		return nil, nil, err
	}

	props, err := loadGachaServiceProps(r.MySQLItem, r.RedisItem)
	if err != nil {
		return nil, nil, err
	}
	s.Gacha.SetGachaProps(props)

	return h, m, nil
}

// setDataToCache settingとitemを取得し、キャッシュに保存する
func setDataToCache(s *service.Services) error {
	if err := s.Setting.SetSettingToCache(); err != nil {
		log.Printf("Failed to set settings to cache: %v", err)
		return err
	}

	if err := s.Item.SetItemToCache(); err != nil {
		log.Printf("Failed to set items to cache: %v", err)
		return err
	}

	return nil
}

func loadGachaServiceProps(
		mySqlItemRepo repository.MySQLItemRepo,
		redisItemRepo repository.RedisItemRepo,
	) (*service.GachaServiceProps, error) {
	items, err := repository.FindItems(mySqlItemRepo, redisItemRepo)
	if err != nil {
		return nil, err
	}

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
