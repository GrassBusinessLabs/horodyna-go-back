package controllers

import (
	"boilerplate/internal/app"
	"boilerplate/internal/infra/http/resources"
	"net/http"
)

type CategoryController struct {
	catService app.CategoryService
}

func NewCategoryController(us app.CategoryService) CategoryController {
	return CategoryController{
		catService: us,
	}
}

func (c CategoryController) FindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Success(w, resources.CategoryDto{}.DomainToDtoCollection(c.catService.FindAll()))
	}
}
