package server

import (
	"log"
	"road2ca/pkg/http/middleware"
	"road2ca/pkg/minigin"
	"road2ca/pkg/server/handler"

	_ "github.com/go-sql-driver/mysql"
)

// Serve HTTPサーバを起動する
func Serve(addr string, h *handler.Handler, m *middleware.Middleware) {
	// ルーターの初期化
	router := minigin.New()

	// CORS対応など共通の設定を適用
	router.Use(m.CommonConfig)

	/* ===== URLマッピングを行う ===== */
	router.GET("/setting/get", h.HandleSettingGet)

	router.POST("/user/create", h.HandleUserCreate)

	userGroup := router.Group("/user")
	{
		// 認証ミドルウェアを適用
		userGroup.Use(m.Authenticate)

		userGroup.GET("/get", h.HandleUserGet)
		// userGroup.POST("/update", h.HandleUserUpdate) // ex04
	}

	/* ===== サーバの起動 ===== */
	log.Println("Server running...")
	if err := router.Run(addr); err != nil {
		log.Fatalf("Server failed to start on %s: %+v", addr, err)
	}
}
