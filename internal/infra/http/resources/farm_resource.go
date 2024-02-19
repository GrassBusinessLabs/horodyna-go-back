package resources

import (
	"boilerplate/internal/domain"
)

type FarmDto struct {
	Id        uint64  `json:"id"`
	Name      *string `json:"name"`
	City      string  `json:"city"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	User      UserDto `json:"user"`
}

type FarmWithOutDto struct {
	Id        uint64  `json:"id"`
	Name      *string `json:"name"`
	City      string  `json:"city"`
	Address   string  `json:"address"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	UserId    uint64  `json:"user_id"`
}

type FarmsDto struct {
	Items []FarmDto `json:"items"`
	Pages uint      `json:"pages"`
	Total uint64    `json:"total"`
}

func (d FarmDto) DomainToDto(farm domain.Farm) FarmDto {
	return FarmDto{
		Id:        farm.Id,
		Name:      farm.Name,
		City:      farm.City,
		Address:   farm.Address,
		Latitude:  farm.Latitude,
		Longitude: farm.Longitude,
		User:      UserDto{}.DomainToDto(farm.User),
	}
}

func (d FarmWithOutDto) DomainToDto(farm domain.Farm) FarmWithOutDto {

	return FarmWithOutDto{
		Id:        farm.Id,
		Name:      farm.Name,
		City:      farm.City,
		Address:   farm.Address,
		Latitude:  farm.Latitude,
		Longitude: farm.Longitude,
		UserId:    farm.User.Id,
	}
}

func (d FarmDto) DomainToDtoPaginatedCollection(farms domain.Farms) FarmsDto {
	result := make([]FarmDto, len(farms.Items))

	for i := range farms.Items {
		result[i] = d.DomainToDto(farms.Items[i])
	}

	return FarmsDto{Items: result, Pages: farms.Pages, Total: farms.Total}
}
