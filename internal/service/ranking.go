package service

import (
	"fmt"
	"road2ca/internal/repository"
	"road2ca/internal/entity"
	"github.com/google/uuid"
)

type GetRankingListResponseDTO struct {
	Rankings []*RankingItemDTO `json:"ranks"`
}

type RankingItemDTO struct {
	UserID   string `json:"userId"`
	UserName string `json:"userName"`
	Rank     int    `json:"rank"`
	Score    int    `json:"score"`
}

type RankingService interface {
	GetRanking(start int) ([]*RankingItemDTO, error)
}

type rankingService struct {
	rankingRepo      repository.RankingRepo
	userRepo         repository.UserRepo
	mysqlSettingRepo repository.MySQLSettingRepo
	redisSettingRepo repository.RedisSettingRepo
}

func NewRankingService(userRepo repository.UserRepo, rankingRepo repository.RankingRepo, mysqlSettingRepo repository.MySQLSettingRepo, redisSettingRepo repository.RedisSettingRepo) RankingService {
	return &rankingService{
		rankingRepo:      rankingRepo,
		userRepo:         userRepo,
		mysqlSettingRepo: mysqlSettingRepo,
		redisSettingRepo: redisSettingRepo,
	}
}

func (s *rankingService) GetRanking(start int) ([]*RankingItemDTO, error) {
	setting, err := repository.FindSetting(s.mysqlSettingRepo, s.redisSettingRepo)
	if err != nil {
		return nil, err
	}

	// start 位から GetRankingLimit 件のランキングを取得
	end := start + setting.GetRankingLimit - 1
	if start < 0 || setting.GetRankingLimit < 0 || start >= end {
		return nil, fmt.Errorf("invalid range: start=%d, end=%d", start, end)
	}

	rankings, err := s.rankingRepo.FindInRange(start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get rankings in range %d-%d: %w", start, end, err)
	}

	if len(rankings) == 0 {
		return []*RankingItemDTO{}, nil
	}

	userIDs := make([][]byte, 0, len(rankings))
	for _, r := range rankings {
		userIDs = append(userIDs, r.UserID)
	}

	// n+1対策で、取得したランキングのユーザーIDを一度に取得
	users, err := s.userRepo.FindByIDs(userIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to find users by IDs: %w", err)
	}

	userMap := make(map[uuid.UUID]*entity.User, len(users))
	for _, user := range users {
		uuid, err := uuid.FromBytes(user.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse user ID: %w", err)
		}
		userMap[uuid] = user
	}

	var result []*RankingItemDTO
	for _, r := range rankings {
		uuid, err := uuid.FromBytes(r.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse user ID: %w", err)
		}
		user, ok := userMap[uuid]
		if !ok {
			return nil, fmt.Errorf("user with ID %s not found in user map", uuid.String())
		}

		result = append(result, &RankingItemDTO{
			UserID:   uuid.String(),
			UserName: user.Name,
			Rank:     r.Rank,
			Score:    r.Score,
		})
	}

	return result, nil
}
