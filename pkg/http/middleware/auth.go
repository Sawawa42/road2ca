package middleware

import (
	"net/http"
	"road2ca/pkg/server/minigin"
	"log"
	"context"
)

// Authenticate ユーザ認証を行ってContextへユーザID情報を保存する
func (m *Middleware) Authenticate(c *minigin.Context) {
	// 認証情報を取得
	token := c.Request.Header.Get("x-token")
	if (len(token) < 1) {
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Writer.Write([]byte(`{"error": "Unauthorized"}`))
		return
	}

	// ユーザ情報をDBから取得
	user, err := m.userDAO.GetByToken(token)
	if (user == nil || err != nil) {
		if (err != nil) {
			log.Println("Error getting user by token:", err)
		} else {
			log.Println("User not found for token:", token)
		}
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Writer.Write([]byte(`{"error": "Unauthorized"}`))
		return
	}

	// ユーザIDをContextに保存
	type contextKey string
	const tokenKey contextKey = "token" // should not use built-in type string as key for value; define your own type to avoid collisions (SA1029)go-staticcheck
	ctx := context.WithValue(c.Request.Context(), tokenKey, user.Token)
	c.Request = c.Request.WithContext(ctx)

	// 認証成功で次へ
	c.Next()
}
