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
	OrderId     uint64
	OfferId     uint64
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

type OrderItems struct {
	Items []OrderItem
	Total uint64
	Pages uint
}
