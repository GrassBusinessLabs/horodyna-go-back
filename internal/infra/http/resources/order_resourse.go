package resources

import (
	"boilerplate/internal/domain"
)

type OrderDto struct {
	Id              uint64  `json:"id"`
	OrderItemsCount uint64  `json:"order_items_count"`
	Status          string  `json:"status"`
	Comment         string  `json:"comment"`
	Address         *string `json:"address"`
	UserId          uint64  `json:"user_id"`
	ProductPrice    float64 `json:"product_price"`
	ShippingPrice   float64 `json:"shipping_price"`
	TotalPrice      float64 `json:"total_price"`
	CreatedDate     string  `json:"created_data"`
}

type OrderDtoWithOrderItems struct {
	Id            uint64         `json:"id"`
	OrderItems    []OrderItemDto `json:"order_items"`
	Status        string         `json:"status"`
	Comment       string         `json:"comment"`
	Address       *string        `json:"address"`
	UserId        uint64         `json:"user_id"`
	ProductPrice  float64        `json:"product_price"`
	ShippingPrice float64        `json:"shipping_price"`
	TotalPrice    float64        `json:"total_price"`
	CreatedDate   string         `json:"created_data"`
}

type OrdersDto struct {
	Items []OrderDto `json:"items"`
	Pages uint       `json:"pages"`
	Total uint64     `json:"total"`
}

func (d OrderDtoWithOrderItems) DomainToDto(order domain.Order, ori []domain.OrderItem, imageDto ImageMDto) OrderDtoWithOrderItems {
	orderItems := make([]OrderItemDto, len(ori))
	for i, item := range ori {
		orderItems[i] = OrderItemDto{}.DomainToDto(item, imageDto)
	}

	return OrderDtoWithOrderItems{
		Id:            order.Id,
		OrderItems:    orderItems,
		Status:        string(order.Status),
		Comment:       order.Comment,
		Address:       order.Address,
		UserId:        order.UserId,
		ProductPrice:  order.ProductsPrice,
		ShippingPrice: order.ShippingPrice,
		TotalPrice:    order.TotalPrice,
		CreatedDate:   order.CreatedDate.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (d OrderDtoWithOrderItems) DomainWithOrderItemsToDtoCollection(orders []domain.Order, imageDto ImageMDto) []OrderDtoWithOrderItems {
	result := make([]OrderDtoWithOrderItems, len(orders))

	for i := range orders {
		result[i] = d.DomainToDto(orders[i], orders[i].OrderItems, imageDto)
	}

	return result
}

func (d OrderDto) DomainToDto(order domain.Order) OrderDto {
	return OrderDto{
		Id:              order.Id,
		OrderItemsCount: order.OrderItemsCount,
		Status:          string(order.Status),
		Comment:         order.Comment,
		Address:         order.Address,
		UserId:          order.UserId,
		ProductPrice:    order.ProductsPrice,
		ShippingPrice:   order.ShippingPrice,
		TotalPrice:      order.TotalPrice,
		CreatedDate:     order.CreatedDate.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (d OrderDto) DomainToDtoCollection(orders []domain.Order) []OrderDto {
	result := make([]OrderDto, len(orders))

	for i := range orders {
		result[i] = d.DomainToDto(orders[i])
	}

	return result
}

func (d OrderDto) DomainToDtoPaginatedCollection(orders domain.Orders) OrdersDto {
	result := make([]OrderDto, len(orders.Items))

	for i := range orders.Items {
		result[i] = d.DomainToDto(orders.Items[i])
	}

	return OrdersDto{Items: result, Pages: orders.Pages, Total: orders.Total}
}
