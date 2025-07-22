package main

import (
	"flag"

	"database/sql"
	"road2ca/internal/entity"
	"road2ca/internal/handler"
	"road2ca/internal/middleware"
	"road2ca/internal/repository"
	"road2ca/internal/server"
	"road2ca/internal/service"

	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"log"
	"math/rand"
	"time"
	"fmt"
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

	h, m, err := initServer(db, rdb)
	if err != nil {
		log.Fatalf("Failed to initialize server: %+v", err)
	}

	server.Serve(addr, h, m)
}

// connectDB MySQLデータベースに接続する
func initMySQL() *sql.DB {
	// TODO: .envからDSNを取得するようにする(正直今回は簡単のためハードコーディングでもいい気がする)
	dsn := "root:ca-tech-dojo@tcp(localhost:3306)/road2ca?parseTime=true"

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
	if rdb == nil {
		log.Fatal("Failed to create Redis client")
	}
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

	// シードデータの追加
	if err := seed(r); err != nil {
		return nil, nil, err
	}

	// itemをキャッシュ
	if err := s.Item.SetItemToCache(); err != nil {
		return nil, nil, err
	}

	return h, m, nil
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

// TODO: mainとは別cmdに切り出す
func seed(r *repository.Repositories) error {
	users := []*entity.User{
		{ID: 2, Name: "Alice", HighScore: 100, Token: "alice"},
		{ID: 3, Name: "Bob", HighScore: 200, Token: "bob"},
		{ID: 4, Name: "Charlie", HighScore: 150, Token: "charlie"},
		{ID: 5, Name: "Dave", HighScore: 300, Token: "dave"},
		{ID: 6, Name: "Eve", HighScore: 250, Token: "eve"},
		{ID: 7, Name: "Frank", HighScore: 400, Token: "frank"},
		{ID: 8, Name: "Grace", HighScore: 350, Token: "grace"},
		{ID: 9, Name: "Heidi", HighScore: 450, Token: "heidi"},
		{ID: 10, Name: "Ivan", HighScore: 500, Token: "ivan"},
		{ID: 11, Name: "Judy", HighScore: 1000, Token: "judy"},
	}

	for _, user := range users {
		if err := r.User.Save(user); err != nil {
			return err
		}
		if err := r.Ranking.Save(user); err != nil {
			return err
		}
	}

	return nil
}
