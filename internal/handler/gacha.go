package handler

import (
	"net/http"
	"road2ca/internal/service"
	"road2ca/pkg/minigin"
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
	c.JSON(http.StatusOK, minigin.H{
		"message": "Gacha draw endpoint is not implemented yet",
	})
}
