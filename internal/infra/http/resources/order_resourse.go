package resources

import (
	"boilerplate/internal/domain"
)

type OrderDto struct {
	Id            uint64         `json:"id"`
	OrderItems    []OrderItemDto `json:"order_items"`
	Status        string         `json:"status"`
	Comment       string         `json:"comment"`
	AddressId     uint64         `json:"address_id"`
	UserId        uint64         `json:"user_id"`
	ProductPrice  float64        `json:"product_price"`
	ShippingPrice float64        `json:"shipping_price"`
	TotalPrice    float64        `json:"total_price"`
}

type OrdersDto struct {
	Items []OrderDto `json:"items"`
	Pages uint       `json:"pages"`
	Total uint64     `json:"total"`
}

func (d OrderDto) DomainToDto(order domain.Order) OrderDto {
	OrderItemsDto := make([]OrderItemDto, len(order.OrderItems))
	for i, item := range order.OrderItems {
		OrderItemsDto[i] = OrderItemDto{}.DomainToDto(item)
	}

	return OrderDto{
		Id:            order.Id,
		OrderItems:    OrderItemsDto,
		Status:        string(order.Status),
		Comment:       order.Comment,
		AddressId:     order.AddressId,
		UserId:        order.UserId,
		ProductPrice:  order.ProductsPrice,
		ShippingPrice: order.ShippingPrice,
		TotalPrice:    order.TotalPrice,
	}
}

func (d OrderDto) DomainToDtoPaginatedCollection(orders domain.Orders) OrdersDto {
	result := make([]OrderDto, len(orders.Items))

	for i := range orders.Items {
		result[i] = d.DomainToDto(orders.Items[i])
	}

	return OrdersDto{Items: result, Pages: orders.Pages, Total: orders.Total}
}
