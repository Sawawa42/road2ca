package repository

import (
	"database/sql"
	"road2ca/internal/entity"
	"strings"
)

type UserRepo interface {
	Save(user *entity.User) error
	SaveTx(tx *sql.Tx, user *entity.User) error
	FindByToken(token string) (*entity.User, error)
	FindByIDs(ids [][]byte) ([]*entity.User, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) Save(user *entity.User) error {
	query := `
		INSERT INTO users (id, name, highscore, coin, token) VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		name = VALUES(name),
		highscore = VALUES(highscore),
		coin = VALUES(coin)`
	_, err := r.db.Exec(query, user.ID, user.Name, user.HighScore, user.Coin, user.Token)
	return err
}

func (r *userRepo) SaveTx(tx *sql.Tx, user *entity.User) error {
	query := `
		INSERT INTO users (id, name, highscore, coin, token) VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		name = VALUES(name),
		highscore = VALUES(highscore),
		coin = VALUES(coin)`
	_, err := tx.Exec(query, user.ID, user.Name, user.HighScore, user.Coin, user.Token)
	return err
}

func (r *userRepo) FindByToken(token string) (*entity.User, error) {
	query := "SELECT id, name, highscore, coin, token FROM users WHERE token = ?"
	row := r.db.QueryRow(query, token)
	user := &entity.User{}
	err := row.Scan(&user.ID, &user.Name, &user.HighScore, &user.Coin, &user.Token)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepo) FindByIDs(ids [][]byte) ([]*entity.User, error) {
	if len(ids) == 0 {
		return []*entity.User{}, nil
	}

	query := "SELECT id, name, highscore, coin, token FROM users WHERE id IN (?" + strings.Repeat(",?", len(ids)-1) + ")"
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		user := &entity.User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.HighScore, &user.Coin, &user.Token); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
