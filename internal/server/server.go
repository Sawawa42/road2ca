package server

import (
	"log"
	"road2ca/internal/handler"
	"road2ca/internal/middleware"
	"road2ca/pkg/minigin"
)

// Serve HTTPサーバを起動する
func Serve(addr string, h *handler.Handler, m *middleware.Middleware) {
	// ルーターの初期化
	router := minigin.New()

	// CORS対応など共通の設定を適用
	router.Use(m.Cors.SettingCors)

	// ロガーを適用
	router.Use(m.Logger.SettingLogger)

	router.GET("/setting/get", h.Setting.HandleGetSetting)

	userGroup := router.Group("/user")
	{
		userGroup.POST("/create", h.User.HandleCreateUser)

		// 認証ミドルウェアを適用
		userGroup.Use(m.Auth.Authenticate)
		userGroup.GET("/get", h.User.HandleGetUser)
		userGroup.POST("/update", h.User.HandleUpdateUser)
	}

	authGroup := router.Group("/")
	{
		// 認証ミドルウェアを適用
		authGroup.Use(m.Auth.Authenticate)
		authGroup.GET("/collection/list", h.Collection.HandleGetCollectionList)
		authGroup.GET("/ranking/list", h.Ranking.HandleGetRankingList)
		authGroup.POST("/game/finish", h.Game.HandleGameFinish)
		authGroup.POST("/gacha/draw", h.Gacha.HandleGachaDraw)
	}


	// サーバを起動
	log.Println("Server running...")
	if err := router.Run(addr); err != nil {
		log.Fatalf("Server failed to start on %s: %+v", addr, err)
	}
}
