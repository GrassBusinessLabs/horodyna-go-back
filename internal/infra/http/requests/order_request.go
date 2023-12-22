package requests

import (
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/database"
)

type OrderRequest struct {
	OrderItems    []OrderItemRequest `json:"order_items"`
	AddressId     uint64             `json:"address_id" validate:"required"`
	Comment       string             `json:"comment"`
	ShippingPrice float64            `json:"shipping_price"`
}

type UpdateOrderRequest struct {
	AddressId     uint64  `json:"address_id" validate:"required"`
	Comment       string  `json:"comment"`
	ShippingPrice float64 `json:"shipping_price"`
	Status        string  `json:"status"`
}

func (m UpdateOrderRequest) ToDomainModel() (interface{}, error) {
	status, err := database.FindOrderStatus(m.Status)
	if err != nil {
		return domain.Order{}, err
	}

	return domain.Order{
		AddressId:     m.AddressId,
		Comment:       m.Comment,
		ShippingPrice: m.ShippingPrice,
		Status:        status,
	}, nil
}

func (m OrderRequest) ToDomainModel() (interface{}, error) {
	var err error
	new := make([]domain.OrderItem, len(m.OrderItems))
	for i, item := range m.OrderItems {
		new[i], err = item.ToDomainModelNotInterface()
		if err != nil {
			return domain.Order{}, err
		}
	}

	return domain.Order{
		AddressId:     m.AddressId,
		Comment:       m.Comment,
		ShippingPrice: m.ShippingPrice,
		OrderItems:    new,
		Status:        domain.DRAFT,
	}, nil
}
