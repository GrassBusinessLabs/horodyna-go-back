package resources

import (
	"boilerplate/internal/domain"
)

type CategoryDto struct {
	Name string `json:"name"`
}

type CategoriesDto struct {
	Data []CategoryDto `json:"data"`
}

func (d CategoryDto) DomainToDto(cat domain.Category) CategoryDto {
	return CategoryDto{
		Name: string(cat),
	}
}

func (d CategoryDto) DomainToDtoCollection(cats []domain.Category) CategoriesDto {
	result := make([]CategoryDto, len(cats))

	for i := range cats {
		result[i] = CategoryDto{}.DomainToDto(cats[i])
	}

	return CategoriesDto{Data: result}
}
