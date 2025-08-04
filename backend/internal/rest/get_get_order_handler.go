package rest

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/ilam072/wbtech-l0/backend/internal/service"
)

// @Summary Get order by ID
// @Description Returns order details by given ID
// @Tags order
// @Param id path string true "order uid"
// @Success 200 {object} dto.Order
// @Failure 404 {object} ErrorResp "order not found"
// @Failure 500 {object} ErrorResp "internal server error"
// @Router /api/order/{id} [get]
func (h *Handler) GetOrderHandler(ctx *fiber.Ctx) error {
	orderId := ctx.Params("id")

	order, err := h.s.GetOrder(ctx.Context(), orderId)
	if err != nil {
		if errors.Is(err, service.ErrOrderNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(
				errorResponse("order not found"))
		}
		if errors.Is(err, service.ErrInvalidUUID) {
			return ctx.Status(fiber.StatusBadRequest).JSON(
				errorResponse("invalid uuid"))
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			errorResponse("something went wrong, try again later"))
	}

	return ctx.Status(fiber.StatusOK).JSON(order)
}
