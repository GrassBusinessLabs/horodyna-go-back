package requests

import (
	"boilerplate/internal/domain"
)

type AddressRequest struct {
	Street  string `json:"street" validate:"required"`
	City    string `json:"city" validate:"required"`
	State   string `json:"state" validate:"required"`
	ZipCode string `json:"zip_code" validate:"required"`
}

func (m AddressRequest) ToDomainModel() (interface{}, error) {
	return domain.Address{
		Street:  m.Street,
		City:    m.City,
		State:   m.State,
		ZipCode: m.ZipCode,
	}, nil
}
