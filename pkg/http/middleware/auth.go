package middleware

import (
	"context"
	"log"
	"net/http"
	"road2ca/internal/model"
	"road2ca/pkg/contextKey"
	"road2ca/pkg/server/minigin"
)

// Authenticate ユーザ認証を行ってContextへユーザID情報を保存する
func (m *Middleware) Authenticate(c *minigin.Context) {
	// 認証情報を取得
	token := c.Request.Header.Get("x-token")
	if len(token) < 1 {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Writer.Write([]byte(`{"error": "Unauthorized"}`))
		return
	}

	// ユーザ情報をDBから取得
	user, err := m.userDAO.GetByToken(token)
	if user == nil || err != nil {
		if err != nil {
			log.Println("Error getting user by token:", err)
		} else {
			log.Println("User not found for token:", token)
		}
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Writer.Write([]byte(`{"error": "Unauthorized"}`))
		return
	}

	// ユーザ情報をContextに保存
	ctx := c.Request.Context()
	ctx = setUserToContext(ctx, user)
	c.Request = c.Request.Clone(ctx)

	// 認証成功で次へ
	c.Next()
}

func setUserToContext(parents context.Context, user *model.User) context.Context {
	if user == nil {
		return parents
	}
	return context.WithValue(parents, contextkey.ContextKey, user)
}
