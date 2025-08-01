package repository

import (
	"database/sql"
	"road2ca/internal/entity"
	"strings"
	"log"
)

type MySQLItemRepo interface {
	// Save ItemsをDBに保存する
	Save(items []*entity.Item) error
	// Find 全てのItemsを取得する
	Find() ([]*entity.Item, error)
	// Truncate テーブルを空にする
	Truncate() error
}

type mysqlItemRepo struct {
	db *sql.DB
}

func NewMySQLItemRepo(db *sql.DB) MySQLItemRepo {
	return &mysqlItemRepo{
		db: db,
	}
}

// Save ItemsをDBに保存する
func (r *mysqlItemRepo) Save(items []*entity.Item) error {
	query := "INSERT INTO items (id, name, rarity, weight) VALUES "
	var placeholders []string
	var args []interface{}
	for _, item := range items {
		placeholders = append(placeholders, "(?, ?, ?, ?)")
		uuidBytes, err := GetUUIDv7Bytes()
		if err != nil {
			return err
		}
		args = append(args, uuidBytes, item.Name, item.Rarity, item.Weight)
	}
	query += strings.Join(placeholders, ", ")
	query += " ON DUPLICATE KEY UPDATE name = VALUES(name), rarity = VALUES(rarity), weight = VALUES(weight)"
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

// Find 全てのItemsを取得する
func (r *mysqlItemRepo) Find() ([]*entity.Item, error) {
	query := "SELECT * FROM items"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*entity.Item
	for rows.Next() {
		var item entity.Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Rarity, &item.Weight); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

// Truncate テーブルを空にする
func (r *mysqlItemRepo) Truncate() error {
	query := "SET FOREIGN_KEY_CHECKS = 0"
	success := false
	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}
	success = true
	defer func() {
		if !success {
			return
		}
		_, err := r.db.Exec("SET FOREIGN_KEY_CHECKS = 1")
		if err != nil {
			log.Printf("Error resetting foreign key checks: %v", err)
		}
	}()

	query = "TRUNCATE TABLE items"
	_, err = r.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
