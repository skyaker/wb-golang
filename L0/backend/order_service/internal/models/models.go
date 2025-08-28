package models

import "time"

type Order struct {
	OrderUID    string `json:"order_uid"    validate:"required,uuid4"`
	TrackNumber string `json:"track_number" validate:"required"`
	Entry       string `json:"entry"        validate:"required"`

	Delivery *Delivery `json:"delivery" validate:"required"`
	Payment  *Payment  `json:"payment"  validate:"required"`
	Items    []Item    `json:"items"    validate:"required,min=1,dive"`

	Locale       string    `json:"locale"             validate:"required"`
	InternalSig  string    `json:"internal_signature"`
	CustomerID   string    `json:"customer_id"        validate:"required,uuid4"`
	DeliveryServ string    `json:"delivery_service"`
	ShardKey     string    `json:"shardkey"           validate:"required"`
	SMID         int       `json:"sm_id"              validate:"gt=0"`
	DateCreated  time.Time `json:"date_created"       validate:"required"`
	OOFShard     string    `json:"oof_shard"          validate:"required"`
}

type AggregatedOrder struct {
	Order    Order    `json:"order"`
	Delivery Delivery `json:"delivery"`
	Payment  Payment  `json:"payment"`
	Items    []Item   `json:"items"`
}

type Delivery struct {
	DeliveryUID string
	OrderUID    string
	Name        string `json:"name"    validate:"required"`
	Phone       string `json:"phone"`
	Zip         string `json:"zip"`
	City        string `json:"city"`
	Address     string `json:"address" validate:"required"`
	Region      string `json:"region"`
	Email       string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"   validate:"required,uuid4"`
	RequestID    string `json:"request_id"    validate:"omitempty,uuid4"`
	Currency     string `json:"currency"      validate:"required,len=3"`
	Provider     string `json:"provider"      validate:"required"`
	Amount       int    `json:"amount"        validate:"gte=0"`
	PaymentDT    int    `json:"payment_dt"    validate:"required"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost" validate:"gte=0"`
	GoodsTotal   int    `json:"goods_total"   validate:"gte=0"`
	CustomFee    int    `json:"custom_fee"    validate:"gte=0"`
}

type Item struct {
	ChrtID      int    `json:"chrt_id"      validate:"required"`
	TrackNumber string `json:"track_number" validate:"required"`
	Price       int    `json:"price"        validate:"gte=0"`
	Rid         string `json:"rid"          validate:"required,uuid4"`
	Name        string `json:"name"         validate:"required"`
	Sale        int    `json:"sale"         validate:"gte=0"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"  validate:"gte=0"`
	NmID        int    `json:"nm_id"        validate:"required"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"       validate:"gte=0"`
}
