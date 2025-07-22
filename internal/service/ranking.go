package service

import (
	"fmt"
	"road2ca/internal/repository"
)

type GetRankingListResponseDTO struct {
	Rankings []*RankingItemDTO `json:"ranks"`
}

type RankingItemDTO struct {
	UserID   int    `json:"userId"`
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
	settingRepo repository.SettingRepo
}

func NewRankingService(userRepo repository.UserRepo, rankingRepo repository.RankingRepo, settingRepo repository.SettingRepo) RankingService {
	return &rankingService{
		rankingRepo: rankingRepo,
		userRepo:    userRepo,
		settingRepo: settingRepo,
	}
}

func (s *rankingService) GetRanking(start int) ([]*RankingItemDTO, error) {
	setting, err := s.settingRepo.FindLatest()
	if err != nil {
		return nil, err
	}

	end := setting.GetRankingLimit
	if start < 0 || end < 0 || start > end {
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
		result = append(result, &RankingItemDTO{
			UserID:   r.UserID,
			UserName: user.Name,
			Rank:     r.Rank,
			Score:    r.Score,
		})
	}

	return result, nil
}
