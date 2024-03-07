package requests

import (
	"boilerplate/internal/domain"
)

type AddressRequest struct {
	City       string  `json:"city" validate:"required"`
	Country    string  `json:"country" validate:"required"`
	Address    string  `json:"address" validate:"required"`
	Department string  `json:"department" validate:"required"`
	Lat        float64 `json:"lat" validate:"required"`
	Lon        float64 `json:"lon" validate:"required"`
	CityRef    *string `json:"city_ref" validate:"required"`
}

func (r AddressRequest) ToDomainModel() (interface{}, error) {
	return domain.Address{
		City:       r.City,
		Country:    r.Country,
		Address:    r.Address,
		Department: r.Department,
		Lat:        r.Lat,
		Lon:        r.Lon,
		CityRef:    r.CityRef,
	}, nil
}
