package controllers

import (
	"boilerplate/internal/app"
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/http/requests"
	"boilerplate/internal/infra/http/resources"
	"log"
	"net/http"
)

type OrderItemController struct {
	orderItemService  app.OrderItemsService
	imageModelService app.ImageModelService
}

func NewOrderItemController(os app.OrderItemsService, ims app.ImageModelService) OrderItemController {
	return OrderItemController{
		orderItemService:  os,
		imageModelService: ims,
	}
}

func (c OrderItemController) AddItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		o := r.Context().Value(OrderKey).(domain.Order)
		orderI, err := requests.Bind(r, requests.OrderItemRequest{}, domain.OrderItem{})
		if err != nil {
			log.Printf("OrderController: %s", err)
			BadRequest(w, err)
			return
		}

		orderI, err = c.orderItemService.Save(orderI, o.Id)
		if err != nil {
			log.Printf("OrderController: %s", err)
			BadRequest(w, err)
			return
		}

		Created(w, resources.OrderItemDto{}.DomainToDto(orderI, c.imageModelService))
	}
}

func (c OrderItemController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		o := r.Context().Value(OrderItemKey).(domain.OrderItem)
		order, err := requests.Bind(r, requests.OrderItemUpdateRequest{}, domain.OrderItem{})
		if err != nil {
			log.Printf("OfferController: %s", err)
			InternalServerError(w, err)
			return
		}

		newOrder, err := c.orderItemService.Update(o, order)
		if err != nil {
			log.Printf("OfferController: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.OrderItemDto{}.DomainToDto(newOrder, c.imageModelService))
	}
}

func (c OrderItemController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		o := r.Context().Value(OrderItemKey).(domain.OrderItem)
		err := c.orderItemService.Delete(o)
		if err != nil {
			log.Printf("OfferController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}
