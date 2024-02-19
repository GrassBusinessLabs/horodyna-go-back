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
}

func NewAddressController(fr app.AddressService) AddressController {
	return AddressController{
		addresservice: fr,
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

		address.User.Id = u.Id
		address, err = c.addresservice.Create(address)
		if err != nil {
			log.Printf("AddressController: %s", err)
			BadRequest(w, err)
			return
		}

		Created(w, resources.AddressDto{}.DomainToDto(address))
	}
}

func (c AddressController) Read() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f := r.Context().Value(AddressKey).(domain.Address)
		Success(w, resources.AddressDto{}.DomainToDto(f))
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
		pagination, err := requests.DecodePaginationQuery(r)
		if err != nil {
			log.Printf("AddressController: %s", err)
			InternalServerError(w, err)
			return
		}

		address, err := c.addresservice.FindAll(pagination)
		if err != nil {
			log.Printf("AddressController: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.AddressDto{}.DomainToDtoPaginatedCollection(address, pagination))
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
		address.User.Id = f.User.Id
		newAddress, err := c.addresservice.Update(address)
		if err != nil {
			log.Printf("AddressController: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.AddressDto{}.DomainToDto(newAddress))
	}
}
