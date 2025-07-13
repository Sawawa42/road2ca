package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"road2ca/internal/service"
	"road2ca/pkg/minigin"
)

type UserHandler interface {
	HandleCreateUser(c *minigin.Context)
	HandleGetUser(c *minigin.Context)
	HandleUpdateUser(c *minigin.Context)
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

// HandleCreateUser ユーザ登録処理
func (h *userHandler) HandleCreateUser(c *minigin.Context) {
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

	res, err := h.userService.Create(req.Name)
	if err != nil {
		log.Printf("ERROR: %v", err)
		c.JSON(http.StatusInternalServerError, minigin.H{
			"error": "Internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

// HandleGetUser ユーザ情報取得処理
func (h *userHandler) HandleGetUser(c *minigin.Context) {
	res, err := h.userService.Get(c)
	if err != nil {
		log.Printf("ERROR: %v", err)
		c.JSON(http.StatusInternalServerError, minigin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

// HandleUpdateUser ユーザ情報更新処理
func (h *userHandler) HandleUpdateUser(c *minigin.Context) {
	var req service.UserUpdateRequestDTO

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

	if err := h.userService.Update(c, req.Name); err != nil {
		log.Printf("ERROR: %v", err)
		c.JSON(http.StatusInternalServerError, minigin.H{
			"error": "Internal server error",
		})
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
