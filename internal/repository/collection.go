package repository

import (
	"database/sql"
	"road2ca/internal/entity"
)

type CollectionRepository interface {
	Save(collection *entity.Collection) error
	FindAllByUserID(userID int) ([]*entity.Collection, error)
}

type collectionRepository struct {
	db *sql.DB
}

func NewCollectionRepository(db *sql.DB) CollectionRepository {
	return &collectionRepository{db: db}
}

// Save コレクションをDBに保存する
func (r *collectionRepository) Save(collection *entity.Collection) error {
	query := `
		INSERT INTO collections (userId, itemId) VALUES (?, ?)
		ON DUPLICATE KEY UPDATE
		itemId = VALUES(itemId)`
	_, err := r.db.Exec(query, collection.UserID, collection.ItemID)
	if err != nil {
		return err
	}
	return nil
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
