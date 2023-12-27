package resources

import (
	"boilerplate/internal/app"
	"boilerplate/internal/domain"
	"log"
)

type AddressDto struct {
	Id      uint64  `json:"id"`
	Title   string  `json:"title"`
	City    string  `json:"city"`
	Country string  `json:"country"`
	Address string  `json:"address"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	User    UserDto `json:"user"`
}

type AddressesDto struct {
	Items []AddressDto `json:"items"`
	Pages uint         `json:"pages"`
	Total uint64       `json:"total"`
}

func (d AddressDto) DomainToDto(address domain.Address, userService app.UserService) AddressDto {
	user, err := userService.FindById(address.UserID)
	if err != nil {
		log.Println(err)
	}

	return AddressDto{
		Id:      address.ID,
		Title:   address.Title,
		City:    address.City,
		Country: address.Country,
		Address: address.Address,
		Lat:     address.Lat,
		Lon:     address.Lon,
		User:    UserDto{}.DomainToDto(user),
	}
}

func (d AddressDto) DomainToDtoPaginatedCollection(addresses domain.Addresses, pag domain.Pagination, us app.UserService) AddressesDto {
	result := make([]AddressDto, len(addresses.Items))

	for i := range addresses.Items {
		result[i] = d.DomainToDto(addresses.Items[i], us)
	}

	return AddressesDto{Items: result, Pages: addresses.Pages, Total: addresses.Total}
}
