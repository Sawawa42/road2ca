package middleware

import (
	"road2ca/internal/service"
)

type Middleware struct {
	Auth   AuthMiddleware
	Cors   CorsMiddleware
	Logger LoggerMiddleware
}

func New(s *service.Services, slogs *SlogInstances) *Middleware {
	return &Middleware{
		Auth:   NewAuthMiddleware(s.Auth),
		Cors:   NewCorsMiddleware(),
		Logger: NewLoggerMiddleware(slogs),
	}
}
