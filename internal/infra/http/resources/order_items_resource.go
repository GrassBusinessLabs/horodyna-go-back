package resources

import (
	"boilerplate/internal/app"
	"boilerplate/internal/domain"
)

type OrderItemDto struct {
	Id         uint64         `json:"id"`
	OrderId    uint64         `json:"order_id"`
	Offer      OfferDto       `json:"offer"`
	Title      string         `json:"title"`
	Price      float64        `json:"price"`
	TotalPrice float64        `json:"total_price"`
	Amount     uint32         `json:"amount"`
	Farm       FarmWithOutDto `json:"farm"`
}

func (d OrderItemDto) DomainToDto(o domain.OrderItem, imageModelService app.ImageModelService) OrderItemDto {
	return OrderItemDto{
		Id:         o.Id,
		OrderId:    o.Order.Id,
		Offer:      OfferDto{}.DomainToDto(o.Offer, imageModelService),
		Title:      o.Title,
		Price:      o.Price,
		TotalPrice: o.TotalPrice,
		Amount:     o.Amount,
		Farm:       FarmWithOutDto{}.DomainToDto(o.Farm),
	}
}
