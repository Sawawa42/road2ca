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

type GameFinishResponseDTO struct {
	Coin int `json:"coin"`
}

type GameService interface {
	FinalizeGame(c *minigin.Context, score int) (*GameFinishResponseDTO, error)
}

type gameService struct {
	userRepo    repository.UserRepo
	rankingRepo repository.RankingRepo
}

func NewGameService(userRepo repository.UserRepo, rankingRepo repository.RankingRepo) GameService {
	return &gameService{
		userRepo:    userRepo,
		rankingRepo: rankingRepo,
	}
}

func (s *gameService) FinalizeGame(c *minigin.Context, score int) (*GameFinishResponseDTO, error) {
	user, ok := c.Request.Context().Value(constants.ContextKey).(*entity.User)
	if !ok {
		return nil, fmt.Errorf("failed to get user")
	}

	if score > user.HighScore {
		user.HighScore = score
	}
	user.Coin += 100

	if err := s.userRepo.Save(user); err != nil {
		return nil, err
	}

	if err := s.rankingRepo.Save(user); err != nil {
		return nil, err
	}

	return &GameFinishResponseDTO{
		Coin: user.Coin,
	}, nil
}
