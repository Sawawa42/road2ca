package handler

import (
	"encoding/json"
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
	var req service.CreateUserRequestDTO

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
		c.JSON(http.StatusInternalServerError, minigin.H{
			"error": "Internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

// HandleGetUser ユーザ情報取得処理
func (h *userHandler) HandleGetUser(c *minigin.Context) {
	res, err := h.userService.GetUser(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, minigin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

// HandleUpdateUser ユーザ情報更新処理
func (h *userHandler) HandleUpdateUser(c *minigin.Context) {
	var req service.UpdateUserRequestDTO

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

	if err := h.userService.UpdateUser(c, req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, minigin.H{
			"error": "Internal server error",
		})
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
