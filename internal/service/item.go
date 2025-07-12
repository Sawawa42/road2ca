package service

import (
	"fmt"
	"road2ca/internal/repository"
)

type ItemService interface {
	CacheItems() error
}

type itemService struct {
	itemRepo repository.ItemRepository
}

func NewItemService(itemRepo repository.ItemRepository) ItemService {
	return &itemService{
		itemRepo: itemRepo,
	}
}

func (s *itemService) CacheItems() error {
	items, err := s.itemRepo.FindAllItemsFromDB()
	if err != nil {
		return fmt.Errorf("failed to find items from MySQL: %w", err)
	}
	if len(items) == 0 {
		return fmt.Errorf("no items found in MySQL")
	}

	if err := s.itemRepo.SaveItemsToCache(items); err != nil {
		return fmt.Errorf("failed to cache items to Redis: %w", err)
	}

	return nil
}
