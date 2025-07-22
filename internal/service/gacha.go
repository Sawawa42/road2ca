package service

import (
	"database/sql"
	"fmt"
	"math/rand"
	"road2ca/internal/constants"
	"road2ca/internal/entity"
	"road2ca/internal/repository"
	"road2ca/pkg/minigin"
)

type DrawGachaRequestDTO struct {
	Times int `json:"times"` // ガチャを引く回数
}

type DrawGachaResponseDTO struct {
	Results []GachaItemDTO `json:"results"`
}

type GachaItemDTO struct {
	CollectionID int    `json:"collectionID"`
	Name         string `json:"name"`
	Rarity       int    `json:"rarity"`
	IsNew        bool   `json:"isNew"`
}

type GachaServiceProps struct {
	TotalWeight int
	RandGen     *rand.Rand
}

type GachaService interface {
	DrawGacha(c *minigin.Context, times int) (*DrawGachaResponseDTO, error)
}

type gachaService struct {
	itemRepo       repository.ItemRepo
	collectionRepo repository.CollectionRepo
	userRepo       repository.UserRepo
	db             *sql.DB
	totalWeight    int
	randGen        *rand.Rand
}

func NewGachaService(
	itemRepo repository.ItemRepo,
	collectionRepo repository.CollectionRepo,
	userRepo repository.UserRepo,
	db *sql.DB,
	gachaProps *GachaServiceProps,
) GachaService {
	return &gachaService{
		itemRepo:       itemRepo,
		collectionRepo: collectionRepo,
		userRepo:       userRepo,
		db:             db,
		totalWeight:    gachaProps.TotalWeight,
		randGen:        gachaProps.RandGen,
	}
}

func (s *gachaService) DrawGacha(c *minigin.Context, times int) (*DrawGachaResponseDTO, error) {
	user, ok := c.Request.Context().Value(constants.ContextKey).(*entity.User)
	if !ok {
		return nil, fmt.Errorf("failed to get user from context")
	}

	// TODO: ガチャの消費コイン数は設定から取得する
	const GachaCoinConsumption = 100
	if user.Coin < GachaCoinConsumption*times {
		return nil, fmt.Errorf("not enough coins")
	}

	// アイテムを取得
	items, err := s.itemRepo.Find()
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

	collections, err := s.collectionRepo.FindByUserID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}

	var hasItemsMap = make(map[int]bool)
	for _, collection := range collections {
		hasItemsMap[collection.ItemID] = true
	}

	var insertNewCollections []*entity.Collection // 新規コレクションを格納するスライス
	var results []GachaItemDTO                    // 結果を格納するスライス
	for _, item := range pickedItems {
		isNew := !hasItemsMap[item.ID]
		results = append(results, GachaItemDTO{
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
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	err = s.collectionRepo.Save(tx, insertNewCollections)
	if err != nil {
		return nil, fmt.Errorf("failed to save collections: %w", err)
	}
	user.Coin -= GachaCoinConsumption * times
	err = s.userRepo.SaveTx(tx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &DrawGachaResponseDTO{
		Results: results,
	}, nil
}
