package resources

import (
	"boilerplate/internal/app"
	"boilerplate/internal/domain"
	"log"
)

type AddressDto struct {
	Street  string  `json:"street"`
	City    string  `json:"city"`
	State   string  `json:"state"`
	ZipCode string  `json:"zip_code"`
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
		Street:  address.Street,
		City:    address.City,
		State:   address.State,
		ZipCode: address.ZipCode,
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
