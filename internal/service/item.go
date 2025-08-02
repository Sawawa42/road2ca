package service

import (
	"fmt"
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
		return fmt.Errorf("no items found in both MySQL and Redis")
	}

	if err := s.redisItemRepo.Save(items); err != nil {
		return err
	}

	return nil
}
