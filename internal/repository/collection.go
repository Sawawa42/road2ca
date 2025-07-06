package repository

import (
	"database/sql"
	"road2ca/internal/entity"
)

type CollectionRepository interface {
	Save(collection *entity.Collection) error
}

type collectionRepository struct {
	db *sql.DB
}

func NewCollectionRepository(db *sql.DB) CollectionRepository {
	return &collectionRepository{db: db}
}

func (r *collectionRepository) Save(collection *entity.Collection) error {
	query := `
		INSERT INTO collections (user_id, item_id) VALUES (?, ?)
		ON DUPLICATE KEY UPDATE
		item_id = VALUES(item_id)`
	_, err := r.db.Exec(query, collection.UserID, collection.ItemID)
	if err != nil {
		return err
	}
	return nil
}
