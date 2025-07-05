package main

import (
	"flag"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"road2ca/pkg/http/middleware"
	"road2ca/pkg/server"
	"road2ca/pkg/server/handler"
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

	// ハンドラの初期化
	h := handler.New(db)

	// ミドルウェアの初期化
	m := middleware.NewMiddleware(db)

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
