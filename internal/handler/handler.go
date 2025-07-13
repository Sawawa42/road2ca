package handler

import (
	"road2ca/internal/service"
)

type Handler struct {
	User    UserHandler
	Setting SettingHandler
	Collection CollectionHandler
	Ranking RankingHandler
}

func New(s *service.Services) *Handler {
	return &Handler{
		User:    NewUserHandler(s.User),
		Setting: NewSettingHandler(),
		Collection: NewCollectionHandler(s.Collection),
		Ranking: NewRankingHandler(s.Ranking),
	}
}
