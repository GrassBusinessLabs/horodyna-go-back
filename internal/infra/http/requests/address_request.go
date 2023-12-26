package requests

import (
	"boilerplate/internal/domain"
)

type AddressRequest struct {
	Title   string  `json:"title" validate:"required"`
	City    string  `json:"city" validate:"required"`
	Country string  `json:"country" validate:"required"`
	Address string  `json:"address" validate:"required"`
	Lat     float64 `json:"lat" validate:"required"`
	Lon     float64 `json:"lon" validate:"required"`
}

func (m AddressRequest) ToDomainModel() (interface{}, error) {
	return domain.Address{
		Title:   m.Title,
		City:    m.City,
		Country: m.Country,
		Address: m.Address,
		Lat:     m.Lat,
		Lon:     m.Lon,
	}, nil
}
