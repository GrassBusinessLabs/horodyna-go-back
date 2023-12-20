package controllers

import (
	"boilerplate/internal/app"
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/http/requests"
	"boilerplate/internal/infra/http/resources"
	"log"
	"net/http"
)

type OrderController struct {
	orderService app.OrderService
}

func NewOrderController(os app.OrderService) OrderController {
	return OrderController{
		orderService: os,
	}
}

func (c OrderController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(UserKey).(domain.User)
		order, err := requests.Bind(r, requests.OrderRequest{}, domain.Order{})

		if err != nil {
			log.Printf("OrderController: %s", err)
			BadRequest(w, err)
			return
		}

		order.UserId = u.Id
		order.Status = true

		order, err = c.orderService.Save(order)
		if err != nil {
			log.Printf("OrderController: %s", err)
			BadRequest(w, err)
			return
		}

		Created(w, resources.OrderDto{}.DomainToDto(order))
	}
}

func (c OrderController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		o := r.Context().Value(OrderKey).(domain.Order)
		Success(w, resources.OrderDto{}.DomainToDto(o))
	}
}

func (c OrderController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		o := r.Context().Value(OrderKey).(domain.Order)
		order, err := requests.Bind(r, requests.UpdateOrderRequest{}, domain.Order{})

		if err != nil {
			log.Printf("OfferController: %s", err)
			InternalServerError(w, err)
			return
		}

		newOrder, err := c.orderService.Update(o, order)

		if err != nil {
			log.Printf("OfferController: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.OrderDto{}.DomainToDto(newOrder))
	}
}

func (c OrderController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		o := r.Context().Value(OrderKey).(domain.Order)
		err := c.orderService.Delete(o)

		if err != nil {
			log.Printf("OfferController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}
