package requests

import (
	"boilerplate/internal/domain"
)

type OrderRequest struct {
	OrderItems    []OrderItemRequest `json:"order_items"`
	AddressId     uint64             `json:"address_id" validate:"required"`
	Comment       string             `json:"comment"`
	ShippingPrice float64            `json:"shipping_price"`
}

type UpdateOrderRequest struct {
	AddressId uint64 `json:"address_id" validate:"required"`
	Comment   string `json:"comment"`
}

func (m UpdateOrderRequest) ToDomainModel() (interface{}, error) {
	return domain.Order{
		AddressId: m.AddressId,
		Comment:   m.Comment,
	}, nil
}

func (m OrderRequest) ToDomainModel() (interface{}, error) {
	orderItems, err := OrderItemRequest{}.ToDomainModelArray(m.OrderItems)
	if err != nil {
		return domain.Order{}, err
	}
	return domain.Order{
		AddressId:     m.AddressId,
		Comment:       m.Comment,
		ShippingPrice: m.ShippingPrice,
		OrderItems:    orderItems,
	}, nil
}
