package service

import (
	"fmt"
	"road2ca/internal/repository"

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
	rankingRepo repository.RankingRepo
	userRepo    repository.UserRepo
	mysqlSettingRepo repository.MySQLSettingRepo
	redisSettingRepo repository.RedisSettingRepo
}

func NewRankingService(userRepo repository.UserRepo, rankingRepo repository.RankingRepo, mysqlSettingRepo repository.MySQLSettingRepo, redisSettingRepo repository.RedisSettingRepo) RankingService {
	return &rankingService{
		rankingRepo:   rankingRepo,
		userRepo:      userRepo,
		mysqlSettingRepo: mysqlSettingRepo,
		redisSettingRepo: redisSettingRepo,
	}
}

func (s *rankingService) GetRanking(start int) ([]*RankingItemDTO, error) {
	setting, err := repository.FindSetting(s.mysqlSettingRepo, s.redisSettingRepo)
	if err != nil {
		return nil, err
	}

	end := start + setting.GetRankingLimit - 1
	if start < 0 || setting.GetRankingLimit < 0 || start >= end {
		return nil, fmt.Errorf("invalid range: start=%d, end=%d", start, end)
	}

	rankings, err := s.rankingRepo.FindInRange(start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get rankings in range %d-%d: %w", start, end, err)
	}

	var result []*RankingItemDTO
	for _, r := range rankings {
		user, err := s.userRepo.FindByID(r.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to find user by ID %d: %w", r.UserID, err)
		}
		uuid, err := uuid.FromBytes(user.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse user ID: %w", err)
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
