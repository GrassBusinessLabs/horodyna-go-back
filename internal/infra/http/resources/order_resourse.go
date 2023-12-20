package resources

import (
	"boilerplate/internal/domain"
)

type OrderDto struct {
	Id            uint64         `json:"id"`
	OrderItems    []OrderItemDto `json:"order_items"`
	Status        bool           `json:"status"`
	Comment       string         `json:"comment"`
	AddressId     uint64         `json:"address_id"`
	UserId        uint64         `json:"user_id"`
	ProductPrice  float64        `json:"product_price"`
	ShippingPrice float64        `json:"shipping_price"`
	TotalPrice    float64        `json:"total_price"`
}

func (d OrderDto) DomainToDto(order domain.Order) OrderDto {
	OrderItemsDto := make([]OrderItemDto, len(order.OrderItems))

	for i, item := range order.OrderItems {
		OrderItemsDto[i] = OrderItemDto{}.DomainToDto(item)
	}

	return OrderDto{
		Id:            order.Id,
		OrderItems:    OrderItemsDto,
		Status:        order.Status,
		Comment:       order.Comment,
		AddressId:     order.AddressId,
		UserId:        order.UserId,
		ProductPrice:  order.ProductsPrice,
		ShippingPrice: order.ShippingPrice,
		TotalPrice:    order.TotalPrice,
	}
}
