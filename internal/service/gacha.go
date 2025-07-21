package service

import (
	"road2ca/internal/constants"
	"road2ca/internal/repository"
	"road2ca/pkg/minigin"
	"road2ca/internal/entity"
	"fmt"
	"math/rand"
	"time"
)

type GachaRequestBody struct {
	Times int `json:"times"` // ガチャを引く回数
}

type GachaResult struct {
	CollectionID int    `json:"collectionID"`
	Name         string `json:"name"`
	Rarity      int    `json:"rarity"`
	IsNew       bool   `json:"isNew"`
}

type GachaService interface {
	Draw(c *minigin.Context, times int) ([]GachaResult, error)
}

type gachaService struct {
	repo *repository.Repositories
	totalWeight int
	randGen *rand.Rand
}

func NewGachaService(repo *repository.Repositories) GachaService {
	items, err := repo.Item.FindAllFromCache()
	if err != nil {
		panic(fmt.Sprintf("failed to get items from cache: %v", err))
	}
	if len(items) == 0 {
		panic("no items found in cache")
	}
	totalWeight := 0
	for _, item := range items {
		if item.Weight == 0 {
			continue // 重みが0のアイテムは無視する
		} else if item.Weight < 0 {
			panic(fmt.Sprintf("item %d has invalid weight %d", item.ID, item.Weight))
		}
		totalWeight += item.Weight
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return &gachaService{
		repo: repo,
		totalWeight: totalWeight,
		randGen: r,
	}
}

func (s *gachaService) Draw(c *minigin.Context, times int) ([]GachaResult, error) {
	user, ok := c.Request.Context().Value(constants.ContextKey).(*entity.User)
	if !ok {
		return nil, fmt.Errorf("failed to get user from context")
	}

	const GachaCoinConsumption = 100 // 仮でハードコード
	if user.Coin < GachaCoinConsumption * times {
		return nil, fmt.Errorf("not enough coins")
	}

	// アイテムをキャッシュから取得
	items, err := s.repo.Item.FindAllFromCache()
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}

	// 重み付き抽選でアイテムを選ぶ
	var pickedItems []*entity.Item
	for i := 0; i < times; i++ {
		val := s.randGen.Intn(s.totalWeight)
		for _, item := range items {
			val -= item.Weight
			if val < 0 {
				pickedItems = append(pickedItems, item)
				break
			}
		}
	}

	collections, err := s.repo.Collection.FindByUserID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}

	var hasItemsMap = make(map[int]bool)
	for _, collection := range collections {
		hasItemsMap[collection.ItemID] = true
	}

	var insertNewCollections []*entity.Collection // 新規コレクションを格納するスライス
	var results []GachaResult // 結果を格納するスライス
	for _, item := range pickedItems {
		isNew := !hasItemsMap[item.ID]
		results = append(results, GachaResult{
			CollectionID: item.ID,
			Name:         item.Name,
			Rarity:       item.Rarity,
			IsNew:        isNew,
		})
		if isNew {
			insertNewCollections = append(insertNewCollections, &entity.Collection{
				UserID: user.ID,
				ItemID: item.ID,
			})
		}
	}

	// トランザクション開始
	tx, err := s.repo.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	err = s.repo.Collection.Save(tx, insertNewCollections)
	if err != nil {
		return nil, fmt.Errorf("failed to save collections: %w", err)
	}
	user.Coin -= GachaCoinConsumption * times
	err = s.repo.User.SaveTx(tx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return results, nil
}
