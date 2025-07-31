package domain

import (
	"github.com/google/uuid"
	"time"
)

type FullOrder struct {
	Order    Order
	Delivery Delivery
	Payment  Payment
	Items    []Item
}

type Order struct {
	ID                uuid.UUID `db:"id"`
	TrackNumber       string    `db:"track_number"`
	Entry             string    `db:"entry"`
	Locale            string    `db:"locale"`
	InternalSignature string    `db:"internal_signature"`
	CustomerID        string    `db:"customer_id"`
	DeliveryService   string    `db:"delivery_service"`
	ShardKey          string    `db:"shardkey"`
	SmID              int       `db:"sm_id"`
	DateCreated       time.Time `db:"date_created"`
	OofShard          string    `db:"oof_shard"`
}

type Delivery struct {
	ID      int       `db:"-"`
	OrderID uuid.UUID `db:"order_id"`
	Name    string    `db:"name"`
	Phone   string    `db:"phone"`
	Zip     string    `db:"zip"`
	City    string    `db:"city"`
	Address string    `db:"address"`
	Region  string    `db:"region"`
	Email   string    `db:"email"`
}

type Payment struct {
	Transaction  uuid.UUID `db:"transaction"`
	OrderID      uuid.UUID `db:"order_id"`
	RequestID    string    `db:"request_id"`
	Currency     string    `db:"currency"`
	Provider     string    `db:"provider"`
	Amount       int       `db:"amount"`
	PaymentDt    time.Time `db:"payment_dt"`
	Bank         string    `db:"bank"`
	DeliveryCost int       `db:"delivery_cost"`
	GoodsTotal   int       `db:"goods_total"`
	CustomFee    int       `db:"custom_fee"`
}

type Item struct {
	ChrtID      int64     `db:"chrt_id"`
	OrderID     uuid.UUID `db:"order_id"`
	TrackNumber string    `db:"track_number"`
	Price       int       `db:"price"`
	Rid         string    `db:"rid"`
	Name        string    `db:"name"`
	Sale        int       `db:"sale"`
	Size        string    `db:"size"`
	TotalPrice  int       `db:"total_price"`
	NmID        int64     `db:"nm_id"`
	Brand       string    `db:"brand"`
	Status      int       `db:"status"`
}
