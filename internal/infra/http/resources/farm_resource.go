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
	Longitude float64 `json:"logitude"`
	User      UserDto `json:"user"`
}

type FarmsDto struct {
	Items []FarmDto `json:"items"`
	Total uint64    `json:"total"`
	Pages uint      `json:"pages"`
}

func (d FarmDto) DomainToDto(farm domain.Farm, us app.UserService) FarmDto {
	user, err := us.FindById(farm.User_id)
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

func (d FarmDto) DomainToDtoCollection(farms domain.Farms, us app.UserService) FarmsDto {
	result := make([]FarmDto, len(farms.Items))

	for i := range farms.Items {
		result[i] = d.DomainToDto(farms.Items[i], us)
	}

	return FarmsDto{Items: result, Pages: farms.Pages, Total: farms.Total}
}
