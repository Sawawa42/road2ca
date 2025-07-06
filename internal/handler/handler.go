package handler

import (
	"road2ca/internal/service"
)

type Handler struct {
	User UserHandler
	Setting SettingHandler
}

func New(s *service.Services) *Handler {
	return &Handler{
		User:    NewUserHandler(s.User),
		Setting: NewSettingHandler(),
	}
}
