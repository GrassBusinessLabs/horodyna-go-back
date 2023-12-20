package domain

import (
	"time"
)

type Order struct {
	Id            uint64
	Comment       string
	UserId        uint64
	AddressId     uint64
	OrderItems    []OrderItem
	ProductsPrice float64
	ShippingPrice float64
	TotalPrice    float64
	Status        bool
	CreatedDate   time.Time
	UpdatedDate   time.Time
	DeletedDate   *time.Time
}

type Orders struct {
	Items []Order
	Total uint64
	Pages uint
}

func (o Order) GetUserId() uint64 {
	return o.UserId
}
