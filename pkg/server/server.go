package server

import (
	"log"
	// "net/http"

	"road2ca/pkg/http/middleware"
	"road2ca/pkg/server/handler"
	"road2ca/pkg/server/minigin"
)

// Serve HTTPサーバを起動する
func Serve(addr string) {
	router := minigin.New()

	// CORS対応など共通の設定を適用
	router.Use(middleware.CommonConfig())

	/* ===== URLマッピングを行う ===== */
	router.GET("/setting/get", handler.HandleSettingGet())

	router.POST("/user/create", handler.HandleUserCreate())

	// TODO: 認証を行うmiddlewareを実装する
	// middlewareは pkg/http/middleware パッケージを利用する
	// http.HandleFunc("/user/get",
	//   get(middleware.Authenticate(handler.HandleUserGet())))

	/* ===== サーバの起動 ===== */
	log.Println("Server running...")
	err := router.Run(addr)
	if err != nil {
		log.Fatalf("Server failed to start on %s: %+v", addr, err)
	}
}
