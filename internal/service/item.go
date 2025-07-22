package service

import (
	"fmt"
	"road2ca/internal/repository"
)

type ItemService interface {
	SetItemToCache() error
}

type itemService struct {
	itemRepo repository.ItemRepo
}

func NewItemService(itemRepo repository.ItemRepo) ItemService {
	return &itemService{
		itemRepo: itemRepo,
	}
}

func (s *itemService) SetItemToCache() error {
	items, err := s.itemRepo.Find()
	if err != nil {
		return fmt.Errorf("failed to find items from MySQL: %w", err)
	}
	if len(items) == 0 {
		return fmt.Errorf("no items found in MySQL")
	}

	if err := s.itemRepo.Save(items); err != nil {
		return fmt.Errorf("failed to cache items to Redis: %w", err)
	}

	return nil
}
