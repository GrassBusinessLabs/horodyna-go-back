package resources

import (
	"boilerplate/internal/app"
	"boilerplate/internal/domain"
	"log"
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
	Farm        FarmDto `json:"farm"`
}

type OffersDto struct {
	Items []OfferDto `json:"items"`
	Pages uint       `json:"pages"`
	Total uint64     `json:"total"`
}

func (d OfferDto) DomainToDto(offer domain.Offer, fs app.FarmService, us app.UserService) OfferDto {
	farm, err := fs.FindById(offer.FarmId)
	if err != nil {
		log.Println(err)
	}

	return OfferDto{
		Id:          offer.Id,
		Title:       offer.Title,
		Description: offer.Description,
		Category:    offer.Category,
		Price:       offer.Price,
		Unit:        offer.Unit,
		Stock:       offer.Stock,
		Cover:       offer.Cover,
		Status:      offer.Status,
		Farm:        FarmDto{}.DomainToDto(farm, us),
	}
}

func (d OfferDto) DomainToDtoPaginatedCollection(offers domain.Offers, pag domain.Pagination, fs app.FarmService, us app.UserService) OffersDto {
	result := make([]OfferDto, len(offers.Items))

	for i := range offers.Items {
		result[i] = d.DomainToDto(offers.Items[i], fs, us)
	}

	return OffersDto{Items: result, Pages: offers.Pages, Total: offers.Total}
}
