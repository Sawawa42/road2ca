package handler

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"road2ca/internal/model"
	"road2ca/pkg/constants"
	"road2ca/pkg/minigin"
)

type UserCreateRequest struct {
	Name string `json:"name"`
}

// HandleUserCreate ユーザ登録処理
func (h *Handler) HandleUserCreate(c *minigin.Context) {
	var req UserCreateRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		http.Error(c.Writer, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(c.Writer, "Invalid request", http.StatusBadRequest)
		return
	}

	// 簡易的にmd5ハッシュをトークンとして使用
	token := fmt.Sprintf("%x", md5.Sum([]byte(req.Name)))

	_, err := h.userDAO.Create(&model.User{
		Name:  req.Name,
		Token: token,
	})
	if err != nil {
		log.Printf("Failed to HandleUserCreate: %v", err)
		http.Error(c.Writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, &model.UserCreateResponse{
		Token: token,
	})
}

// HandleUserGet ユーザ情報取得処理
func (h *Handler) HandleUserGet(c *minigin.Context) {
	user, ok := c.Request.Context().Value(constants.ContextKey).(*model.User)
	if !ok {
		log.Println("Failed to HandleUserGet: user not found in context")
		http.Error(c.Writer, "Unauthorized", http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, &model.UserGetResponse{
		ID:        user.ID,
		Name:      user.Name,
		HighScore: user.HighScore,
		Coin:      user.Coin,
	})
}
