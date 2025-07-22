package service

import (
	"road2ca/internal/repository"
)

type Services struct {
	User       UserService
	Auth       AuthService
	Setting    SettingService
	Item       ItemService
	Collection CollectionService
	Ranking    RankingService
	Game       GameService
	Gacha      GachaService
}

type contextKeyType string

const ContextKey contextKeyType = "contextKey"

func New(repo *repository.Repositories, gachaProps *GachaServiceProps) *Services {
	return &Services{
		User:       NewUserService(repo.User),
		Auth:       NewAuthService(repo.User),
		Item:       NewItemService(repo.Item),
		Collection: NewCollectionService(repo.Collection, repo.Item),
		Ranking:    NewRankingService(repo.User, repo.Ranking, repo.Setting),
		Game:       NewGameService(repo.User, repo.Ranking, repo.Setting),
		Gacha:      NewGachaService(repo.Item, repo.Collection, repo.User, repo.Setting, repo.DB, gachaProps),
		Setting:    NewSettingService(repo.Setting),
	}
}
