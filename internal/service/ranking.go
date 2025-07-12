package service

import (
	"fmt"
	"road2ca/internal/repository"
)

type Ranking struct {
	UserID    int `json:"userId"`
	UserName  string `json:"userName"`
	Rank	  int    `json:"rank"`
	Score	 int `json:"score"`
}

type RankingService interface {
	// Create() error
	GetInRange(start, end int) ([]*Ranking, error)
}

type rankingService struct {
	rankingRepo repository.RankingRepository
	userRepo	repository.UserRepository
}

func NewRankingService(userRepo repository.UserRepository, rankingRepo repository.RankingRepository) RankingService {
	return &rankingService{
		rankingRepo: rankingRepo,
		userRepo:    userRepo,
	}
}

// GetInRange returns rankings in the specified range.
func (s *rankingService) GetInRange(start, end int) ([]*Ranking, error) {
	if start < 0 || end < 0 || start > end {
		return nil, fmt.Errorf("invalid range: start=%d, end=%d", start, end)
	}

	rankings, err := s.rankingRepo.FindInRangeFromCache(start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get rankings in range %d-%d: %w", start, end, err)
	}

	var result []*Ranking
	for _, r := range rankings {
		user, err := s.userRepo.FindByID(r.UserID)
		if err != nil {
			return nil, fmt.Errorf("failed to find user by ID %d: %w", r.UserID, err)
		}
		result = append(result, &Ranking{
			UserID:   r.UserID,
			UserName: user.Name,
			Rank:     r.Rank,
			Score:    r.Score,
		})
	}

	return result, nil
}
