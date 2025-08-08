package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/ilam072/wbtech-l0/backend/internal/service"
	"github.com/ilam072/wbtech-l0/backend/internal/types/dto"
	kafkamocks "github.com/ilam072/wbtech-l0/backend/mocks/kafka"
	"github.com/ilam072/wbtech-l0/backend/pkg/logger/handlers/slogdiscard"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestOrderConsumerHandler_Start_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := slogdiscard.NewDiscardLogger()
	consumer := kafkamocks.NewMockConsumer(ctrl)
	mockService := kafkamocks.NewMockService(ctrl)
	validator := kafkamocks.NewMockValidator(ctrl)

	h := NewOrderConsumerHandler(log, consumer, mockService, validator)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	order := dto.Order{
		OrderUID: uuid.New().String(),
	}
	bytes, err := json.Marshal(order)
	require.NoError(t, err)

	gomock.InOrder(
		consumer.EXPECT().Consume(gomock.Any()).Return(kafka.Message{Value: bytes}, nil),
		validator.EXPECT().Validate(order).Return(nil),
		mockService.EXPECT().CreateOrder(gomock.Any(), order).Return(nil),
		consumer.EXPECT().Consume(gomock.Any()).DoAndReturn(func(ctx context.Context) (kafka.Message, error) {
			cancel()
			return kafka.Message{}, context.Canceled
		}),
	)

	err = h.Start(ctx)
	assert.NoError(t, err)
}

func TestOrderConsumerHandler_Start_ContextCancelled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	consumer := kafkamocks.NewMockConsumer(ctrl)
	validator := kafkamocks.NewMockValidator(ctrl)
	mockService := kafkamocks.NewMockService(ctrl)
	logger := slogdiscard.NewDiscardLogger()

	handler := &OrderConsumerHandler{
		consumer:  consumer,
		validator: validator,
		service:   mockService,
		log:       logger,
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := handler.Start(ctx)
	assert.NoError(t, err)

	consumer.EXPECT().Consume(gomock.Any()).Times(0)
}

func TestOrderConsumerHandler_Start_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	consumer := kafkamocks.NewMockConsumer(ctrl)
	validator := kafkamocks.NewMockValidator(ctrl)
	mockService := kafkamocks.NewMockService(ctrl)
	logger := slogdiscard.NewDiscardLogger()

	h := &OrderConsumerHandler{
		consumer:  consumer,
		validator: validator,
		service:   mockService,
		log:       logger,
	}

	ctx, cancel := context.WithCancel(context.Background())

	invalidMessage := []byte("invalid json")
	consumer.EXPECT().Consume(gomock.Any()).Return(kafka.Message{Value: invalidMessage}, nil)
	consumer.EXPECT().Consume(gomock.Any()).DoAndReturn(func(ctx context.Context) (kafka.Message, error) {
		cancel()
		return kafka.Message{}, context.Canceled
	})

	err := h.Start(ctx)
	assert.NoError(t, err)

}

func TestOrderConsumerHandler_Start_ValidationFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	consumer := kafkamocks.NewMockConsumer(ctrl)
	validator := kafkamocks.NewMockValidator(ctrl)
	mockService := kafkamocks.NewMockService(ctrl)
	logger := slogdiscard.NewDiscardLogger()

	h := &OrderConsumerHandler{
		consumer:  consumer,
		validator: validator,
		service:   mockService,
		log:       logger,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	order := dto.Order{
		OrderUID: "invalid uuid",
	}
	message, err := json.Marshal(order)
	require.NoError(t, err)

	gomock.InOrder(
		consumer.EXPECT().Consume(gomock.Any()).Return(kafka.Message{Value: message}, nil),
		validator.EXPECT().Validate(order).Return(errors.New("validation error")),
		consumer.EXPECT().Consume(gomock.Any()).DoAndReturn(func(ctx context.Context) (kafka.Message, error) {
			cancel()
			return kafka.Message{}, context.Canceled
		}),
	)

	err = h.Start(ctx)
	assert.NoError(t, err)
}

func TestOrderConsumerHandler_Start_CreateOrderFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	consumer := kafkamocks.NewMockConsumer(ctrl)
	validator := kafkamocks.NewMockValidator(ctrl)
	mockService := kafkamocks.NewMockService(ctrl)
	logger := slogdiscard.NewDiscardLogger()

	h := &OrderConsumerHandler{
		consumer:  consumer,
		validator: validator,
		service:   mockService,
		log:       logger,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	order := dto.Order{
		OrderUID: uuid.New().String(),
	}
	bytes, err := json.Marshal(order)
	require.NoError(t, err)

	gomock.InOrder(
		consumer.EXPECT().Consume(gomock.Any()).Return(kafka.Message{Value: bytes}, nil),
		validator.EXPECT().Validate(order).Return(nil),
		mockService.EXPECT().CreateOrder(gomock.Any(), order).Return(service.ErrOrderExists),
		consumer.EXPECT().Consume(gomock.Any()).DoAndReturn(func(ctx context.Context) (kafka.Message, error) {
			cancel()
			return kafka.Message{}, context.Canceled
		}),
	)

	err = h.Start(ctx)
}
