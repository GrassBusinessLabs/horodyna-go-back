package app

import (
	"boilerplate/internal/domain"
)

type CategoryService interface {
	FindAll() []domain.Category
}

type categoryService struct {
	repositoryCat []domain.Category
}

func NewCategoryService() CategoryService {
	return categoryService{
		repositoryCat: domain.GetCategoriesList(),
	}
}

func (s categoryService) FindAll() []domain.Category {
	return s.repositoryCat
}
