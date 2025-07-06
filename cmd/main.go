package main

import (
	"flag"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"road2ca/internal/handler"
	"road2ca/internal/middleware"
	"road2ca/internal/server"
	"road2ca/internal/repository"
	"road2ca/internal/service"
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
	// データベース接続の初期化
	db := connectDB()
	defer db.Close()

	h, m := initServer(db)

	server.Serve(addr, h, m)
}

// connectDB MySQLデータベースに接続する
func connectDB() *sql.DB {
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

func initServer(db *sql.DB) (*handler.Handler, *middleware.Middleware) {
	r := repository.New(db)
	s := service.New(r)
	h := handler.New(s)
	m := middleware.New(s)

	return h, m
}
