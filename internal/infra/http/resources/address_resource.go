package resources

import (
	"boilerplate/internal/domain"
)

type AddressDto struct {
	Id         uint64  `json:"id"`
	UserId     uint64  `json:"user_id"`
	City       string  `json:"name"`
	Coutry     string  `json:"country"`
	Address    string  `json:"address"`
	Department string  `json:"department"`
	Lat        float64 `json:"lat"`
	Lon        float64 `json:"lon"`
}

func (d AddressDto) DomainToDto(address domain.Address) AddressDto {
	return AddressDto{
		Id:         address.Id,
		UserId:     address.User.Id,
		City:       address.City,
		Coutry:     address.Country,
		Address:    address.Address,
		Department: address.Department,
		Lat:        address.Lat,
		Lon:        address.Lon,
	}
}

func (d AddressDto) DomainToDtoCollection(addresses []domain.Address) []AddressDto {
	result := make([]AddressDto, len(addresses))

	for i := range addresses {
		result[i] = d.DomainToDto(addresses[i])
	}

	return result
}
