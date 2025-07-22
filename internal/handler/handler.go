package handler

import (
	"road2ca/internal/service"
)

type Handler struct {
	User       UserHandler
	Setting    SettingHandler
	Collection CollectionHandler
	Ranking    RankingHandler
	Game       GameHandler
	Gacha      GachaHandler
}

func New(s *service.Services) *Handler {
	return &Handler{
		User:       NewUserHandler(s.User),
		Setting:    NewSettingHandler(s.Setting),
		Collection: NewCollectionHandler(s.Collection),
		Ranking:    NewRankingHandler(s.Ranking),
		Game:       NewGameHandler(s.User, s.Game),
		Gacha:      NewGachaHandler(s.Gacha),
	}
}
