package middleware

import (
	"fmt"
	"net/http"
	"road2ca/internal/service"
	"road2ca/pkg/minigin"
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
		c.Error(fmt.Errorf("missing x-token header"))
		c.JSON(http.StatusUnauthorized, minigin.H{
			"error": "Unauthorized",
		})
		return
	}

	err := m.authService.SaveTokenToContext(token, c)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusUnauthorized, minigin.H{
			"error": "Unauthorized",
		})
		return
	}

	// 認証成功で次へ
	c.Next()
}
