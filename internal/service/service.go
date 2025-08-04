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
const ReqIdContextKey contextKeyType = "reqIdContextKey"

func New(repo *repository.Repositories, props *GachaProperties) *Services {
	return &Services{
		User:       NewUserService(repo.User),
		Auth:       NewAuthService(repo.User),
		Item:       NewItemService(),
		Collection: NewCollectionService(repo.Collection, repo.MySQLItem, repo.RedisItem),
		Ranking:    NewRankingService(repo.User, repo.Ranking, repo.MySQLSetting, repo.RedisSetting),
		Game:       NewGameService(repo.User, repo.Ranking, repo.MySQLSetting, repo.RedisSetting),
		Gacha:      NewGachaService(repo.MySQLItem, repo.RedisItem, repo.MySQLSetting, repo.RedisSetting, repo.Collection, repo.User, repo.DB, props),
		Setting:    NewSettingService(repo.MySQLSetting, repo.RedisSetting),
	}
}
