package requests

import (
	"boilerplate/internal/domain"
)

type AddressRequest struct {
	UserID  uint64  `json:"user_id" validate:"required"`
	Title   string  `json:"title" validate:"required"`
	City    string  `json:"city" validate:"required"`
	Country string  `json:"country" validate:"required"`
	Address string  `json:"address" validate:"required"`
	Lat     float64 `json:"lat" validate:"required"`
	Lon     float64 `json:"lon" validate:"required"`
}

func (m AddressRequest) ToDomainModel() (interface{}, error) {
	return domain.Address{
		UserID:  m.UserID,
		Title:   m.Title,
		City:    m.City,
		Country: m.Country,
		Address: m.Address,
		Lat:     m.Lat,
		Lon:     m.Lon,
	}, nil
}
