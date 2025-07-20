package repository

import (
	"database/sql"
	"road2ca/internal/entity"
)

type UserRepository interface {
	Save(tx *sql.Tx, user *entity.User) error
	FindByToken(token string) (*entity.User, error)
	FindByID(id int) (*entity.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Save(tx *sql.Tx, user *entity.User) error {
	query := `
		INSERT INTO users (name, highscore, coin, token) VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		name = VALUES(name),
		highscore = VALUES(highscore),
		coin = VALUES(coin)`
	if tx == nil {
		_, err := tx.Exec(query, user.Name, user.HighScore, user.Coin, user.Token)
		return err
	}
	_, err := r.db.Exec(query, user.Name, user.HighScore, user.Coin, user.Token)
	return err
}

func (r *userRepository) FindByToken(token string) (*entity.User, error) {
	query := "SELECT id, name, highscore, coin, token FROM users WHERE token = ?"
	row := r.db.QueryRow(query, token)
	user := &entity.User{}
	err := row.Scan(&user.ID, &user.Name, &user.HighScore, &user.Coin, &user.Token)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) FindByID(id int) (*entity.User, error) {
	query := "SELECT id, name, highscore, coin, token FROM users WHERE id = ?"
	row := r.db.QueryRow(query, id)
	user := &entity.User{}
	err := row.Scan(&user.ID, &user.Name, &user.HighScore, &user.Coin, &user.Token)
	if err != nil {
		return nil, err
	}
	return user, nil
}
