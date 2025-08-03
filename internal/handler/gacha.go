package handler

import (
	"net/http"
	"road2ca/internal/service"
	"road2ca/pkg/minigin"

	"encoding/json"
	"fmt"
	"errors"
)

type GachaHandler interface {
	HandleGachaDraw(c *minigin.Context)
}

type gachaHandler struct {
	gachaService service.GachaService
}

func NewGachaHandler(gachaService service.GachaService) GachaHandler {
	return &gachaHandler{
		gachaService: gachaService,
	}
}

// HandleGachaDraw ガチャを引く処理
func (h *gachaHandler) HandleGachaDraw(c *minigin.Context) {
	var req service.DrawGachaRequestDTO

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, minigin.H{
			"error": "Invalid request body",
		})
		return
	}

	if req.Times < 1 || req.Times > 100 {
		c.Error(fmt.Errorf("times must be between 1 and 100"))
		c.JSON(http.StatusBadRequest, minigin.H{
			"error": "Times must be between 1 and 100",
		})
		return
	}

	res, err := h.gachaService.DrawGacha(c, req.Times)
	if err != nil {
		c.Error(err)
		if errors.Is(err, service.ErrNotEnoughCoins) {
			c.JSON(http.StatusBadRequest, minigin.H{
				"error": "Not enough coins",
			})
		} else {
			c.JSON(http.StatusInternalServerError, minigin.H{
				"error": "Failed to draw gacha",
			})
		}
		return
	}

	c.JSON(http.StatusOK, res)
}
