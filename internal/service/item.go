package service

type ItemService interface {
}

type itemService struct {
}

func NewItemService() ItemService {
	return &itemService{}
}
