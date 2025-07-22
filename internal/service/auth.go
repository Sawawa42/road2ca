package service

import (
	"context"
	"fmt"
	"road2ca/internal/repository"
	"road2ca/pkg/minigin"
)

type AuthService interface {
	SaveTokenToContext(token string, c *minigin.Context) error
}

type authService struct {
	userRepo repository.UserRepo
}

func NewAuthService(userRepo repository.UserRepo) AuthService {
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
	ctx = context.WithValue(ctx, ContextKey, user)
	c.Request = c.Request.Clone(ctx)

	return nil
}
