package handler

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"road2ca/pkg/server/minigin"
	"road2ca/internal/model"
	"road2ca/pkg/contextKey"
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

func (h *Handler) HandleUserGet(c *minigin.Context) {
	token, ok := c.Request.Context().Value(contextkey.AuthToken).(string)
	if !ok {
		// 仮
		log.Println("No token found in context")
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.Writer.Write([]byte(`{"error": "Unauthorized"}`))
		return
	}
	log.Printf("HandleUserGet called with token: %s", token)
	
	// 仮で固定データを返す
	user := &model.User{
		ID:    1,
		Name:  "John Doe",
		HighScore: 1000,
		Coin: 500,
	}
	
	response, err := json.Marshal(user)
	if err != nil {
		log.Printf("Failed to marshal user data: %v", err)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Write([]byte(`{"error": "Internal server error"}`))
		return
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Write(response)
	c.Next()
}
