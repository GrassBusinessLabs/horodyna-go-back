package controllers

import (
	"boilerplate/internal/app"
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/http/requests"
	"boilerplate/internal/infra/http/resources"
	"log"
	"net/http"
)

type ImageModelController struct {
	imageModelService app.ImageModelService
}

func NewImageModelController(ir app.ImageModelService) ImageModelController {
	return ImageModelController{
		imageModelService: ir,
	}
}

func (c ImageModelController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decodedBytes, err := requests.Bind(r, requests.RegisterRequest{}, domain.ImageModel{})
		if err != nil {
			log.Printf("ImageModelController: %s", err)
			BadRequest(w, err)
		}

		imageM, err = c.imageModelService.Save(imageM)
		if err != nil {
			log.Printf("ImageModelController: %s", err)
			BadRequest(w, err)
			return
		}

		Created(w, resources.ImageMDto{}.DomainToDto(imageM))
	}
}

func (c ImageModelController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		imageM, err := requests.Bind(r, requests.ImageRequest{}, domain.Image{})
		if err != nil {
			log.Printf("ImageModelController: %s", err)
			BadRequest(w, err)
			return
		}

		i := r.Context().Value(UserKey).(domain.ImageModel)
		imageM, err = c.imageModelService.Update(i, domain.ImageModel{Entity: })
		if err != nil {
			log.Printf("ImageModelController: %s", err)
			InternalServerError(w, err)
			return
		}

		var ImageMDto resources.ImageMDto
		Success(w, ImageMDto.DomainToDto(imageM))
	}
}

func (c ImageModelController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(UserKey).(domain.ImageModel)

		err := c.imageModelService.Delete(u.Id)
		if err != nil {
			log.Printf("ImageModelController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}
