package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"road2ca/pkg/server/minigin"
	"database/sql"
)

const (
	// ガチャ1回あたりのコイン消費量
	GachaCoinConsumption = 100
)

// HandleSettingGet ゲーム設定情報取得処理
func HandleSettingGet(db *sql.DB) minigin.HandlerFunc {
	return func(c *minigin.Context) {
		data, err := json.Marshal(&settingGetResponse{
			GachaCoinConsumption: GachaCoinConsumption,
		})
		if err != nil {
			log.Println(err)
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Write(data)
		c.Next()
	}
}

type settingGetResponse struct {
	GachaCoinConsumption int32 `json:"gachaCoinConsumption"`
}

