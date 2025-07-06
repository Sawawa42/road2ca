package service

import (
	"road2ca/internal/repository"
	"fmt"
	"road2ca/pkg/minigin"
	"road2ca/internal/constants"
	"context"
)

type AuthService interface {
	SaveTokenToContext(token string, c *minigin.Context) error
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) SaveTokenToContext(token string, c *minigin.Context) error {
	user, err := s.userRepo.FindByToken(token)
	if err != nil {
		return fmt.Errorf("internal server error: %w", err)
	}

	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, constants.ContextKey, user)
	c.Request = c.Request.Clone(ctx)

	return nil
}
