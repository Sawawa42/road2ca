package handler

import (
	"log"
	"net/http"
	"road2ca/internal/service"
	"road2ca/pkg/minigin"
)

type SettingHandler interface {
	HandleGetSetting(c *minigin.Context)
}

type settingHandler struct {
	settingService service.SettingService
}

func NewSettingHandler(settingService service.SettingService) SettingHandler {
	return &settingHandler{
		settingService: settingService,
	}
}

// HandleGetSetting 設定情報を取得する
func (h *settingHandler) HandleGetSetting(c *minigin.Context) {
	res, err := h.settingService.GetSettings()
	if err != nil {
		log.Printf("ERROR: %v", err)
		c.JSON(http.StatusInternalServerError, minigin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, res)
}
