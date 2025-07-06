package service

import (
	"road2ca/internal/entity"
	"road2ca/internal/repository"
	"crypto/md5"
	"fmt"
	"road2ca/pkg/minigin"
	"road2ca/internal/constants"
)

type UserDTO struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	HighScore int    `json:"highScore"`
	Coin      int    `json:"coin"`
}

type UserCreateResponseDTO struct {
	Token string `json:"token"`
}

type UserCreateRequestDTO struct {
	Name string `json:"name"`
}

type UserService interface {
	CreateUser(name string) (*UserCreateResponseDTO, error)
	GetUser(c *minigin.Context) (*UserDTO, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(name string) (*UserCreateResponseDTO, error) {
	token := fmt.Sprintf("%x", md5.Sum([]byte(name)))
	user := &entity.User{
		Name:  name,
		HighScore: 0,
		Coin: 0,
		Token: token,
	}
	if err := s.userRepo.Save(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &UserCreateResponseDTO{
		Token: token,
	}, nil
}

func (s *userService) GetUser(c *minigin.Context) (*UserDTO, error) {
	user, ok := c.Request.Context().Value(constants.ContextKey).(*entity.User)
	if !ok {
		return nil, fmt.Errorf("failed to get user")
	}

	return &UserDTO{
		ID:        user.ID,
		Name:      user.Name,
		HighScore: user.HighScore,
		Coin:      user.Coin,
	}, nil
}
