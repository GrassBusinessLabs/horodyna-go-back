package controllers

import (
	"boilerplate/internal/app"
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/http/requests"
	"boilerplate/internal/infra/http/resources"
	"errors"
	"log"
	"net/http"
	"strconv"
)

type ImageModelController struct {
	imageModelService app.ImageModelService
}

func NewImageModelController(ir app.ImageModelService) ImageModelController {
	return ImageModelController{
		imageModelService: ir,
	}
}

func (c ImageModelController) FindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entityStr := r.URL.Query().Get("entity")
		if entityStr == "" {
			BadRequest(w, errors.New("Parameter entity is required!"))
			return
		}
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			BadRequest(w, errors.New("Parameter id is required!"))
			return
		}
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			InternalServerError(w, err)
			return
		}

		images, err := c.imageModelService.FindAll(entityStr, id)
		if err != nil {
			InternalServerError(w, err)
			return
		}

		Created(w, resources.ImageMDto{}.DomainToDtoMass(images))
	}
}

func (c ImageModelController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		imageM, err := requests.Bind(r, requests.ImageRequest{}, domain.Image{})
		if err != nil {
			log.Printf("ImageModelController: 1 %s", err)
			BadRequest(w, err)
			return
		}

		imageM, err = c.imageModelService.Save(imageM)
		if err != nil {
			log.Printf("ImageModelController: 2 %s", err)
			BadRequest(w, err)
			return
		}

		Created(w, resources.ImageMDto{}.DomainToDto(imageM))
	}
}

func (c ImageModelController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		image := r.Context().Value(ImageKey).(domain.Image)
		Success(w, resources.ImageMDto{}.DomainToDto(image))
	}
}

func (c ImageModelController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(ImageKey).(domain.Image)

		err := c.imageModelService.Delete(u.Id)
		if err != nil {
			log.Printf("ImageModelController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}
