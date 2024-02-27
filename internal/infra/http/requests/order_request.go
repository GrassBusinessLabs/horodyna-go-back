package requests

import (
	"boilerplate/internal/domain"
)

type OrderRequest struct {
	OrderItems    []OrderItemRequest `json:"order_items"`
	Address       *string            `json:"address" validate:"required"`
	Comment       string             `json:"comment"`
	ShippingPrice float64            `json:"shipping_price"`
	PostOffice    *string            `json:"post_office"`
	Ttn           *string            `json:"ttn"`
}

type UpdateOrderRequest struct {
	Address       *string `json:"address" validate:"required"`
	Comment       string  `json:"comment"`
	ShippingPrice float64 `json:"shipping_price"`
	PostOffice    *string `json:"post_office"`
	Ttn           *string `json:"ttn"`
}

type OrderStatusRequest struct {
	Status string `json:"status"`
}

func (m UpdateOrderRequest) ToDomainModel() (interface{}, error) {
	return domain.Order{
		Address:       m.Address,
		Comment:       m.Comment,
		ShippingPrice: m.ShippingPrice,
		PostOffice:    m.PostOffice,
		Ttn:           m.Ttn,
	}, nil
}

func (m OrderRequest) ToDomainModel() (interface{}, error) {
	orderItems, err := OrderItemRequest{}.ToDomainModelArray(m.OrderItems)
	if err != nil {
		return domain.Order{}, err
	}

	return domain.Order{
		Address:       m.Address,
		Comment:       m.Comment,
		ShippingPrice: m.ShippingPrice,
		OrderItems:    orderItems,
		PostOffice:    m.PostOffice,
		Ttn:           m.Ttn,
	}, nil
}

func (m OrderStatusRequest) ToDomainModel() (interface{}, error) {
	return domain.Order{
		Status: domain.OrderStatus(m.Status),
	}, nil
}
