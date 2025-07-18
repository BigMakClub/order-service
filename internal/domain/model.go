package domain

import (
	"github.com/google/uuid"
	"time"
)

//
//type Order struct {
//	OrderId           uuid.UUID
//	TrackNumber       string
//	Entry             string
//	Locale            string
//	InternalSignature string
//	CustomerId        string
//	DeliveryService   string
//	ShardKey          int16
//	SmId              int
//	DateCreated       time.Time
//
//	Delivery Delivery
//	Payment  Payment
//	Items    []Item
//}
//
//type Delivery struct {
//	Name    string
//	Phone   string
//	Zip     string
//	City    string
//	Address string
//	Region  string
//	Email   string
//}
//
//type Payment struct {
//	TransactionId string
//	RequestId     string
//	Currency      string
//	Provider      string
//	Amount        int64
//	PaymentDt     int64
//	Bank          string
//	DeliveryCost  int64
//	GoodsTotal    int64
//	CustomFee     int64
//}
//
//type Item struct {
//	ID          int64
//	ChrtId      int64
//	TrackNumber string
//	Price       int64
//	RID         string
//	Name        string
//	Sale        int
//	Size        string
//	TotalPrice  int64
//	NmID        int64
//	Brand       string
//	Status      string
//}

type Order struct {
	OrderId           uuid.UUID `json:"order_uid"`
	TrackNumber       string    `json:"track_number"`
	Entry             string    `json:"entry"`
	Locale            string    `json:"locale"`
	InternalSignature string    `json:"internal_signature"`
	CustomerId        string    `json:"customer_id"`
	DeliveryService   string    `json:"delivery_service"`
	ShardKey          int64     `json:"shardkey"`
	SmId              int       `json:"sm_id"`
	DateCreated       time.Time `json:"date_created"`

	Delivery Delivery `json:"delivery"`
	Payment  Payment  `json:"payment"`
	Items    []Item   `json:"items"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	TransactionId string `json:"transaction_id"`
	RequestId     string `json:"request_id"`
	Currency      string `json:"currency"`
	Provider      string `json:"provider"`
	Amount        int64  `json:"amount"`
	PaymentDt     int64  `json:"payment_dt"`
	Bank          string `json:"bank"`
	DeliveryCost  int64  `json:"delivery_cost"`
	GoodsTotal    int64  `json:"goods_total"`
	CustomFee     int64  `json:"custom_fee"`
}

type Item struct {
	ID          int64  `json:"-"`
	ChrtId      int64  `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int64  `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int64  `json:"total_price"`
	NmID        int64  `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      string `json:"status"`
}
