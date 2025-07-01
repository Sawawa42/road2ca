package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"road2ca/pkg/server/minigin"
	"database/sql"
)

type UserCreateRequest struct {
	Name string `json:"name"`
}

// HandleUserCreate ユーザ登録処理
func HandleUserCreate(db *sql.DB) minigin.HandlerFunc {
	return func(c *minigin.Context) {
		var req UserCreateRequest

		if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
			log.Printf("Failed to decode request body: %v", err)
			c.Writer.WriteHeader(http.StatusBadRequest)
			c.Writer.Write([]byte(`{"error": "Invalid request"}`))
			return
		}

		if req.Name == "" {
			c.Writer.WriteHeader(http.StatusBadRequest)
			c.Writer.Write([]byte(`{"error": "Invalid request"}`))
			return
		}

		// TODO: tokenの生成処理を実装

		// TODO: DAOの実装
		_, err := db.Exec("INSERT INTO users (name, token) VALUES (?, ?)", req.Name, "example_token")
		if err != nil {
			log.Printf("Failed to create user: %v", err)
			c.Writer.WriteHeader(http.StatusInternalServerError)
			c.Writer.Write([]byte(`{"error": "Internal server error"}`))
			return
		}

		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write([]byte(`{"token": "example_token"}`))
		c.Next()
	}
}
