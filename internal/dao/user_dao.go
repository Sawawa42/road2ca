package dao

import (
	"database/sql"
	"fmt"
	"road2ca/internal/model"
)

type UserDAO interface {
	Create(user *model.User) (int64, error)
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
