package rest

import (
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/ilam072/wbtech-l0/backend/internal/service"
	"github.com/ilam072/wbtech-l0/backend/internal/types/dto"
	httpmock "github.com/ilam072/wbtech-l0/backend/mocks/http"
	"github.com/ilam072/wbtech-l0/backend/pkg/logger/handlers/slogdiscard"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetOrderHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := slogdiscard.NewDiscardLogger()
	mockService := httpmock.NewMockOrderService(ctrl)
	h := NewHandler(log, mockService)

	app := fiber.New()
	app.Get("/api/order/:id", h.GetOrderHandler)

	orderId := uuid.New().String()
	expectedOrder := dto.Order{OrderUID: orderId}

	mockService.EXPECT().GetOrder(gomock.Any(), orderId).Return(expectedOrder, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/order/"+orderId, nil)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var order dto.Order
	err := json.NewDecoder(resp.Body).Decode(&order)
	assert.NoError(t, err)
	assert.Equal(t, expectedOrder.OrderUID, order.OrderUID)
}

func TestHandler_GetOrderHandler_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := slogdiscard.NewDiscardLogger()
	mockService := httpmock.NewMockOrderService(ctrl)
	h := NewHandler(log, mockService)

	app := fiber.New()
	app.Get("/api/order/:id", h.GetOrderHandler)

	orderId := uuid.New().String()

	mockService.EXPECT().GetOrder(gomock.Any(), orderId).Return(dto.Order{}, service.ErrOrderNotFound)

	req := httptest.NewRequest(http.MethodGet, "/api/order/"+orderId, nil)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestHandler_GetOrderHandler_InvalidUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := slogdiscard.NewDiscardLogger()
	mockService := httpmock.NewMockOrderService(ctrl)
	h := NewHandler(log, mockService)

	app := fiber.New()
	app.Get("/api/order/:id", h.GetOrderHandler)

	invalidOrderID := "invalid-uuid"

	mockService.EXPECT().GetOrder(gomock.Any(), invalidOrderID).Return(dto.Order{}, service.ErrInvalidUUID)

	req := httptest.NewRequest(http.MethodGet, "/api/order/"+invalidOrderID, nil)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestHandler_GetOrderHandler_InternalError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := slogdiscard.NewDiscardLogger()
	mockService := httpmock.NewMockOrderService(ctrl)
	h := NewHandler(log, mockService)

	app := fiber.New()
	app.Get("/api/order/:id", h.GetOrderHandler)

	orderId := uuid.New().String()

	mockService.EXPECT().GetOrder(gomock.Any(), orderId).Return(dto.Order{}, errors.New("internal error"))

	req := httptest.NewRequest(http.MethodGet, "/api/order/"+orderId, nil)
	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
