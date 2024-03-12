package resources

import (
	"boilerplate/internal/app"
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
	PostOffice      *string `json:"post_office"`
	PostOfficeCity  *string `json:"post_office_city"`
	Ttn             *string `json:"ttn"`
	CreatedDate     string  `json:"created_data"`
}

type OrdersDto struct {
	Items []OrderDto `json:"items"`
	Pages uint       `json:"pages"`
	Total uint64     `json:"total"`
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
		PostOffice:      order.PostOffice,
		PostOfficeCity:  order.PostOfficeCity,
		Ttn:             order.Ttn,
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

type OrderDtoWithOrderItems struct {
	Id             uint64         `json:"id"`
	OrderItems     []OrderItemDto `json:"order_items"`
	Status         string         `json:"status"`
	Comment        string         `json:"comment"`
	Address        *string        `json:"address"`
	User           UserDto        `json:"user"`
	ProductPrice   float64        `json:"product_price"`
	ShippingPrice  float64        `json:"shipping_price"`
	TotalPrice     float64        `json:"total_price"`
	PostOffice     *string        `json:"post_office"`
	PostOfficeCity *string        `json:"post_office_city"`
	Ttn            *string        `json:"ttn"`
	CreatedDate    string         `json:"created_data"`
}

func (d OrderDtoWithOrderItems) DomainToDto(order domain.Order, imageModelService app.ImageModelService) OrderDtoWithOrderItems {
	orderItems := make([]OrderItemDto, len(order.OrderItems))
	for i, item := range order.OrderItems {
		orderItems[i] = OrderItemDto{}.DomainToDto(item, imageModelService)
	}

	return OrderDtoWithOrderItems{
		Id:             order.Id,
		OrderItems:     orderItems,
		Status:         string(order.Status),
		Comment:        order.Comment,
		Address:        order.Address,
		User:           UserDto{}.DomainToDto(order.User),
		ProductPrice:   order.ProductsPrice,
		ShippingPrice:  order.ShippingPrice,
		TotalPrice:     order.TotalPrice,
		PostOffice:     order.PostOffice,
		PostOfficeCity: order.PostOfficeCity,
		Ttn:            order.Ttn,
		CreatedDate:    order.CreatedDate.Format("2006-01-02T15:04:05Z07:00"),
	}
}

type OrdersDtoWithOrderItems struct {
	Items []OrderDtoWithOrderItems `json:"items"`
	Pages uint                     `json:"pages"`
	Total uint64                   `json:"total"`
}

func (d OrderDtoWithOrderItems) DomainToDtoCollection(orders []domain.Order, imageModelService app.ImageModelService) []OrderDtoWithOrderItems {
	result := make([]OrderDtoWithOrderItems, len(orders))

	for i := range orders {
		result[i] = d.DomainToDto(orders[i], imageModelService)
	}

	return result
}

func (d OrderDtoWithOrderItems) DomainToDtoPaginatedCollection(orders domain.Orders, imageModelService app.ImageModelService) OrdersDtoWithOrderItems {
	result := make([]OrderDtoWithOrderItems, len(orders.Items))

	for i := range orders.Items {
		result[i] = d.DomainToDto(orders.Items[i], imageModelService)
	}

	return OrdersDtoWithOrderItems{Items: result, Pages: orders.Pages, Total: orders.Total}
}

type OrderDtoWithPercentage struct {
	Id              uint64   `json:"id"`
	OrderItemsCount uint64   `json:"order_items_count"`
	Status          string   `json:"status"`
	Comment         string   `json:"comment"`
	Address         *string  `json:"address"`
	User            UserDto  `json:"user"`
	ProductPrice    float64  `json:"product_price"`
	ShippingPrice   float64  `json:"shipping_price"`
	TotalPrice      float64  `json:"total_price"`
	PostOffice      *string  `json:"post_office"`
	PostOfficeCity  *string  `json:"post_office_city"`
	Ttn             *string  `json:"ttn"`
	CreatedDate     string   `json:"created_data"`
	Percenatge      *float64 `json:"percentage"`
}

type OrdersDtoWithPercentage struct {
	Total float64                  `json:"total"`
	Items []OrderDtoWithPercentage `json:"items"`
}

func (d OrdersDtoWithPercentage) DomainToDto(orders []domain.Order, total float64) OrdersDtoWithPercentage {
	return OrdersDtoWithPercentage{
		Total: total,
		Items: OrderDtoWithPercentage{}.DomainToDtoCollection(orders),
	}
}

func (d OrderDtoWithPercentage) DomainToDto(order domain.Order) OrderDtoWithPercentage {
	return OrderDtoWithPercentage{
		Id:              order.Id,
		OrderItemsCount: order.OrderItemsCount,
		Status:          string(order.Status),
		Comment:         order.Comment,
		Address:         order.Address,
		User:            UserDto{}.DomainToDto(order.User),
		ProductPrice:    order.ProductsPrice,
		ShippingPrice:   order.ShippingPrice,
		TotalPrice:      order.TotalPrice,
		PostOffice:      order.PostOffice,
		PostOfficeCity:  order.PostOfficeCity,
		Ttn:             order.Ttn,
		CreatedDate:     order.CreatedDate.Format("2006-01-02T15:04:05Z07:00"),
		Percenatge:      order.Percentage,
	}
}

func (d OrderDtoWithPercentage) DomainToDtoCollection(orders []domain.Order) []OrderDtoWithPercentage {
	result := make([]OrderDtoWithPercentage, len(orders))

	for i := range orders {
		result[i] = d.DomainToDto(orders[i])
	}

	return result
}

type SplitedOrdersDto struct {
	SplitedOrders map[uint64]OrderDtoWithOrderItems `json:"splited_orders"`
}

func (d SplitedOrdersDto) DomainToDto(splitedOrders map[uint64]domain.Order, imageModelService app.ImageModelService) SplitedOrdersDto {
	splitedOrdersDto := make(map[uint64]OrderDtoWithOrderItems, len(splitedOrders))

	for key, order := range splitedOrders {
		splitedOrdersDto[key] = OrderDtoWithOrderItems{}.DomainToDto(order, imageModelService)
	}

	return SplitedOrdersDto{SplitedOrders: splitedOrdersDto}
}
