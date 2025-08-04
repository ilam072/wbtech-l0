package rest

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/ilam072/wbtech-l0/backend/internal/types/dto"
	"log/slog"
)

type OrderService interface {
	GetOrder(ctx context.Context, orderId string) (dto.Order, error)
}

type Handler struct {
	log *slog.Logger
	api *fiber.App
	s   OrderService
}

func NewHandler(log *slog.Logger, s OrderService) *Handler {
	api := fiber.New()

	api.Get("/swagger/*", swagger.HandlerDefault)

	h := &Handler{
		log: log,
		api: api,
		s:   s,
	}
	h.api.Get("/api/order/:id", h.GetOrderHandler)

	return h
}

func (h *Handler) Listen(addr string) error {
	return h.api.Listen(addr)
}

func (h *Handler) Shutdown() error {
	return h.api.Shutdown()
}
