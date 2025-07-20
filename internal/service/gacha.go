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
	userRepo repository.UserRepository
	itemRepo repository.ItemRepository
	collectionRepo repository.CollectionRepository
	totalWeight int
	randGen *rand.Rand
}

func NewGachaService(userRepo repository.UserRepository, itemRepo repository.ItemRepository, collectionRepo repository.CollectionRepository) GachaService {
	items, err := itemRepo.FindAllFromCache()
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
		userRepo: userRepo,
		itemRepo: itemRepo,
		collectionRepo: collectionRepo,
		totalWeight: totalWeight,
		randGen: r,
	}
}

func (s *gachaService) Draw(c *minigin.Context, times int) ([]GachaResult, error) {
	// やること
	// 5. 抽選結果をユーザのコレクションに追加
	// 6. ユーザのコインを減らす
	// (5, 6はトランザクションでまとめる)
	user, ok := c.Request.Context().Value(constants.ContextKey).(*entity.User)
	if !ok {
		return nil, fmt.Errorf("failed to get user from context")
	}

	const GachaCoinConsumption = 100 // 仮でハードコード
	if user.Coin < GachaCoinConsumption * times {
		return nil, fmt.Errorf("not enough coins")
	}

	// アイテムをキャッシュから取得
	items, err := s.itemRepo.FindAllFromCache()
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

	// ユーザの現在のコレクションを取得
	// pickedItemsと比較して新規アイテムを判定
	// コレクションに追加するアイテムを決定
	// pickedItemsのうち、新規アイテムをisNew=trueに設定したものをresultsとして返す
	collections, err := s.collectionRepo.FindAllByUserID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}

	var newItemsMap = make(map[int]bool)
	for _, collection := range collections {
		newItemsMap[collection.ItemID] = true
	}
	var insertItems []*entity.Item // 新規アイテムを格納するスライス
	var results []GachaResult // 結果を格納するスライス
	for _, item := range pickedItems {
		isNew := !newItemsMap[item.ID]
		results = append(results, GachaResult{
			CollectionID: item.ID,
			Name:         item.Name,
			Rarity:       item.Rarity,
			IsNew:        isNew,
		})
		if isNew {
			insertItems = append(insertItems, item)
		}
	}
	
	return results, nil
}
