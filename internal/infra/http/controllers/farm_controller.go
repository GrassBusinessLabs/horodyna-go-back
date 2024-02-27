package controllers

import (
	"boilerplate/internal/app"
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/http/requests"
	"boilerplate/internal/infra/http/resources"
	"log"
	"net/http"
)

type FarmController struct {
	farmService app.FarmService
}

func NewFarmController(fr app.FarmService) FarmController {
	return FarmController{
		farmService: fr,
	}
}

func (c FarmController) FindAllByCoords() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pagination, err := requests.DecodePaginationQuery(r)
		if err != nil {
			log.Printf("FarmController: %s", err)
			BadRequest(w, err)
			return
		}

		req, err := requests.Bind(r, requests.PointsRequest{}, domain.Points{})
		if err != nil {
			log.Printf("FarmController: %s", err)
			InternalServerError(w, err)
			return
		}

		farms, err := c.farmService.FindAllByCoords(req, pagination)
		if err != nil {
			log.Printf("FarmController: %s", err)
			InternalServerError(w, err)
			return
		}
		Success(w, resources.FarmDto{}.DomainToDtoPaginatedCollection(farms, resources.ImageMDto{}))
	}
}

func (c FarmController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(UserKey).(domain.User)
		farm, err := requests.Bind(r, requests.FarmRequest{}, domain.Farm{})
		farm.User.Id = u.Id
		if err != nil {
			log.Printf("FarmController req: %s", err)
			BadRequest(w, err)
			return
		}

		farm, err = c.farmService.Save(farm)
		if err != nil {
			log.Printf("FarmController save: %s", err)
			BadRequest(w, err)
			return
		}

		Created(w, resources.FarmDto{}.DomainToDto(farm, resources.ImageMDto{}))
	}
}

func (c FarmController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f := r.Context().Value(FarmKey).(domain.Farm)
		Success(w, resources.FarmDto{}.DomainToDto(f, resources.ImageMDto{}))
	}
}

func (c FarmController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f := r.Context().Value(FarmKey).(domain.Farm)
		farm, err := requests.Bind(r, requests.FarmRequest{}, domain.Farm{})
		if err != nil {
			log.Printf("FarmController: %s", err)
			InternalServerError(w, err)
			return
		}

		farm.Id = f.Id
		farm.User.Id = f.User.Id
		newfarm, err := c.farmService.Update(farm, farm)
		if err != nil {
			log.Printf("FarmController: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.FarmDto{}.DomainToDto(newfarm, resources.ImageMDto{}))
	}
}

func (c FarmController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f := r.Context().Value(FarmKey).(domain.Farm)
		err := c.farmService.Delete(f.Id)
		if err != nil {
			log.Printf("FarmController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}

func (c FarmController) ListView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pagination, err := requests.DecodePaginationQuery(r)
		if err != nil {
			log.Printf("FarmController: %s", err)
			InternalServerError(w, err)
			return
		}

		farms, err := c.farmService.FindAll(pagination)
		if err != nil {
			log.Printf("FarmController: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.FarmDto{}.DomainToDtoPaginatedCollection(farms, resources.ImageMDto{}))
	}
}
