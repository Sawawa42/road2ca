package service

import (
	"fmt"
	"road2ca/internal/repository"
)

type ItemService interface {
}

type itemService struct {
	itemRepo repository.ItemRepository
}

func NewItemService(itemRepo repository.ItemRepository) ItemService {
	if err := setToCache(itemRepo); err != nil {
		panic(fmt.Sprintf("failed to set items to cache: %v", err))
	}

	return &itemService{
		itemRepo: itemRepo,
	}
}

func setToCache(itemRepo repository.ItemRepository) error {
	items, err := itemRepo.FindAllFromDB()
	if err != nil {
		return fmt.Errorf("failed to find items from MySQL: %w", err)
	}
	if len(items) == 0 {
		return fmt.Errorf("no items found in MySQL")
	}

	if err := itemRepo.SaveToCache(items); err != nil {
		return fmt.Errorf("failed to cache items to Redis: %w", err)
	}

	return nil
}
