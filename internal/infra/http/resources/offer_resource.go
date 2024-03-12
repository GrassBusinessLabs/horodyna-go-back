package resources

import (
	"boilerplate/internal/app"
	"boilerplate/internal/domain"
	"math"
)

type OfferDto struct {
	Id               uint64      `json:"id"`
	Title            string      `json:"title"`
	Description      string      `json:"description"`
	Category         string      `json:"category"`
	Price            float64     `json:"price"`
	Unit             string      `json:"unit"`
	Stock            uint        `json:"stock"`
	Status           bool        `json:"status"`
	Cover            string      `json:"image"`
	AdditionalImages []ImageMDto `json:"additional_images"`
	User             UserDto     `json:"user"`
	FarmId           uint64      `json:"farm_id"`
}

type OffersDto struct {
	Items []OfferDto `json:"items"`
	Pages uint       `json:"pages"`
	Total uint64     `json:"total"`
}

func (d OfferDto) DomainToDto(offer domain.Offer, imageModelService app.ImageModelService) OfferDto {
	additionalImages, _ := imageModelService.FindAll("offers", offer.Id)
	return OfferDto{
		Id:               offer.Id,
		Title:            offer.Title,
		Description:      offer.Description,
		Category:         offer.Category,
		Price:            math.Round(offer.Price*100) / 100,
		Unit:             offer.Unit,
		Stock:            offer.Stock,
		Cover:            offer.Cover.Name,
		AdditionalImages: ImageMDto{}.DomainToDtoMass(additionalImages).Items,
		Status:           offer.Status,
		FarmId:           offer.Farm.Id,
		User:             UserDto{}.DomainToDto(offer.User),
	}
}

func (d OfferDto) DomainToDtoPaginatedCollection(offers domain.Offers, imageModelService app.ImageModelService) OffersDto {
	result := make([]OfferDto, len(offers.Items))

	for i := range offers.Items {
		result[i] = d.DomainToDto(offers.Items[i], imageModelService)
	}

	return OffersDto{Items: result, Pages: offers.Pages, Total: offers.Total}
}
