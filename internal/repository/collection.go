package repository

import (
	"database/sql"
	"road2ca/internal/entity"
	"strings"
)

type CollectionRepo interface {
	Save(tx *sql.Tx, collections []*entity.Collection) error
	FindByUserID(userID int) ([]*entity.Collection, error)
}

type collectionRepo struct {
	db *sql.DB
}

func NewCollectionRepo(db *sql.DB) CollectionRepo {
	return &collectionRepo{db: db}
}

// Save コレクションをDBに保存する
func (r *collectionRepo) Save(tx *sql.Tx, collections []*entity.Collection) error {
	if len(collections) == 0 {
		return nil
	}

	query := "INSERT INTO collections (userId, itemId) VALUES "
	var placeholders []string
	var args []interface{}
	for _, collection := range collections {
		placeholders = append(placeholders, "(?, ?)")
		args = append(args, collection.UserID, collection.ItemID)
	}
	query += strings.Join(placeholders, ", ")
	query += " ON DUPLICATE KEY UPDATE itemId = VALUES(itemId)"

	_, err := tx.Exec(query, args...)
	return err
}

// FindByUserID ユーザーIDに紐づくコレクションを取得する
func (r *collectionRepo) FindByUserID(userID int) ([]*entity.Collection, error) {
	query := `SELECT id, userId, itemId FROM collections WHERE userId = ?`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collections []*entity.Collection
	for rows.Next() {
		var collection entity.Collection
		if err := rows.Scan(&collection.ID, &collection.UserID, &collection.ItemID); err != nil {
			return nil, err
		}
		collections = append(collections, &collection)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return collections, nil
}
