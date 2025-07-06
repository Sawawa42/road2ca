package main

import (
	"flag"

	"database/sql"
	"log"
	"road2ca/internal/handler"
	"road2ca/internal/middleware"
	"road2ca/internal/repository"
	"road2ca/internal/server"
	"road2ca/internal/service"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"context"
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

	h, m := initServer(db)

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
	addr := ""
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
	})
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %+v", err)
	}
	return rdb
}

func initServer(db *sql.DB) (*handler.Handler, *middleware.Middleware) {
	r := repository.New(db)
	s := service.New(r)
	h := handler.New(s)
	m := middleware.New(s)

	return h, m
}
