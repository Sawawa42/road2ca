package handler

import (
	"net/http"
	"road2ca/internal/service"
	"road2ca/pkg/minigin"
)

type CollectionHandler interface {
	HandleGetCollectionList(c *minigin.Context)
}

type collectionHandler struct {
	collectionService service.CollectionService
}

func NewCollectionHandler(collectionService service.CollectionService) CollectionHandler {
	return &collectionHandler{
		collectionService: collectionService,
	}
}

// HandleGetCollectionList コレクション一覧取得処理
func (h *collectionHandler) HandleGetCollectionList(c *minigin.Context) {
	res, err := h.collectionService.GetCollectionList(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, minigin.H{
			"error": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, &service.CollectionListResponseDTO{
		Collections: res,
	})
}
