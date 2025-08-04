package cache

import (
	"context"
	"github.com/ilam072/wbtech-l0/backend/internal/types/domain"
	"github.com/ilam072/wbtech-l0/backend/internal/types/dto"
	"github.com/ilam072/wbtech-l0/backend/pkg/e"
	"github.com/maypok86/otter/v2"
)

type OrderRepo interface {
	GetLastOrders(ctx context.Context, limit int) ([]domain.FullOrder, error)
}

type OrderConverter interface {
	DtoToDomainOrder(dto dto.Order) (domain.Order, domain.Delivery, domain.Payment, []domain.Item, error)
	DomainToDtoOrder(fullOrder domain.FullOrder) dto.Order
}

type OrderCache struct {
	store     *otter.Cache[string, dto.Order]
	orderRepo OrderRepo
	converter OrderConverter
}

func New(orderRepo OrderRepo, converter OrderConverter) *OrderCache {
	opts := &otter.Options[string, dto.Order]{
		MaximumSize: 1000,
	}

	cache := otter.Must[string, dto.Order](opts)

	return &OrderCache{
		store:     cache,
		orderRepo: orderRepo,
		converter: converter,
	}
}

func (c *OrderCache) Set(key string, order dto.Order) {
	c.store.Set(key, order)
}

func (c *OrderCache) Get(key string) (dto.Order, bool) {
	order, ok := c.store.GetIfPresent(key)

	return order, ok
}

func (c *OrderCache) Preload(ctx context.Context, limit int) error {
	const op = "cache.Preload()"

	orders, err := c.orderRepo.GetLastOrders(ctx, limit)
	if err != nil {
		return e.Wrap(op, err)
	}

	if orders == nil {
		return nil
	}

	for _, order := range orders {
		o := order
		c.Set(o.Order.ID.String(), c.converter.DomainToDtoOrder(o))
	}
	return nil
}
