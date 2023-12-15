package requests

import (
	"boilerplate/internal/domain"
)

type ImageRequest struct {
	Name string `json:"name" validate:"required"`
	Data string `json:"data" validate:"required"`
}

type OfferRequest struct {
	Title       string        `json:"title" validate:"required,gte=1,max=40"`
	Description string        `json:"description" validate:"required"`
	Category    string        `json:"category" validate:"required"`
	Price       float64       `json:"price" validate:"required"`
	Unit        string        `json:"unit" validate:"required"`
	Stock       uint          `json:"stock" validate:"required"`
	FarmId      uint64        `json:"farm_id" validate:"required"`
	Status      bool          `json:"status"`
	Cover       *ImageRequest `json:"image"`
}

func (m ImageRequest) ToDomainModel() domain.Image {
	return domain.Image{
		Name: m.Name,
		Data: m.Data,
	}
}

func (m OfferRequest) ToDomainModel() (interface{}, error) {
	var img domain.Image
	if m.Cover != nil {
		img = m.Cover.ToDomainModel()
	}

	return domain.Offer{
		Title:       m.Title,
		Description: m.Description,
		Category:    m.Category,
		Price:       m.Price,
		Unit:        m.Unit,
		Stock:       m.Stock,
		Status:      m.Status,
		Farm:        domain.Farm{Id: m.FarmId},
		Cover:       img,
	}, nil
}
