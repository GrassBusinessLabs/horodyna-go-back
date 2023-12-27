package resources

import (
	"boilerplate/internal/app"
	"boilerplate/internal/domain"
	"log"
)

type FarmDto struct {
	Id        uint64  `json:"id"`
	Name      string  `json:"name"`
	City      string  `json:"city"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	User      UserDto `json:"user"`
}

type FarmsDto struct {
	Items []FarmDto `json:"items"`
	Pages uint      `json:"pages"`
	Total uint64    `json:"total"`
}

func (d FarmDto) DomainToDto(farm domain.Farm, us app.UserService) FarmDto {
	user, err := us.FindById(farm.UserId)
	if err != nil {
		log.Println(err)
	}

	return FarmDto{
		Id:        farm.Id,
		Name:      farm.Name,
		City:      farm.City,
		Address:   farm.Address,
		Latitude:  farm.Latitude,
		Longitude: farm.Longitude,
		User:      UserDto{}.DomainToDto(user),
	}
}

func (d FarmDto) DomainToDtoPaginatedCollection(farms domain.Farms, us app.UserService) FarmsDto {
	result := make([]FarmDto, len(farms.Items))

	for i := range farms.Items {
		result[i] = d.DomainToDto(farms.Items[i], us)
	}

	return FarmsDto{Items: result, Pages: farms.Pages, Total: farms.Total}
}
