package middleware

import (
	"log"
	"net/http"
	"road2ca/pkg/minigin"

	"road2ca/internal/service"
)

type AuthMiddleware interface {
	Authenticate(c *minigin.Context)
}

type authMiddleware struct {
	authService service.AuthService
}

func NewAuthMiddleware(authService service.AuthService) AuthMiddleware {
	return &authMiddleware{
		authService: authService,
	}
}

// Authenticate ユーザ認証を行ってContextへユーザID情報を保存する
func (m *authMiddleware) Authenticate(c *minigin.Context) {
	// 認証情報を取得
	token := c.Request.Header.Get("x-token")
	if len(token) < 1 {
		http.Error(c.Writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := m.authService.SaveTokenToContext(token, c)
	if err != nil {
		log.Printf("Authentication failed: %v", err)
		http.Error(c.Writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 認証成功で次へ
	c.Next()
}
