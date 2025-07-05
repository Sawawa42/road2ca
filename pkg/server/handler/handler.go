package handler

import (
	"database/sql"
	"road2ca/pkg/dao"
)

type Handler struct {
	userDAO dao.UserDAO
}

func New(db *sql.DB) *Handler {
	return &Handler{
		userDAO: dao.NewUserDAO(db),
	}
}
