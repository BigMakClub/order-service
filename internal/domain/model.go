package domain

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderId           uuid.UUID
	TrackNumber       string
	Entry             string
	Locale            string
	InternalSignature string
	CustomerId        string
	DeliveryService   string
	ShardKey          string
	SmId              int
	DateCreated       time.Time

	Delivery Delivery
	Payment  Payment
	Items    []Item
}

type Delivery struct {
	Name    string
	Phone   string
	Zip     string
	City    string
	Address string
	Region  string
	Email   string
}

type Payment struct {
	TransactionId string
	RequestId     string
	Currency      string
	Provider      string
	Amount        int64
	PaymentDt     int64
	Bank          string
	DeliveryCost  int64
	GoodsTotal    int64
	CustomFee     int64
}

type Item struct {
	ID          int64
	ChrtId      int64
	TrackNumber string
	Price       int64
	RID         string
	Name        string
	Sale        int
	Size        string
	TotalPrice  int64
	NmID        int64
	Brand       string
	Status      string
}
