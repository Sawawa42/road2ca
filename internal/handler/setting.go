package handler

import (
	"net/http"
	"road2ca/pkg/minigin"
	"road2ca/internal/service"
	"log"
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
		http.Error(c.Writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, res)
}
