package service

import (
	"road2ca/internal/repository"
)

type ItemService interface {
	SetItemToCache() error
}

type itemService struct {
	mySqlItemRepo repository.MySQLItemRepo
	redisItemRepo repository.RedisItemRepo
}

func NewItemService(mySqlItemRepo repository.MySQLItemRepo, redisItemRepo repository.RedisItemRepo) ItemService {
	return &itemService{
		mySqlItemRepo: mySqlItemRepo,
		redisItemRepo: redisItemRepo,
	}
}

func (s *itemService) SetItemToCache() error {
	items, err := repository.FindItems(s.mySqlItemRepo, s.redisItemRepo)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}

	if err := s.redisItemRepo.Save(items); err != nil {
		return err
	}

	return nil
}
