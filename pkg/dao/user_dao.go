package dao

import (
	"database/sql"
	"fmt"
	"road2ca/pkg/model"
)

type UserDAO interface {
	Create(user *model.User) (int64, error)
	GetByToken(token string) (*model.User, error)
}

type userDAOImpl struct {
	db *sql.DB
}

func NewUserDAO(db *sql.DB) UserDAO {
	return &userDAOImpl{db: db}
}

func (dao *userDAOImpl) Create(user *model.User) (int64, error) {
	query := "INSERT INTO users (name, token) VALUES (?, ?)"
	result, err := dao.db.Exec(query, user.Name, user.Token)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	userID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}
	return userID, nil
}

func (dao *userDAOImpl) GetByToken(token string) (*model.User, error) {
	query := "SELECT id, name, highscore, coin, token FROM users WHERE token = ?"
	row := dao.db.QueryRow(query, token)
	user := &model.User{}
	err := row.Scan(&user.ID, &user.Name, &user.HighScore, &user.Coin, &user.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // ユーザが存在しない場合はnilを返す
		}
		return nil, fmt.Errorf("failed to get user by token: %w", err)
	}
	return user, nil
}
