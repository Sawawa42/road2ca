package service

import (
	"fmt"
	"road2ca/internal/constants"
	"road2ca/internal/entity"
	"road2ca/internal/repository"
	"road2ca/pkg/minigin"
)

type CollectionsResponse struct {
	CollectionID int    `json:"collectionID"`
	Name         string `json:"name"`
	Rarity       int    `json:"rarity"`
	HasItem      bool   `json:"hasItem"`
}

type CollectionService interface {
	GetCollectionList(c *minigin.Context) ([]*CollectionsResponse, error)
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

func (s *collectionService) GetCollectionList(c *minigin.Context) ([]*CollectionsResponse, error) {
	// ユーザー情報をコンテキストから取得
	user, ok := c.Request.Context().Value(constants.ContextKey).(*entity.User)
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
	collectionItemMap := make(map[int]bool)
	for _, collection := range collections {
		collectionItemMap[collection.ItemID] = true
	}

	res := make([]*CollectionsResponse, 0, len(items))
	for _, item := range items {
		// アイテム所持を判定
		hasItem := collectionItemMap[item.ID]
		res = append(res, &CollectionsResponse{
			CollectionID: item.ID,
			Name:         item.Name,
			Rarity:       item.Rarity,
			HasItem:      hasItem,
		})
	}

	return res, nil
}
