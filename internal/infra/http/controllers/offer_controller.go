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

type OfferController struct {
	offerService app.OfferService
	farmService  app.FarmService
}

func NewOfferController(os app.OfferService, fr app.FarmService) OfferController {
	return OfferController{
		offerService: os,
		farmService:  fr,
	}
}

func (c OfferController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(UserKey).(domain.User)
		offer, err := requests.Bind(r, requests.OfferRequest{}, domain.Offer{})
		if err != nil {
			log.Printf("OfferController: %s", err)
			BadRequest(w, err)
			return
		}
		offer.UserId = u.Id
		offer.Status = true

		farm, err := c.farmService.FindById(offer.Farm.Id)
		if err != nil {
			log.Printf("OfferController: %s", err)
			BadRequest(w, err)
			return
		}
		if farm.GetUserId() != u.Id {
			err := errors.New("User is not a farm owner.")
			log.Printf("OfferController: %s", err)
			BadRequest(w, err)
			return
		}

		offer.Farm = farm
		offer, err = c.offerService.Save(offer)
		if err != nil {
			log.Printf("OfferController: %s", err)
			BadRequest(w, err)
			return
		}

		Created(w, resources.OfferDto{}.DomainToDto(offer))
	}
}

func (c OfferController) FindByFarmId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		farmId, err := strconv.ParseUint(chi.URLParam(r, "farmId"), 10, 64)
		if err != nil {
			log.Printf("OfferController: %s", err)
			InternalServerError(w, err)
			return
		}

		pagination, err := requests.DecodePaginationQuery(r)
		if err != nil {
			log.Printf("OfferController: %s", err)
			InternalServerError(w, err)
			return
		}

		offers, err := c.offerService.FindAllByFarmId(farmId, pagination)
		Success(w, resources.OfferDto{}.DomainToDtoPaginatedCollection(offers))
	}
}

func (c OfferController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		o := r.Context().Value(OfferKey).(domain.Offer)
		Success(w, resources.OfferDto{}.DomainToDto(o))
	}
}

func (c OfferController) ListView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(UserKey).(domain.User)
		pagination, err := requests.DecodePaginationQuery(r)
		if err != nil {
			log.Printf("OfferController: %s", err)
			InternalServerError(w, err)
			return
		}

		offers, err := c.offerService.FindAll(u, pagination)
		if err != nil {
			log.Printf("OfferController: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.OfferDto{}.DomainToDtoPaginatedCollection(offers))
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

		newOffer, err := c.offerService.Update(o, offer)
		if err != nil {
			log.Printf("OfferController: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.OfferDto{}.DomainToDto(newOffer))
	}
}

func (c OfferController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		o := r.Context().Value(OfferKey).(domain.Offer)
		err := c.offerService.Delete(o)
		if err != nil {
			log.Printf("OfferController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}
