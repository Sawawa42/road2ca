package handler

import (
	"log"
	"net/http"
	"road2ca/internal/service"
	"road2ca/pkg/minigin"
)

type SettingHandler interface {
	HandleSettingGet(c *minigin.Context)
}

type settingHandler struct {
	settingService service.SettingService
}

func NewSettingHandler() SettingHandler {
	return &settingHandler{
		settingService: service.NewSettingService(),
	}
}

// HandleSettingGet 設定情報を取得する
func (h *settingHandler) HandleSettingGet(c *minigin.Context) {
	res, err := h.settingService.GetSetting()
	if err != nil {
		log.Printf("ERROR: %v", err)
		c.JSON(http.StatusInternalServerError, minigin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, res)
}
