package service

import (
	"fmt"
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
	userRepo         repository.UserRepo
	rankingRepo      repository.RankingRepo
	mysqlSettingRepo repository.MySQLSettingRepo
	redisSettingRepo repository.RedisSettingRepo
}

func NewGameService(userRepo repository.UserRepo, rankingRepo repository.RankingRepo, mysqlSettingRepo repository.MySQLSettingRepo, redisSettingRepo repository.RedisSettingRepo) GameService {
	return &gameService{
		userRepo:         userRepo,
		rankingRepo:      rankingRepo,
		mysqlSettingRepo: mysqlSettingRepo,
		redisSettingRepo: redisSettingRepo,
	}
}

func (s *gameService) FinalizeGame(c *minigin.Context, score int) (*GameFinishResponseDTO, error) {
	user, ok := c.Request.Context().Value(ContextKey).(*entity.User)
	if !ok {
		return nil, fmt.Errorf("failed to get user")
	}

	if score > user.HighScore {
		user.HighScore = score
	}

	setting, err := repository.FindSetting(s.mysqlSettingRepo, s.redisSettingRepo)
	if err != nil {
		return nil, err
	}

	user.Coin += setting.RewardCoin

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
