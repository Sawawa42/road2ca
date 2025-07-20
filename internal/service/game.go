package service

import (
	"fmt"
	"road2ca/internal/constants"
	"road2ca/internal/entity"
	"road2ca/internal/repository"
	"road2ca/pkg/minigin"
)

type GameFinishRequestDTO struct {
	Score int `json:"score"`
}

type GameService interface {
	Finish(c *minigin.Context, score int) (int, error)
}

type gameService struct {
	userRepo    repository.UserRepository
	rankingRepo repository.RankingRepository
}

func NewGameService(userRepo repository.UserRepository, rankingRepo repository.RankingRepository) GameService {
	return &gameService{
		userRepo:    userRepo,
		rankingRepo: rankingRepo,
	}
}

func (s *gameService) Finish(c *minigin.Context, score int) (int, error) {
	user, ok := c.Request.Context().Value(constants.ContextKey).(*entity.User)
	if !ok {
		return 0, fmt.Errorf("failed to get user")
	}

	if score > user.HighScore {
		user.HighScore = score
	}
	user.Coin += 100

	if err := s.userRepo.Save(nil, user); err != nil {
		return 0, err
	}

	if err := s.rankingRepo.SaveToCache(user); err != nil {
		return 0, err
	}

	return user.Coin, nil
}
