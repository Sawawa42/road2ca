package handler

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"road2ca/pkg/server/minigin"
	"road2ca/internal/model"
)

type UserCreateRequest struct {
	Name string `json:"name"`
}

// HandleUserCreate ユーザ登録処理
func (h *Handler) HandleUserCreate(c *minigin.Context) {
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

	// 簡易的にmd5ハッシュをトークンとして使用
	token := fmt.Sprintf("%x", md5.Sum([]byte(req.Name)))

	_, err := h.userDAO.Create(&model.User{
		Name:  req.Name,
		Token: token,
	})
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.Write([]byte(`{"error": "Internal server error"}`))
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write([]byte(`{"token": "` + token + `"}`))
	c.Next()
}
