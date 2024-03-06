package controllers

import (
	"boilerplate/internal/app"
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/http/requests"
	"boilerplate/internal/infra/http/resources"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type AddressController struct {
	addressService app.AddressService
}

func NewAddressController(ar app.AddressService) AddressController {
	return AddressController{
		addressService: ar,
	}
}

func (c AddressController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address, err := requests.Bind(r, requests.AddressRequest{}, domain.Address{})
		if err != nil {
			log.Printf("AddressController: %s", err)
			BadRequest(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)
		address.User.Id = user.Id
		address, err = c.addressService.Save(address)
		if err != nil {
			log.Printf("AddressController: %s", err)
			InternalServerError(w, err)
			return
		}

		Created(w, resources.AddressDto{}.DomainToDto(address))
	}
}

func (c AddressController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address, err := requests.Bind(r, requests.AddressRequest{}, domain.Address{})
		if err != nil {
			log.Printf("AddressController: %s", err)
			BadRequest(w, err)
			return
		}

		addressInstance := r.Context().Value(AddressKey).(domain.Address)
		addressInstance.City = address.City
		addressInstance.Country = address.Country
		addressInstance.Address = address.Address
		addressInstance.Department = address.Department
		addressInstance.Lat = address.Lat
		addressInstance.Lon = address.Lon
		address, err = c.addressService.Update(addressInstance)
		if err != nil {
			log.Printf("AddressController: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.AddressDto{}.DomainToDto(address))
	}
}

func (c AddressController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := r.Context().Value(AddressKey).(domain.Address)
		err := c.addressService.Delete(address.Id)
		if err != nil {
			log.Printf("AddressController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}

func (c AddressController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := r.Context().Value(AddressKey).(domain.Address)
		Success(w, resources.AddressDto{}.DomainToDto(address))
	}
}

func (c AddressController) FindByUserId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.ParseUint(chi.URLParam(r, "userId"), 10, 64)
		if err != nil {
			log.Printf("OrderController: %s", err)
			BadRequest(w, err)
			return
		}

		address, err := c.addressService.FindByUserId(userId)
		if err != nil {
			log.Printf("OrderController: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.AddressDto{}.DomainToDto(address))
	}
}
