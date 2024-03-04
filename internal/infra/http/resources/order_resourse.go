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
	User            UserDto `json:"user"`
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
	User          UserDto        `json:"user"`
	ProductPrice  float64        `json:"product_price"`
	ShippingPrice float64        `json:"shipping_price"`
	TotalPrice    float64        `json:"total_price"`
	CreatedDate   string         `json:"created_data"`
}

type SplitedOrdersDto struct {
	SplitedOrders map[uint64]OrderDtoWithOrderItems `json:"splited_orders"`
}

type OrdersPercentageDto struct {
	Total            float64            `json:"total"`
	OrdersPercentage map[uint64]float64 `json:"orders_percentage"`
}

type OrdersDto struct {
	Items []OrderDto `json:"items"`
	Pages uint       `json:"pages"`
	Total uint64     `json:"total"`
}

type OrdersDtoWithOrderItems struct {
	Items []OrderDtoWithOrderItems `json:"items"`
	Pages uint                     `json:"pages"`
	Total uint64                   `json:"total"`
}

func (d OrderDtoWithOrderItems) DomainToDto(order domain.Order) OrderDtoWithOrderItems {
	orderItems := make([]OrderItemDto, len(order.OrderItems))
	for i, item := range order.OrderItems {
		orderItems[i] = OrderItemDto{}.DomainToDto(item)
	}

	return OrderDtoWithOrderItems{
		Id:            order.Id,
		OrderItems:    orderItems,
		Status:        string(order.Status),
		Comment:       order.Comment,
		Address:       order.Address,
		User:          UserDto{}.DomainToDto(order.User),
		ProductPrice:  order.ProductsPrice,
		ShippingPrice: order.ShippingPrice,
		TotalPrice:    order.TotalPrice,
		CreatedDate:   order.CreatedDate.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (d OrderDtoWithOrderItems) DomainToDtoCollection(orders []domain.Order) []OrderDtoWithOrderItems {
	result := make([]OrderDtoWithOrderItems, len(orders))

	for i := range orders {
		result[i] = d.DomainToDto(orders[i])
	}

	return result
}

func (d OrderDtoWithOrderItems) DomainToDtoPaginatedCollection(orders domain.Orders) OrdersDtoWithOrderItems {
	result := make([]OrderDtoWithOrderItems, len(orders.Items))

	for i := range orders.Items {
		result[i] = d.DomainToDto(orders.Items[i])
	}

	return OrdersDtoWithOrderItems{Items: result, Pages: orders.Pages, Total: orders.Total}
}

func (d OrderDto) DomainToDto(order domain.Order) OrderDto {
	return OrderDto{
		Id:              order.Id,
		OrderItemsCount: order.OrderItemsCount,
		Status:          string(order.Status),
		Comment:         order.Comment,
		Address:         order.Address,
		User:            UserDto{}.DomainToDto(order.User),
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

func (d SplitedOrdersDto) DomainToDto(splitedOrders map[uint64]domain.Order) SplitedOrdersDto {
	splitedOrdersDto := make(map[uint64]OrderDtoWithOrderItems, len(splitedOrders))

	for key, order := range splitedOrders {
		splitedOrdersDto[key] = OrderDtoWithOrderItems{}.DomainToDto(order)
	}

	return SplitedOrdersDto{SplitedOrders: splitedOrdersDto}
}

func (d OrdersPercentageDto) DataToDto(ordersPercentage map[uint64]float64, total float64) OrdersPercentageDto {
	return OrdersPercentageDto{
		Total:            total,
		OrdersPercentage: ordersPercentage,
	}
}
