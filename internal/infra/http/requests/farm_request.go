package requests

import (
	"boilerplate/internal/domain"
)

type FarmRequest struct {
	Name      string  `json:"name" validate:"required,gte=1,max=40"`
	City      string  `json:"city" validate:"required"`
	Address   string  `json:"address" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

func (m FarmRequest) ToDomainModel() (interface{}, error) {
	return domain.Farm{
		Name:      m.Name,
		City:      m.City,
		Address:   m.Address,
		Longitude: m.Longitude,
		Latitude:  m.Latitude,
	}, nil
}
