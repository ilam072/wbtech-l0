package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ilam072/wbtech-l0/backend/internal/service"
	"github.com/ilam072/wbtech-l0/backend/internal/types/dto"
	"github.com/ilam072/wbtech-l0/backend/pkg/logger/sl"
	"github.com/segmentio/kafka-go"
	"log/slog"
)

//go:generate mockgen -source=order_handler.go -destination=../../../../mocks/kafka/mock_order_handler.go -package kafka
type Consumer interface {
	Consume(context.Context) (kafka.Message, error)
	Close() error
}

type Service interface {
	CreateOrder(context.Context, dto.Order) error
}

type Validator interface {
	Validate(i interface{}) error
}

type OrderConsumerHandler struct {
	log       *slog.Logger
	consumer  Consumer
	service   Service
	validator Validator
}

func NewOrderConsumerHandler(
	log *slog.Logger,
	c Consumer,
	s Service,
	v Validator,
) *OrderConsumerHandler {
	return &OrderConsumerHandler{
		log:       log,
		consumer:  c,
		service:   s,
		validator: v,
	}
}

func (h *OrderConsumerHandler) Start(ctx context.Context) error {
	const op = "kafka.handler.Start()"

	log := h.log.With(
		slog.String("op", op),
	)

	for {
		select {
		case <-ctx.Done():
			log.Info("kafka consumer shutting down...")
			return nil
		default:
			message, err := h.consumer.Consume(ctx)
			if err != nil {
				log.Warn("failed to read message", slog.String("error", err.Error()))
				continue
			}

			order := dto.Order{}

			if err := json.Unmarshal(message.Value, &order); err != nil {
				log.Error("failed to decode json message to order", sl.Err(err))
				continue
			}

			if err := h.validator.Validate(order); err != nil {
				log.Warn("failed to validate order", slog.String("error", err.Error()))
				continue
			}

			if err = h.service.CreateOrder(ctx, order); err != nil {
				if errors.Is(err, service.ErrOrderExists) {
					log.Warn("order with such uid already exists", slog.String("error", err.Error()))
					continue
				}
				log.Error("failed to create order", sl.Err(err))
			}
		}
	}
}
