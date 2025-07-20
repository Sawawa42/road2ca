package repository

import (
	"database/sql"
	"road2ca/internal/entity"
	"strings"
)

type CollectionRepository interface {
	Save(tx *sql.Tx, collections []*entity.Collection) error
	FindAllByUserID(userID int) ([]*entity.Collection, error)
}

type collectionRepository struct {
	db *sql.DB
}

func NewCollectionRepository(db *sql.DB) CollectionRepository {
	return &collectionRepository{db: db}
}

// Save コレクションをDBに保存する
func (r *collectionRepository) Save(tx *sql.Tx, collections []*entity.Collection) error {
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

	if tx == nil {
		_, err := r.db.Exec(query, args...)
		return err
	}
	_, err := tx.Exec(query, args...)
	return err
}

// FindAllByUserID ユーザーIDに紐づくコレクションを全て取得する
func (r *collectionRepository) FindAllByUserID(userID int) ([]*entity.Collection, error) {
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
