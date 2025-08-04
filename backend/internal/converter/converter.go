package converter

import (
	"github.com/google/uuid"
	"github.com/ilam072/wbtech-l0/backend/internal/types/domain"
	"github.com/ilam072/wbtech-l0/backend/internal/types/dto"
	"time"
)

type Converter struct{}

func New() *Converter {
	return &Converter{}
}

func (c *Converter) DtoToDomainOrder(dto dto.Order) (domain.Order, domain.Delivery, domain.Payment, []domain.Item, error) {
	orderUID, err := uuid.Parse(dto.OrderUID)
	if err != nil {
		return domain.Order{}, domain.Delivery{}, domain.Payment{}, nil, err
	}

	order := domain.Order{
		ID:                orderUID,
		TrackNumber:       dto.TrackNumber,
		Entry:             dto.Entry,
		Locale:            dto.Locale,
		InternalSignature: dto.InternalSignature,
		CustomerID:        dto.CustomerID,
		DeliveryService:   dto.DeliveryService,
		ShardKey:          dto.Shardkey,
		SmID:              dto.SmID,
		DateCreated:       dto.DateCreated,
		OofShard:          dto.OofShard,
	}

	delivery := domain.Delivery{
		OrderID: order.ID,
		Name:    dto.Delivery.Name,
		Phone:   dto.Delivery.Phone,
		Zip:     dto.Delivery.Zip,
		City:    dto.Delivery.City,
		Address: dto.Delivery.Address,
		Region:  dto.Delivery.Region,
		Email:   dto.Delivery.Email,
	}

	payment := domain.Payment{
		Transaction:  order.ID,
		OrderID:      order.ID,
		RequestID:    dto.Payment.RequestID,
		Currency:     dto.Payment.Currency,
		Provider:     dto.Payment.Provider,
		Amount:       dto.Payment.Amount,
		PaymentDt:    time.Unix(dto.Payment.PaymentDt, 0),
		Bank:         dto.Payment.Bank,
		DeliveryCost: dto.Payment.DeliveryCost,
		GoodsTotal:   dto.Payment.GoodsTotal,
		CustomFee:    dto.Payment.CustomFee,
	}

	var items []domain.Item
	for _, itm := range dto.Items {
		item := domain.Item{
			ChrtID:      int64(itm.ChrtID),
			OrderID:     order.ID,
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

	return order, delivery, payment, items, nil
}

func (c *Converter) DomainToDtoOrder(fullOrder domain.FullOrder) dto.Order {
	delivery := dto.Delivery{
		Name:    fullOrder.Delivery.Name,
		Phone:   fullOrder.Delivery.Phone,
		Zip:     fullOrder.Delivery.Zip,
		City:    fullOrder.Delivery.City,
		Address: fullOrder.Delivery.Address,
		Region:  fullOrder.Delivery.Region,
		Email:   fullOrder.Delivery.Email,
	}

	payment := dto.Payment{
		Transaction:  fullOrder.Payment.Transaction.String(),
		RequestID:    fullOrder.Payment.RequestID,
		Currency:     fullOrder.Payment.Currency,
		Provider:     fullOrder.Payment.Provider,
		Amount:       fullOrder.Payment.Amount,
		PaymentDt:    fullOrder.Payment.PaymentDt.Unix(),
		Bank:         fullOrder.Payment.Bank,
		DeliveryCost: fullOrder.Payment.DeliveryCost,
		GoodsTotal:   fullOrder.Payment.GoodsTotal,
		CustomFee:    fullOrder.Payment.CustomFee,
	}

	var items []dto.Item
	for _, itm := range fullOrder.Items {
		item := dto.Item{
			ChrtID:      int(itm.ChrtID),
			TrackNumber: itm.TrackNumber,
			Price:       itm.Price,
			Rid:         itm.Rid,
			Name:        itm.Name,
			Sale:        itm.Sale,
			Size:        itm.Size,
			TotalPrice:  itm.TotalPrice,
			NmID:        int(itm.NmID),
			Brand:       itm.Brand,
			Status:      itm.Status,
		}
		items = append(items, item)
	}

	return dto.Order{
		OrderUID:          fullOrder.Order.ID.String(),
		TrackNumber:       fullOrder.Order.TrackNumber,
		Entry:             fullOrder.Order.Entry,
		Delivery:          delivery,
		Payment:           payment,
		Items:             items,
		Locale:            fullOrder.Order.Locale,
		InternalSignature: fullOrder.Order.InternalSignature,
		CustomerID:        fullOrder.Order.CustomerID,
		DeliveryService:   fullOrder.Order.DeliveryService,
		Shardkey:          fullOrder.Order.ShardKey,
		SmID:              fullOrder.Order.SmID,
		DateCreated:       fullOrder.Order.DateCreated,
		OofShard:          fullOrder.Order.OofShard,
	}
}
