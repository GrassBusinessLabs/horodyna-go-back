package resources

import (
	"boilerplate/internal/domain"
)

type FarmDto struct {
	Id        uint64      `json:"id"`
	Name      *string     `json:"name"`
	City      string      `json:"city"`
	Address   string      `json:"address"`
	Latitude  float64     `json:"latitude"`
	Longitude float64     `json:"longitude"`
	AllImages []ImageMDto `json:"all_images"`
	User      UserDto     `json:"user"`
}

type FarmWithOutDto struct {
	Id        uint64      `json:"id"`
	Name      *string     `json:"name"`
	City      string      `json:"city"`
	Address   string      `json:"address"`
	Latitude  float64     `json:"latitude"`
	Longitude float64     `json:"longitude"`
	AllImages []ImageMDto `json:"all_images"`
	UserId    uint64      `json:"user_id"`
}

type FarmsDto struct {
	Items []FarmDto `json:"items"`
	Pages uint      `json:"pages"`
	Total uint64    `json:"total"`
}

func (d FarmDto) DomainToDto(farm domain.Farm, imageDto ImageMDto) FarmDto {
	return FarmDto{
		Id:        farm.Id,
		Name:      farm.Name,
		City:      farm.City,
		Address:   farm.Address,
		Latitude:  farm.Latitude,
		Longitude: farm.Longitude,
		AllImages: imageDto.DomainToDtoMass(farm.AllImages).Items,
		User:      UserDto{}.DomainToDto(farm.User),
	}
}

func (d FarmWithOutDto) DomainToDto(farm domain.Farm, imageDto ImageMDto) FarmWithOutDto {

	return FarmWithOutDto{
		Id:        farm.Id,
		Name:      farm.Name,
		City:      farm.City,
		Address:   farm.Address,
		Latitude:  farm.Latitude,
		Longitude: farm.Longitude,
		AllImages: imageDto.DomainToDtoMass(farm.AllImages).Items,
		UserId:    farm.User.Id,
	}
}

func (d FarmDto) DomainToDtoPaginatedCollection(farms domain.Farms, imageDto ImageMDto) FarmsDto {
	result := make([]FarmDto, len(farms.Items))

	for i := range farms.Items {
		result[i] = d.DomainToDto(farms.Items[i], imageDto)
	}

	return FarmsDto{Items: result, Pages: farms.Pages, Total: farms.Total}
}
