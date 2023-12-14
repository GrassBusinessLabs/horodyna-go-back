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
	userService app.UserService
}

func NewFarmController(fr app.FarmService, us app.UserService) FarmController {
	return FarmController{
		farmService: fr,
		userService: us,
	}
}

func (c FarmController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(UserKey).(domain.User)
		farm, err := requests.Bind(r, requests.FarmRequest{}, domain.Farm{})
		farm.UserId = u.Id

		if err != nil {
			log.Printf("FarmController: %s", err)
			BadRequest(w, err)
<<<<<<< HEAD
			return
=======
>>>>>>> 0155198fc8277ca536fde512340a43011cd41860
		}

		farm, err = c.farmService.Save(farm)

		if err != nil {
			log.Printf("FarmController: %s", err)
			BadRequest(w, err)
			return
		}

		Created(w, resources.FarmDto{}.DomainToDto(farm, c.userService))
	}
}

func (c FarmController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f := r.Context().Value(FarmKey).(domain.Farm)
		Success(w, resources.FarmDto{}.DomainToDto(f, c.userService))
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
		farm.UserId = f.UserId
		newfarm, err := c.farmService.Update(farm, farm)

		if err != nil {
			log.Printf("FarmController: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.FarmDto{}.DomainToDto(newfarm, c.userService))
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

		Success(w, resources.FarmDto{}.DomainToDtoPaginatedCollection(farms, pagination, c.userService))
	}
}
