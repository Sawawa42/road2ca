package service

import (
	"fmt"
	"road2ca/internal/entity"
	"road2ca/internal/repository"
	"road2ca/pkg/minigin"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type CollectionListResponseDTO struct {
	Collections []*CollectionListItemDTO `json:"collections"`
}

type CollectionListItemDTO struct {
	CollectionID string `json:"collectionID"`
	Name         string `json:"name"`
	Rarity       int    `json:"rarity"`
	HasItem      bool   `json:"hasItem"`
}

type CollectionService interface {
	GetCollectionList(c *minigin.Context) ([]*CollectionListItemDTO, error)
}

type collectionService struct {
	collectionRepo repository.CollectionRepo
	mysqlItemRepo  repository.MySQLItemRepo
	redisItemRepo  repository.RedisItemRepo
}

func NewCollectionService(collectionRepo repository.CollectionRepo, mysqlItemRepo repository.MySQLItemRepo, redisItemRepo repository.RedisItemRepo) CollectionService {
	return &collectionService{
		collectionRepo: collectionRepo,
		mysqlItemRepo:  mysqlItemRepo,
		redisItemRepo:  redisItemRepo,
	}
}

func (s *collectionService) GetCollectionList(c *minigin.Context) ([]*CollectionListItemDTO, error) {
	// ユーザー情報をコンテキストから取得
	user, ok := c.Request.Context().Value(ContextKey).(*entity.User)
	if !ok {
		return nil, fmt.Errorf("failed to get user from context")
	}

	// アイテムをキャッシュから取得
	items, err := s.redisItemRepo.Find()
	if err != nil {
		// キャッシュにアイテムがない場合はMySQLから取得
		if err == redis.Nil {
			items, err = s.mysqlItemRepo.Find()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// ユーザーのコレクションを取得
	collections, err := s.collectionRepo.FindByUserID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}

	res := make([]*CollectionListItemDTO, 0, len(items))
	for _, collection := range collections {
		uuid, err := uuid.FromBytes(collection.ID)
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			hasItem := false
			if item.ID == collection.ItemID {
				hasItem = true
			}
			res = append(res, &CollectionListItemDTO{
				CollectionID: uuid.String(),
				Name:         item.Name,
				Rarity:       item.Rarity,
				HasItem:      hasItem,
			})
			break
		}
	}

	return res, nil
}
