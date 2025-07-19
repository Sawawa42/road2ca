package handler

import (
	"net/http"
	"road2ca/internal/service"
	"road2ca/pkg/minigin"

	"encoding/json"
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
	var req service.GachaRequestBody

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.JSON(http.StatusBadRequest, minigin.H{
			"error": "Invalid request body",
		})
		return
	}

	results, err := h.gachaService.Draw(c, req.Times)
	if err != nil {
		c.JSON(http.StatusInternalServerError, minigin.H{
			"error": "Failed to draw gacha",
		})
		return
	}

	c.JSON(http.StatusOK, minigin.H{
		"results": results,
	})
}
