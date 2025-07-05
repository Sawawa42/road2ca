package handler

import (
	"net/http"
	"road2ca/internal/model"
	"road2ca/pkg/minigin"
)

const (
	// ガチャ1回あたりのコイン消費量
	GachaCoinConsumption = 100
)

// HandleSettingGet ゲーム設定情報取得処理
func (h *Handler) HandleSettingGet(c *minigin.Context) {
	c.JSON(http.StatusOK, &model.SettingGetResponse{
		GachaCoinConsumption: GachaCoinConsumption,
	})
}
