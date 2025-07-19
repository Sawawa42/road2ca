package service

import (
	"road2ca/internal/repository"
	"road2ca/pkg/minigin"
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
	itmRepo repository.ItemRepository
}

func NewGachaService(userRepo repository.UserRepository, itmRepo repository.ItemRepository) GachaService {
	return &gachaService{
		userRepo: userRepo,
		itmRepo: itmRepo,
	}
}

func (s *gachaService) Draw(c *minigin.Context, times int) ([]GachaResult, error) {
	// ガチャのロジックを実装
	// ここでは仮の実装として空のスライスを返す
	results := make([]GachaResult, times)
	for i := 0; i < times; i++ {
		results[i] = GachaResult{
			CollectionID: i + 1,
			Name:         "Item" + string(i+1),
			Rarity:      i % 5, // 仮のレアリティ
			IsNew:       true,  // 仮の新規フラグ
		}
	}
	return results, nil
}
