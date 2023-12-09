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

	"github.com/go-chi/chi/v5"
)

const PER_PAGE = 4

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
		farm.User_id = u.Id

		if err != nil {
			log.Printf("FarmController: %s", err)
			BadRequest(w, err)
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
		idParam := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idParam, 0, 64)

		if err != nil {
			log.Printf("FarmController: %s", err)
			InternalServerError(w, err)
			return
		}

		farm, err := c.farmService.FindById(id)

		if err != nil {
			log.Printf("UserController: %s", err)
			BadRequest(w, err)
			return
		}

		Success(w, resources.FarmDto{}.DomainToDto(farm, c.userService))
	}
}

func (c FarmController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(UserKey).(domain.User)
		farm, err := requests.Bind(r, requests.FarmRequest{}, domain.Farm{})

		if err != nil {
			log.Printf("FarmController: %s", err)
			InternalServerError(w, err)
			return
		}

		idParam := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idParam, 0, 64)

		if err != nil {
			log.Printf("FarmController: %s", err)
			InternalServerError(w, err)
			return
		}

		if u.Id != farm.User_id {
			err = errors.New("Only owner can update!")
			log.Printf("FarmController: %s", err)
			BadRequest(w, err)
			return
		}

		farmold, err := c.farmService.FindById(id)
		farm.Id = farmold.Id
		farm.User_id = farmold.User_id

		if err != nil {
			log.Printf("FarmController: %s", err)
			InternalServerError(w, err)
			return
		}

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
		u := r.Context().Value(UserKey).(domain.User)
		idParam := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idParam, 0, 64)

		if err != nil {
			log.Printf("FarmController: %s", err)
			InternalServerError(w, err)
			return
		}

		farm, err := c.farmService.FindById(id)
		if err != nil {
			log.Printf("FarmController: %s", err)
			InternalServerError(w, err)
			return
		}

		if u.Id != farm.User_id {
			err = errors.New("Only owner can delete!")
			log.Printf("FarmController: %s", err)
			BadRequest(w, err)
			return
		}

		err = c.farmService.Delete(id)

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
		farms, err := c.farmService.FindAll()
		log.Println(farms)
		if err != nil {
			log.Printf("FarmController: %s", err)
			InternalServerError(w, err)
			return
		}

		pageParam := r.URL.Query().Get("page")
		page, err := strconv.Atoi(pageParam)
		if err != nil || page < 1 {
			page = 1
		}

		offset := (page - 1) * PER_PAGE

		startIndex := offset
		endIndex := offset + PER_PAGE
		if endIndex > len(farms.Items) {
			endIndex = len(farms.Items)
		}

		farms.Items = farms.Items[startIndex:endIndex]
		farms.Total = c.farmService.Count()
		farms.Pages = uint(farms.Total) / PER_PAGE

		if farms.Total%PER_PAGE != 0 {
			farms.Pages++
		}

		Success(w, resources.FarmDto{}.DomainToDtoCollection(farms, c.userService))
	}
}
