package service

import (
	"fmt"
	"road2ca/internal/entity"
	"road2ca/internal/repository"
	"road2ca/pkg/minigin"
	"github.com/google/uuid"
)

type CollectionListResponseDTO struct {
	Collections []*CollectionListItemDTO `json:"collections"`
}

type CollectionListItemDTO struct {
	CollectionID uuid.UUID `json:"collectionID"`
	Name         string    `json:"name"`
	Rarity       int       `json:"rarity"`
	HasItem      bool      `json:"hasItem"`
}

type CollectionService interface {
	GetCollectionList(c *minigin.Context) ([]*CollectionListItemDTO, error)
}

type collectionService struct {
	collectionRepo repository.CollectionRepo
	itemRepo       repository.ItemRepo
}

func NewCollectionService(collectionRepo repository.CollectionRepo, itemRepo repository.ItemRepo) CollectionService {
	return &collectionService{
		collectionRepo: collectionRepo,
		itemRepo:       itemRepo,
	}
}

func (s *collectionService) GetCollectionList(c *minigin.Context) ([]*CollectionListItemDTO, error) {
	// ユーザー情報をコンテキストから取得
	user, ok := c.Request.Context().Value(ContextKey).(*entity.User)
	if !ok {
		return nil, fmt.Errorf("failed to get user from context")
	}

	// アイテムを取得
	items, err := s.itemRepo.Find()
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}
	// ユーザーのコレクションを取得
	collections, err := s.collectionRepo.FindByUserID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}

	// 特定のアイテムIDがユーザのコレクションに含まれているかをチェックするためのマップを作成
	collectionItemMap := make(map[uuid.UUID]bool)
	for _, collection := range collections {
		collectionItemMap[collection.ItemID] = true
	}

	res := make([]*CollectionListItemDTO, 0, len(items))
	for _, item := range items {
		// アイテム所持を判定
		hasItem := collectionItemMap[item.ID]
		res = append(res, &CollectionListItemDTO{
			CollectionID: item.ID,
			Name:         item.Name,
			Rarity:       item.Rarity,
			HasItem:      hasItem,
		})
	}

	return res, nil
}
