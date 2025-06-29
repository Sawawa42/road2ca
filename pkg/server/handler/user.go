package handler

import (
	"net/http"
	"road2ca/pkg/server/minigin"
)

// HandleUserCreate ユーザ登録処理
func HandleUserCreate() minigin.HandlerFunc {
	return func(c *minigin.Context) {
		// ユーザ登録のロジックをここに実装
		// 例えば、リクエストボディからユーザ情報を取得し、データベースに保存するなど

		// レスポンスを返す
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write([]byte(`{"token": "example_token"}`))
		c.Next()
	}
}