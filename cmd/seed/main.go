package main

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"road2ca/internal/entity"
	"road2ca/internal/repository"
)

func seed(r *repository.Repositories) error {
	users := []*entity.User{
		{Name: "Alice", HighScore: 100, Token: "alice"},
		{Name: "Bob", HighScore: 200, Token: "bob"},
		{Name: "Charlie", HighScore: 150, Token: "charlie"},
		{Name: "Dave", HighScore: 300, Token: "dave"},
		{Name: "Eve", HighScore: 250, Token: "eve"},
		{Name: "Frank", HighScore: 400, Token: "frank"},
		{Name: "Grace", HighScore: 350, Token: "grace"},
		{Name: "Heidi", HighScore: 450, Token: "heidi"},
		{Name: "Ivan", HighScore: 500, Token: "ivan"},
		{Name: "Judy", HighScore: 1000, Token: "judy"},
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

func main() {
	// MySQL接続の初期化
	db := initMySQL()
	defer db.Close()

	// Redis接続の初期化
	rdb := initRedis()
	defer rdb.Close()

	r := repository.New(db, rdb)
	if err := seed(r); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}
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
