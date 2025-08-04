package service

import (
	"context"
	"errors"
	"github.com/ilam072/wbtech-l0/backend/internal/types/domain"
	"github.com/ilam072/wbtech-l0/backend/internal/types/dto"
)

type OrderRepo interface {
	CreateOrder(context.Context, domain.Order, domain.Delivery, domain.Payment, []domain.Item) error
	GetOrder(ctx context.Context, ID string) (domain.FullOrder, error)
}

type OrderCache interface {
	Set(key string, order dto.Order)
	Get(key string) (dto.Order, bool)
}

type OrderConverter interface {
	DtoToDomainOrder(dto dto.Order) (domain.Order, domain.Delivery, domain.Payment, []domain.Item, error)
	DomainToDtoOrder(fullOrder domain.FullOrder) dto.Order
}

var (
	ErrOrderExists   = errors.New("order already exists")
	ErrOrderNotFound = errors.New("order not found")
	ErrInvalidUUID   = errors.New("invalid uuid")
)

type OrderService struct {
	orderRepo OrderRepo
	cache     OrderCache
	converter OrderConverter
}

func NewOrderService(repo OrderRepo, cache OrderCache, converter OrderConverter) *OrderService {
	return &OrderService{
		orderRepo: repo,
		cache:     cache,
		converter: converter,
	}
}
