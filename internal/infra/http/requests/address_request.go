package requests

import (
	"boilerplate/internal/domain"
)

type AddressRequest struct {
	Street  string `json:"street" validate:"required"`
	UserID  uint64 `json:"user_id" validate:"required"`
	Title   string `json:"title" validate:"required"`
	City    string `json:"city" validate:"required"`
	Country string `json:"country" validate:"required"`
	Address string `json:"address" validate:"required"`
	Lat     string `json:"lat" validate:"required"`
	Lon     string `json:"lon" validate:"required"`
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
