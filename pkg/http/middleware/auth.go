package middleware

import (
	"net/http"
	"road2ca/pkg/server/minigin"
)

// Authenticate ユーザ認証を行ってContextへユーザID情報を保存する
func Authenticate(c *minigin.Context) {
	// 認証情報を取得
	token := c.Request.Header.Get("x-token")
	if len(token) < 32 { // TODO: 条件を検討する
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Writer.Write([]byte(`{"error": "Unauthorized"}`))
		return
	}

	// TODO: implement authentication logic here

	c.Next()
}
