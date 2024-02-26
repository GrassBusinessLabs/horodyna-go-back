package domain

import (
	"time"
)

type OrderStatus string

const (
	DRAFT     OrderStatus = "DRAFT"
	SUBMITTED OrderStatus = "SUBMITTED"
	APPROVED  OrderStatus = "APPROVED"
	DECLINED  OrderStatus = "DECLINED"
	SHIPPING  OrderStatus = "SHIPPING"
	COMPLETED OrderStatus = "COMPLETED"
)

type Order struct {
	Id              uint64
	Comment         string
	UserId          uint64
	Address         Address
	OrderItems      []OrderItem
	OrderItemsCount uint64
	ProductsPrice   float64
	ShippingPrice   float64
	TotalPrice      float64
	Status          OrderStatus
	PostOffice      *string
	Ttn             *string
	CreatedDate     time.Time
	UpdatedDate     time.Time
	DeletedDate     *time.Time
}

type Orders struct {
	Items []Order
	Total uint64
	Pages uint
}

func (o Order) GetUserId() uint64 {
	return o.UserId
}

func (o Order) IsOrderStatusValid(oldStatus OrderStatus, newStatus OrderStatus) bool {
	if newStatus == SUBMITTED {
		return true
	} else if newStatus == SHIPPING {
		return oldStatus == SUBMITTED
	} else if newStatus == COMPLETED {
		return oldStatus == SHIPPING
	}

	return false
}
