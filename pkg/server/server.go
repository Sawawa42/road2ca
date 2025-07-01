package server

import (
	"log"
	"road2ca/pkg/http/middleware"
	"road2ca/pkg/server/handler"
	"road2ca/pkg/server/minigin"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// Serve HTTPサーバを起動する
func Serve(addr string) {
	// TODO: このあたりをまとめる
	db, err := sql.Open("mysql", "root:ca-tech-dojo@tcp(localhost:3306)/road2ca?parseTime=true")
	if err != nil {
		log.Fatalf("Failed to connect to database: %+v", err)
	}
	defer db.Close()

	// DB接続の確認
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %+v", err)
	}
	// このあたりここまで

	// ハンドラの初期化
	h := handler.New(db)

	// ミドルウェアの初期化
	m := middleware.NewMiddleware(db)

	// ルーターの初期化
	router := minigin.New()

	// CORS対応など共通の設定を適用
	router.Use(middleware.CommonConfig())

	/* ===== URLマッピングを行う ===== */
	router.GET("/setting/get", h.HandleSettingGet)

	router.POST("/user/create", h.HandleUserCreate)

	userGroup := router.Group("/user")
	{
		userGroup.Use(m.Authenticate) // 認証ミドルウェアを適用
		userGroup.GET("/get", h.HandleUserGet)
		// userGroup.POST("/update", h.HandleUserUpdate) // ex04
	}

	// TODO: 認証を行うmiddlewareを実装する
	// middlewareは pkg/http/middleware パッケージを利用する
	// http.HandleFunc("/user/get",
	//   get(middleware.Authenticate(handler.HandleUserGet())))

	/* ===== サーバの起動 ===== */
	log.Println("Server running...")
	err = router.Run(addr) 
	if err != nil {
		log.Fatalf("Server failed to start on %s: %+v", addr, err)
	}
}
