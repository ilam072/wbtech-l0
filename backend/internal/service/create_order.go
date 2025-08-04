package service

import (
	"context"
	"errors"
	"github.com/ilam072/wbtech-l0/backend/internal/repo"
	"github.com/ilam072/wbtech-l0/backend/internal/types/dto"
	"github.com/ilam072/wbtech-l0/backend/pkg/e"
)

func (s OrderService) CreateOrder(ctx context.Context, order dto.Order) error {
	const op = "OrderService.CreateOrder()"

	domainOrder, delivery, payment, items, err := s.converter.DtoToDomainOrder(order)
	if err != nil {
		return e.Wrap(op, err)
	}

	if err := s.orderRepo.CreateOrder(ctx, domainOrder, delivery, payment, items); err != nil {
		if errors.Is(err, repo.ErrOrderExists) {
			return e.Wrap(op, ErrOrderExists)
		}
		return e.Wrap(op, err)
	}

	s.cache.Set(order.OrderUID, order)
	return nil
}
