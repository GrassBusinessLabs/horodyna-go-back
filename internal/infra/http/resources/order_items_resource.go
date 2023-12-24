package resources

import (
	"boilerplate/internal/domain"
)

type OrderItemDto struct {
	Id         uint64  `json:"id"`
	OrderId    uint64  `json:"order_id"`
	OfferId    uint64  `json:"offer_id"`
	Title      string  `json:"title"`
	Price      float64 `json:"price"`
	TotalPrice float64 `json:"total_price"`
	Amount     uint32  `json:"amount"`
}

func (d OrderItemDto) DomainToDto(o domain.OrderItem) OrderItemDto {
	return OrderItemDto{
		Id:         o.Id,
		OrderId:    o.Order.Id,
		OfferId:    o.OfferId,
		Title:      o.Title,
		Price:      o.Price,
		TotalPrice: o.TotalPrice,
		Amount:     o.Amount,
	}
}
