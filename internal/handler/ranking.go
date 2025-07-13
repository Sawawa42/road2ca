package handler

import (
	"log"
	"net/http"
	"road2ca/internal/service"
	"road2ca/pkg/minigin"
)

type RankingHandler interface {
	HandleGetRankingList(c *minigin.Context)
}

type rankingHandler struct {
	rankingService service.RankingService
}

func NewRankingHandler(rankingService service.RankingService) RankingHandler {
	return &rankingHandler{
		rankingService: rankingService,
	}
}

// HandleGetRankingList ランキング一覧取得処理
func (h *rankingHandler) HandleGetRankingList(c *minigin.Context) {
	// クエリパラメータからstartを取得
	start, err := c.QueryInt("start")
	if err != nil {
		log.Printf("ERROR: %v", err)
		c.JSON(http.StatusBadRequest, minigin.H{
			"error": "Invalid start parameter",
		})
		return
	}

	if start <= 0 {
		c.JSON(http.StatusBadRequest, minigin.H{
			"error": "Invalid input",
		})
		return
	}

	// TODO: endは設定から取得する予定
	res, err := h.rankingService.GetInRange(start, 10)
	if err != nil {
		log.Printf("ERROR: %v", err)
		c.JSON(http.StatusInternalServerError, minigin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, minigin.H{
		"ranks": res,
	})
}