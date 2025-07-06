package middleware

import (
	"net/http"
	"road2ca/pkg/minigin"
)

type CorsMiddleware interface {
	SettingCors(c *minigin.Context)
}

type corsMiddleware struct{}

func NewCorsMiddleware() CorsMiddleware {
	return &corsMiddleware{}
}

func (m *corsMiddleware) SettingCors(c *minigin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Add("Access-Control-Allow-Headers", "Content-Type,Accept,Origin,x-token")

	// プリフライトリクエストは処理を通さない
	if c.Request.Method == http.MethodOptions {
		c.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	c.Next()
}
