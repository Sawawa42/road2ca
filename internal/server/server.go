package server

import (
	"log"
	"road2ca/internal/handler"
	"road2ca/internal/middleware"
	"road2ca/pkg/minigin"

	_ "github.com/go-sql-driver/mysql"
)

// Serve HTTPサーバを起動する
func Serve(addr string, h *handler.Handler, m *middleware.Middleware) {
	// ルーターの初期化
	router := minigin.New()

	// CORS対応など共通の設定を適用
	router.Use(m.Cors.SettingCors)

	/* ===== URLマッピングを行う ===== */
	router.GET("/setting/get", h.Setting.HandleSettingGet)

	router.POST("/user/create", h.User.HandleUserCreate)

	userGroup := router.Group("/user")
	{
		// 認証ミドルウェアを適用
		userGroup.Use(m.Auth.Authenticate)

		userGroup.GET("/get", h.User.HandleUserGet)
		// userGroup.POST("/update", h.HandleUserUpdate) // ex04
	}

	/* ===== サーバの起動 ===== */
	log.Println("Server running...")
	if err := router.Run(addr); err != nil {
		log.Fatalf("Server failed to start on %s: %+v", addr, err)
	}
}
