package middleware

import (
	"road2ca/internal/service"
)

type Middleware struct {
	Auth AuthMiddleware
	Cors CorsMiddleware
}

func New(s *service.Services) *Middleware {
	return &Middleware{
		Auth: NewAuthMiddleware(s.Auth),
		Cors: NewCorsMiddleware(),
	}
}
