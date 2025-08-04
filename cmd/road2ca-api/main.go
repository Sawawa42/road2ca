package main

import (
	"context"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"os/signal"
	"road2ca/internal/server"
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
	db := server.InitMySQL()
	defer db.Close()

	rdb := server.InitRedis()
	defer rdb.Close()

	slogs, err := server.NewSlogInstances()
	if err != nil {
		log.Fatalf("Failed to initialize slog instances: %v", err)
	}

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

	h, m, err := server.SetupServer(db, rdb, slogs)
	if err != nil {
		log.Fatalf("Failed to initialize server: %+v", err)
	}

	server.Serve(addr, h, m)
}
