package domain

import (
	"time"
)

type OrderItem struct {
	Id          uint64
	Title       string
	Price       float64
	TotalPrice  float64
	Amount      uint32
	Order       Order
	Offer       Offer
	Farm        Farm
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

type OrderItems struct {
	Items []OrderItem
	Total uint64
	Pages uint
}

func (o OrderItem) GetUserId() uint64 {
	return o.Order.User.Id
}
