package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
	"road2ca/internal/repository"
	"road2ca/internal/seed"
)

func main() {
	db := initMySQL()
	defer db.Close()

	mi := repository.NewMySQLItemRepo(db)
	ms := repository.NewMySQLSettingRepo(db)
	mc := repository.NewCollectionRepo(db)
	if err := seed.Seed(mi, ms, mc); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}
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
