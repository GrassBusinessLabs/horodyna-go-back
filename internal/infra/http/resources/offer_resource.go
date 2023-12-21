package resources

import (
	"boilerplate/internal/domain"
)

type OfferDto struct {
	Id          uint64  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
	Unit        string  `json:"unit"`
	Stock       uint    `json:"stock"`
	Status      bool    `json:"status"`
	Cover       string  `json:"image"`
	UserId      uint64  `json:"user_id"`
	FarmId      uint64  `json:"farm_id"`
}

type OffersDto struct {
	Items []OfferDto `json:"items"`
	Pages uint       `json:"pages"`
	Total uint64     `json:"total"`
}

func (d OfferDto) DomainToDto(offer domain.Offer) OfferDto {
	return OfferDto{
		Id:          offer.Id,
		Title:       offer.Title,
		Description: offer.Description,
		Category:    offer.Category,
		Price:       offer.Price,
		Unit:        offer.Unit,
		Stock:       offer.Stock,
		Cover:       offer.Cover.Name,
		Status:      offer.Status,
		FarmId:      offer.Farm.Id,
		UserId:      offer.UserId,
	}
}

func (d OfferDto) DomainToDtoPaginatedCollection(offers domain.Offers) OffersDto {
	result := make([]OfferDto, len(offers.Items))

	for i := range offers.Items {
		result[i] = d.DomainToDto(offers.Items[i])
	}

	return OffersDto{Items: result, Pages: offers.Pages, Total: offers.Total}
}
