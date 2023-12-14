package controllers

import (
	"boilerplate/internal/app"
	"boilerplate/internal/domain"
	"boilerplate/internal/filesystem"
	"boilerplate/internal/infra/http/requests"
	"boilerplate/internal/infra/http/resources"
	"log"
	"net/http"
)

type OfferController struct {
	offerService app.OfferService
	farmService  app.FarmService
	userService  app.UserService
	fileService  filesystem.ImageStorageService
}

func NewOfferController(os app.OfferService, fr app.FarmService, us app.UserService, fs filesystem.ImageStorageService) OfferController {
	return OfferController{
		offerService: os,
		farmService:  fr,
		userService:  us,
		fileService:  fs,
	}
}

func (c OfferController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(UserKey).(domain.User)

		offer, err := requests.Bind(r, requests.OfferRequest{}, domain.Offer{})
		offer.UserId = u.Id
		offer.Status = true

		if err != nil {
			log.Printf("OfferController: %s", err)
			BadRequest(w, err)
			return
		}

		farm, err := c.farmService.FindById(offer.FarmId)

		if err != nil || farm.GetUserId() != u.Id {
			log.Printf("OfferController: %s", err)
			BadRequest(w, err)
			return
		}

		offer.FarmId = farm.Id
		offer, err = c.offerService.Save(offer, c.fileService)

		if err != nil {
			log.Printf("OfferController: %s", err)
			BadRequest(w, err)
			return
		}

		Created(w, resources.OfferDto{}.DomainToDto(offer, c.farmService, c.userService))
	}
}

func (c OfferController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		o := r.Context().Value(OfferKey).(domain.Offer)
		Success(w, resources.OfferDto{}.DomainToDto(o, c.farmService, c.userService))
	}
}

func (c OfferController) ListView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pagination, err := requests.DecodePaginationQuery(r)

		if err != nil {
			log.Printf("OfferController: %s", err)
			InternalServerError(w, err)
			return
		}

		offers, err := c.offerService.FindAll(pagination)

		if err != nil {
			log.Printf("OfferController: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.OfferDto{}.DomainToDtoPaginatedCollection(offers, pagination, c.farmService, c.userService))
	}
}

func (c OfferController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		o := r.Context().Value(OfferKey).(domain.Offer)

		offer, err := requests.Bind(r, requests.OfferRequest{}, domain.Offer{})

		if err != nil {
			log.Printf("OfferController: %s", err)
			InternalServerError(w, err)
			return
		}

		new_offer, err := c.offerService.Update(o, offer, c.fileService)

		if err != nil {
			log.Printf("OfferController: %s", err)
			InternalServerError(w, err)
			return
		}
		Success(w, resources.OfferDto{}.DomainToDto(new_offer, c.farmService, c.userService))
	}
}

func (c OfferController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		o := r.Context().Value(OfferKey).(domain.Offer)
		err := c.offerService.Delete(o, c.fileService)

		if err != nil {
			log.Printf("OfferController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}
