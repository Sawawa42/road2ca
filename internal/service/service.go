package service

import (
	"road2ca/internal/repository"
)

type Services struct {
	User    UserService
	Auth    AuthService
	Setting SettingService
	Item    ItemService
	Collection CollectionService
}

func New(repo *repository.Repositories) *Services {
	return &Services{
		User:    NewUserService(repo.User),
		Auth:    NewAuthService(repo.User),
		Setting: NewSettingService(),
		Item:    NewItemService(repo.Item),
		Collection: NewCollectionService(repo.Collection, repo.Item),
	}
}
