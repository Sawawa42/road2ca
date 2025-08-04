package handler

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"road2ca/internal/service"
	"road2ca/pkg/minigin"
)

type GameHandler interface {
	HandleGameFinish(c *minigin.Context)
}

type gameHandler struct {
	userService service.UserService
	gameService service.GameService
}

func NewGameHandler(userService service.UserService, gameService service.GameService) GameHandler {
	return &gameHandler{
		userService: userService,
		gameService: gameService,
	}
}

// HandleGameFinish ゲーム終了処理
func (h *gameHandler) HandleGameFinish(c *minigin.Context) {
	var req service.GameFinishRequestDTO
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, minigin.H{
			"error": "Invalid request body",
		})
		return
	}

	if req.Score == nil {
		c.Error(fmt.Errorf("score is required"))
		c.JSON(http.StatusBadRequest, minigin.H{
			"error": "Invalid request body",
		})
		return
	}

	if *req.Score < 0 || *req.Score > math.MaxInt32 {
		c.Error(fmt.Errorf("score must be between 0 and %d", math.MaxInt32))
		c.JSON(http.StatusBadRequest, minigin.H{
			"error": "Invalid request body",
		})
		return
	}

	// ゲーム終了処理を実行
	res, err := h.gameService.FinalizeGame(c, *req.Score)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, minigin.H{
			"error": "Internal server error",
		})
		return
	}

	// 現在のコイン数を送信
	c.JSON(http.StatusOK, res)
}
