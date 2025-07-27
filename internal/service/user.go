package service

import (
	"crypto/md5"
	"fmt"
	"road2ca/internal/entity"
	"road2ca/internal/repository"
	"road2ca/pkg/minigin"

	"github.com/google/uuid"
	"time"
)

type GetUserResponseDTO struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	HighScore int    `json:"highScore"`
	Coin      int    `json:"coin"`
}

type CreateUserResponseDTO struct {
	Token string `json:"token"`
}

type CreateUserRequestDTO struct {
	Name string `json:"name"`
}

type UpdateUserRequestDTO struct {
	Name string `json:"name"`
}

type UserService interface {
	CreateUser(name string) (*CreateUserResponseDTO, error)
	GetUser(c *minigin.Context) (*GetUserResponseDTO, error)
	UpdateUser(c *minigin.Context, name string) error
}

type userService struct {
	userRepo repository.UserRepo
}

func NewUserService(userRepo repository.UserRepo) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(name string) (*CreateUserResponseDTO, error) {
	token := fmt.Sprintf("%x", md5.Sum([]byte(name + fmt.Sprintf("%d", time.Now().UnixNano()))))
	uuidBytes, err := repository.GetUUIDv7Bytes()
	if err != nil {
		return nil, err
	}
	user := &entity.User{
		ID:        uuidBytes,
		Name:      name,
		HighScore: 0,
		Coin:      0,
		Token:     token,
	}
	if err := s.userRepo.Save(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &CreateUserResponseDTO{
		Token: token,
	}, nil
}

func (s *userService) GetUser(c *minigin.Context) (*GetUserResponseDTO, error) {
	user, ok := c.Request.Context().Value(ContextKey).(*entity.User)
	if !ok {
		return nil, fmt.Errorf("failed to get user")
	}

	uuid, err := uuid.FromBytes(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user ID: %w", err)
	}

	return &GetUserResponseDTO{
		ID:        uuid.String(),
		Name:      user.Name,
		HighScore: user.HighScore,
		Coin:      user.Coin,
	}, nil
}

func (s *userService) UpdateUser(c *minigin.Context, name string) error {
	user, ok := c.Request.Context().Value(ContextKey).(*entity.User)
	if !ok {
		return fmt.Errorf("failed to get user")
	}

	user.Name = name
	if err := s.userRepo.Save(user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}
