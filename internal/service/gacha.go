package service

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"road2ca/internal/entity"
	"road2ca/internal/repository"
	"road2ca/pkg/minigin"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type DrawGachaRequestDTO struct {
	Times int `json:"times"` // ガチャを引く回数
}

type DrawGachaResponseDTO struct {
	Results []GachaItemDTO `json:"results"`
}

type GachaItemDTO struct {
	CollectionID string `json:"collectionID"`
	Name         string `json:"name"`
	Rarity       int    `json:"rarity"`
	IsNew        bool   `json:"isNew"`
}

type GachaServiceProps struct {
	TotalWeight int
	RandGen     *rand.Rand
}

var ErrNotEnoughCoins = errors.New("not enough coins")

type GachaService interface {
	DrawGacha(c *minigin.Context, times int) (*DrawGachaResponseDTO, error)
	SetGachaProps(props *GachaServiceProps)
}

type gachaService struct {
	mysqlItemRepo    repository.MySQLItemRepo
	redisItemRepo    repository.RedisItemRepo
	mysqlSettingRepo repository.MySQLSettingRepo
	redisSettingRepo repository.RedisSettingRepo
	collectionRepo   repository.CollectionRepo
	userRepo         repository.UserRepo
	db               *sql.DB
	totalWeight      int
	randGen          *rand.Rand
}

func NewGachaService(
	mysqlItemRepo repository.MySQLItemRepo,
	redisItemRepo repository.RedisItemRepo,
	mysqlSettingRepo repository.MySQLSettingRepo,
	redisSettingRepo repository.RedisSettingRepo,
	collectionRepo repository.CollectionRepo,
	userRepo repository.UserRepo,
	db *sql.DB,
) GachaService {
	return &gachaService{
		mysqlItemRepo:    mysqlItemRepo,
		redisItemRepo:    redisItemRepo,
		mysqlSettingRepo: mysqlSettingRepo,
		redisSettingRepo: redisSettingRepo,
		collectionRepo:   collectionRepo,
		userRepo:         userRepo,
		db:               db,
	}
}

func (s *gachaService) DrawGacha(c *minigin.Context, times int) (*DrawGachaResponseDTO, error) {
	setting, err := repository.FindSetting(s.mysqlSettingRepo, s.redisSettingRepo)
	if err != nil {
		return nil, err
	}

	if times < 1 || times > setting.DrawGachaMaxTimes {
		return nil, fmt.Errorf("times must be between 1 and %d", setting.DrawGachaMaxTimes)
	}

	user, ok := c.Request.Context().Value(ContextKey).(*entity.User)
	if !ok {
		return nil, fmt.Errorf("failed to get user from context")
	}

	if user.Coin < setting.GachaCoinConsumption*times {
		return nil, ErrNotEnoughCoins
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

	var hasItemsMap = make(map[uuid.UUID]bool)
	for _, collection := range collections {
		uuid, err := uuid.FromBytes(collection.ItemID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse item ID: %w", err)
		}
		hasItemsMap[uuid] = true
	}

	var insertNewCollections []*entity.Collection // 新規コレクションを格納するスライス
	var results []GachaItemDTO                    // 結果を格納するスライス
	for _, item := range pickedItems {
		uuid, err := uuid.FromBytes(item.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse item ID: %w", err)
		}

		isNew := !hasItemsMap[uuid]
		results = append(results, GachaItemDTO{
			CollectionID: uuid.String(),
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
	user.Coin -= setting.GachaCoinConsumption * times
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

func (s *gachaService) SetGachaProps(props *GachaServiceProps) {
	s.totalWeight = props.TotalWeight
	s.randGen = props.RandGen
}
