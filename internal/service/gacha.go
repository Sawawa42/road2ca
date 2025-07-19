package service

import (
	"road2ca/internal/repository"
)

type Gacha struct {
	Times int `json:"times"` // ガチャを引く回数
}

type GachaService interface {
	// Draw(c *minigin.Context) (int, error)
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