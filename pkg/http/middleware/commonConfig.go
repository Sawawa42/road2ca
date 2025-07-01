package middleware

import (
	"net/http"
	"road2ca/pkg/server/minigin"
)

func CommonConfig() minigin.HandlerFunc {
	return func(c *minigin.Context) {
		// CORS対応
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Add("Access-Control-Allow-Headers", "Content-Type,Accept,Origin,x-token")

		// プリフライトリクエストは処理を通さない
		if c.Request.Method == http.MethodOptions {
			c.Writer.WriteHeader(http.StatusNoContent)
			return
		}

		// 共通のレスポンスヘッダを設定
		c.Writer.Header().Add("Content-Type", "application/json")
		c.Next()
	}
}

