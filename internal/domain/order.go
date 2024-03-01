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
	User            User
	Address         *string
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
	return o.User.Id
}

func GetActiveOrderStatuses() []OrderStatus {
	return []OrderStatus{SUBMITTED, APPROVED, SHIPPING}
}

func (o Order) IsReceiverStatus(stasus OrderStatus) bool {
	receiverStatuses := []OrderStatus{SUBMITTED, COMPLETED}
	for _, receiverStatus := range receiverStatuses {
		if stasus == receiverStatus {
			return true
		}
	}

	return false
}

func (o Order) IsFarmerStatus(stasus OrderStatus) bool {
	farmerStatuses := []OrderStatus{APPROVED, DECLINED, SHIPPING}
	for _, farmerStatus := range farmerStatuses {
		if stasus == farmerStatus {
			return true
		}
	}

	return false
}
