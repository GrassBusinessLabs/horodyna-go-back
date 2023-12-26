package controllers

import (
	"boilerplate/internal/app"
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/http/requests"
	"boilerplate/internal/infra/http/resources"
	"log"
	"net/http"
)

type AddressController struct {
	addresservice app.AddressService
	userService   app.UserService
}

func NewAddressController(fr app.AddressService, us app.UserService) AddressController {
	return AddressController{
		addresservice: fr,
		userService:   us,
	}
}

func (c AddressController) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(UserKey).(domain.User)
		address, err := requests.Bind(r, requests.AddressRequest{}, domain.Address{})

		if err != nil {
			log.Printf("AddressController: %s", err)
			BadRequest(w, err)

			return
		}

		address.UserID = u.Id

		address, err = c.addresservice.Create(address)

		if err != nil {
			log.Printf("AddressController: %s", err)
			BadRequest(w, err)
			return
		}

		Created(w, resources.AddressDto{}.DomainToDto(address, c.userService))
	}
}

func (c AddressController) Read() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f := r.Context().Value(AddressKey).(domain.Address)
		Success(w, resources.AddressDto{}.DomainToDto(f, c.userService))
	}
}

func (c AddressController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f := r.Context().Value(AddressKey).(domain.Address)

		err := c.addresservice.Delete(f.ID)

		if err != nil {
			log.Printf("AddressController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}

func (c AddressController) FindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(UserKey).(domain.User)
		pagination, err := requests.DecodePaginationQuery(r)

		if err != nil {
			log.Printf("AddressController: %s", err)
			InternalServerError(w, err)
			return
		}

		address, err := c.addresservice.FindAll(pagination, u)
		if err != nil {
			log.Printf("AddressController: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.AddressDto{}.DomainToDtoPaginatedCollection(address, pagination, c.userService))
	}
}

func (c AddressController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f := r.Context().Value(AddressKey).(domain.Address)
		address, err := requests.Bind(r, requests.AddressRequest{}, domain.Address{})

		if err != nil {
			log.Printf("AddressController: %s", err)
			InternalServerError(w, err)
			return
		}

		address.ID = f.ID
		address.UserID = f.UserID
		newAddress, err := c.addresservice.Update(address)

		if err != nil {
			log.Printf("AddressController: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.AddressDto{}.DomainToDto(newAddress, c.userService))
	}
}
