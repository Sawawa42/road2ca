package service

import (
	"crypto/md5"
	"fmt"
	"road2ca/internal/constants"
	"road2ca/internal/entity"
	"road2ca/internal/repository"
	"road2ca/pkg/minigin"
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

type UserUpdateRequestDTO struct {
	Name      string `json:"name"`
}

type UserService interface {
	Create(name string) (*UserCreateResponseDTO, error)
	Get(c *minigin.Context) (*UserDTO, error)
	Update(c *minigin.Context, name string) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Create(name string) (*UserCreateResponseDTO, error) {
	token := fmt.Sprintf("%x", md5.Sum([]byte(name)))
	user := &entity.User{
		Name:      name,
		HighScore: 0,
		Coin:      0,
		Token:     token,
	}
	if err := s.userRepo.Save(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &UserCreateResponseDTO{
		Token: token,
	}, nil
}

func (s *userService) Get(c *minigin.Context) (*UserDTO, error) {
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

func (s *userService) Update(c *minigin.Context, name string) error {
	user, ok := c.Request.Context().Value(constants.ContextKey).(*entity.User)
	if !ok {
		return fmt.Errorf("failed to get user")
	}

	user.Name = name
	if err := s.userRepo.Save(user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}
	
