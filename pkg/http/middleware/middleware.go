package middleware

import (
	"database/sql"
	"road2ca/pkg/dao"
)

type Middleware struct {
	userDAO dao.UserDAO
}

func NewMiddleware(db *sql.DB) *Middleware {
	return &Middleware{
		userDAO: dao.NewUserDAO(db),
	}
}
