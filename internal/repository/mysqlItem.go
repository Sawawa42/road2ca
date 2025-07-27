package repository

import (
	"database/sql"
	"road2ca/internal/entity"
	"strings"
	"github.com/google/uuid"
)

type MySQLItemRepo interface {
	Save(items []*entity.Item) error
	Find() ([]*entity.Item, error)
}

type mysqlItemRepo struct {
	db *sql.DB
}

func NewMySQLItemRepo(db *sql.DB) MySQLItemRepo {
	return &mysqlItemRepo{
		db: db,
	}
}

func (r *mysqlItemRepo) Save(items []*entity.Item) error {
	query := "INSERT INTO items (id, name, rarity, weight) VALUES "
	var placeholders []string
	var args []interface{}
	for _, item := range items {
		placeholders = append(placeholders, "(?, ?, ?, ?)")
		uuid, err := uuid.NewV7()
		if err != nil {
			return err
		}
		args = append(args, uuid, item.Name, item.Rarity, item.Weight)
	}
	query += strings.Join(placeholders, ", ")
	query += " ON DUPLICATE KEY UPDATE name = VALUES(name), rarity = VALUES(rarity), weight = VALUES(weight)"
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

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
