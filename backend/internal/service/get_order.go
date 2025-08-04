package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/ilam072/wbtech-l0/backend/internal/repo"
	"github.com/ilam072/wbtech-l0/backend/internal/types/dto"
	"github.com/ilam072/wbtech-l0/backend/pkg/e"
)

func (s OrderService) GetOrder(ctx context.Context, orderId string) (dto.Order, error) {
	const op = "OrderService.GetOrder()"

	if _, err := uuid.Parse(orderId); err != nil {
		return dto.Order{}, ErrInvalidUUID
	}

	order, ok := s.cache.Get(orderId)
	if ok {
		return order, nil
	}

	fullOrder, err := s.orderRepo.GetOrder(ctx, orderId)
	if err != nil {
		if errors.Is(err, repo.ErrOrderNotFound) {
			return dto.Order{}, e.Wrap(op, ErrOrderNotFound)
		}
		return dto.Order{}, e.Wrap(op, err)
	}

	return s.converter.DomainToDtoOrder(fullOrder), nil
}
