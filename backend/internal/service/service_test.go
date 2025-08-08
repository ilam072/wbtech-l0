package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/ilam072/wbtech-l0/backend/internal/repo"
	"github.com/ilam072/wbtech-l0/backend/internal/types/domain"
	"github.com/ilam072/wbtech-l0/backend/internal/types/dto"
	mocks "github.com/ilam072/wbtech-l0/backend/mocks/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestOrderService_CreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepo(ctrl)
	cache := mocks.NewMockOrderCache(ctrl)
	converter := mocks.NewMockOrderConverter(ctrl)
	service := NewOrderService(mockRepo, cache, converter)

	dtoOrder := dto.Order{
		OrderUID:    uuid.New().String(),
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: dto.Delivery{
			Name:    "John Doe",
			Phone:   "+1234567890",
			Zip:     "123456",
			City:    "Test City",
			Address: "123 Test St",
			Region:  "Test Region",
			Email:   "test@gmail.com",
		},
		Payment: dto.Payment{
			Transaction:  "txn123",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "TestProvider",
			Amount:       1000,
			PaymentDt:    1637907727,
			Bank:         "Test Bank",
			DeliveryCost: 100,
			GoodsTotal:   900,
			CustomFee:    0,
		},
		Items: []dto.Item{
			{
				ChrtID:      123456,
				TrackNumber: "WBILMTESTTRACK",
				Price:       1000,
				Rid:         "rid123",
				Name:        "Test Item",
				Sale:        30,
				Size:        "M",
				TotalPrice:  700,
				NmID:        123456789,
				Brand:       "Test Brand",
				Status:      1,
			},
		},
		Locale:            "en",
		InternalSignature: "signature",
		CustomerID:        "customer123",
		DeliveryService:   "Test Delivery Service",
		Shardkey:          "shardkey123",
		SmID:              1,
		DateCreated:       time.Now(),
		OofShard:          "oofshard123",
	}

	domainOrder := domain.Order{
		ID:                uuid.MustParse(dtoOrder.OrderUID),
		TrackNumber:       dtoOrder.TrackNumber,
		Entry:             dtoOrder.Entry,
		Locale:            dtoOrder.Locale,
		InternalSignature: dtoOrder.InternalSignature,
		CustomerID:        dtoOrder.CustomerID,
		DeliveryService:   dtoOrder.DeliveryService,
		ShardKey:          dtoOrder.Shardkey,
		SmID:              dtoOrder.SmID,
		DateCreated:       dtoOrder.DateCreated,
		OofShard:          dtoOrder.OofShard,
	}

	delivery := domain.Delivery{
		OrderID: domainOrder.ID,
		Name:    dtoOrder.Delivery.Name,
		Phone:   dtoOrder.Delivery.Phone,
		Zip:     dtoOrder.Delivery.Zip,
		City:    dtoOrder.Delivery.City,
		Address: dtoOrder.Delivery.Address,
		Region:  dtoOrder.Delivery.Region,
		Email:   dtoOrder.Delivery.Email,
	}

	payment := domain.Payment{
		Transaction:  domainOrder.ID,
		OrderID:      domainOrder.ID,
		RequestID:    dtoOrder.Payment.RequestID,
		Currency:     dtoOrder.Payment.Currency,
		Provider:     dtoOrder.Payment.Provider,
		Amount:       dtoOrder.Payment.Amount,
		PaymentDt:    time.Unix(dtoOrder.Payment.PaymentDt, 0),
		Bank:         dtoOrder.Payment.Bank,
		DeliveryCost: dtoOrder.Payment.DeliveryCost,
		GoodsTotal:   dtoOrder.Payment.GoodsTotal,
		CustomFee:    dtoOrder.Payment.CustomFee,
	}

	var items []domain.Item
	for _, itm := range dtoOrder.Items {
		item := domain.Item{
			ChrtID:      int64(itm.ChrtID),
			OrderID:     domainOrder.ID,
			TrackNumber: itm.TrackNumber,
			Price:       itm.Price,
			Rid:         itm.Rid,
			Name:        itm.Name,
			Sale:        itm.Sale,
			Size:        itm.Size,
			TotalPrice:  itm.TotalPrice,
			NmID:        int64(itm.NmID),
			Brand:       itm.Brand,
			Status:      itm.Status,
		}
		items = append(items, item)
	}

	converter.EXPECT().DtoToDomainOrder(dtoOrder).Return(domainOrder, delivery, payment, items, nil)

	ctx := context.Background()
	mockRepo.EXPECT().CreateOrder(ctx, domainOrder, delivery, payment, items).Return(nil)
	cache.EXPECT().Set(dtoOrder.OrderUID, dtoOrder)

	err := service.CreateOrder(ctx, dtoOrder)
	assert.NoError(t, err)
}

func TestOrderService_GetOrderFromCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepo(ctrl)
	cache := mocks.NewMockOrderCache(ctrl)
	converter := mocks.NewMockOrderConverter(ctrl)
	service := NewOrderService(mockRepo, cache, converter)

	orderId := uuid.New().String()
	expectedOrder := dto.Order{
		OrderUID: orderId,
	}
	cache.EXPECT().Get(orderId).Return(expectedOrder, true)

	order, err := service.GetOrder(context.Background(), orderId)
	assert.NoError(t, err)
	assert.Equal(t, expectedOrder, order)
}

func TestOrderService_GetOrderFromRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepo(ctrl)
	cache := mocks.NewMockOrderCache(ctrl)
	converter := mocks.NewMockOrderConverter(ctrl)
	service := NewOrderService(mockRepo, cache, converter)

	orderId := uuid.New().String()

	cache.EXPECT().Get(orderId).Return(dto.Order{}, false)

	fullOrder := domain.FullOrder{
		Order: domain.Order{
			ID: uuid.MustParse(orderId),
		},
	}
	mockRepo.EXPECT().GetOrder(gomock.Any(), orderId).Return(fullOrder, nil)

	dtoOrder := dto.Order{
		OrderUID: fullOrder.Order.ID.String(),
	}
	converter.EXPECT().DomainToDtoOrder(fullOrder).Return(dtoOrder)

	order, err := service.GetOrder(context.Background(), orderId)
	assert.NoError(t, err)
	assert.Equal(t, dtoOrder, order)
}

func TestOrderService_GetOrder_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepo(ctrl)
	cache := mocks.NewMockOrderCache(ctrl)
	converter := mocks.NewMockOrderConverter(ctrl)
	service := NewOrderService(mockRepo, cache, converter)

	orderId := uuid.New().String()

	cache.EXPECT().Get(orderId).Return(dto.Order{}, false)
	mockRepo.EXPECT().GetOrder(context.Background(), orderId).Return(domain.FullOrder{}, repo.ErrOrderNotFound)

	dtoOrder, err := service.GetOrder(context.Background(), orderId)
	assert.ErrorIs(t, err, ErrOrderNotFound)
	assert.Empty(t, dtoOrder)
}

func TestOrderService_GetOrder_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderRepo(ctrl)
	cache := mocks.NewMockOrderCache(ctrl)
	converter := mocks.NewMockOrderConverter(ctrl)
	service := NewOrderService(mockRepo, cache, converter)

	orderId := uuid.New().String()
	repoErr := errors.New("mockRepo error")

	cache.EXPECT().Get(orderId).Return(dto.Order{}, false)
	mockRepo.EXPECT().GetOrder(context.Background(), orderId).Return(domain.FullOrder{}, repoErr)

	dtoOrder, err := service.GetOrder(context.Background(), orderId)
	assert.ErrorContains(t, err, repoErr.Error())
	assert.ErrorContains(t, err, "OrderService.GetOrder()")

	assert.Empty(t, dtoOrder)
}
