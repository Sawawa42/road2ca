package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"road2ca/internal/service"
	"road2ca/pkg/minigin"
)

type UserHandler interface {
	HandleUserCreate(c *minigin.Context)
	HandleUserGet(c *minigin.Context)
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

// HandleUserCreate ユーザ登録処理
func (h *userHandler) HandleUserCreate(c *minigin.Context) {
	var req service.UserCreateRequestDTO

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.JSON(http.StatusBadRequest, minigin.H{
			"error": "Invalid request body",
		})
		return
	}

	if req.Name == "" {
		c.JSON(http.StatusBadRequest, minigin.H{
			"error": "Name is required",
		})
		return
	}

	res, err := h.userService.CreateUser(req.Name)
	if err != nil {
		log.Printf("ERROR: %v", err)
		c.JSON(http.StatusInternalServerError, minigin.H{
			"error": "Internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

// HandleUserGet ユーザ情報取得処理
func (h *userHandler) HandleUserGet(c *minigin.Context) {
	res, err := h.userService.GetUser(c)
	if err != nil {
		log.Printf("ERROR: %v", err)
		c.JSON(http.StatusInternalServerError, minigin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, res)
}
