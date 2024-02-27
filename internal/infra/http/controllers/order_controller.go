package controllers

import (
	"boilerplate/internal/app"
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/http/requests"
	"boilerplate/internal/infra/http/resources"
	"errors"
	"log"
	"net/http"
)

type OrderController struct {
	orderService     app.OrderService
	orderItemService app.OrderItemsService
}

func NewOrderController(os app.OrderService, ordItemServ app.OrderItemsService) OrderController {
	return OrderController{
		orderService:     os,
		orderItemService: ordItemServ,
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
		orderItems, err := c.orderItemService.FindAll(o.Id)
		if err != nil {
			log.Printf("OrderController: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.OrderDtoWithOrderItems{}.DomainToDto(o, orderItems, resources.ImageMDto{}))
	}
}

func (c OrderController) FindAllByUserId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(UserKey).(domain.User)
		pagination, err := requests.DecodePaginationQuery(r)
		if err != nil {
			log.Printf("OrderController: %s", err)
			InternalServerError(w, err)
			return
		}

		orders, err := c.orderService.FindAllByUserId(u.Id, pagination)
		if err != nil {
			log.Printf("OrderController: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.OrderDto{}.DomainToDtoPaginatedCollection(orders))
	}
}

func (c OrderController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		o := r.Context().Value(OrderKey).(domain.Order)
		order, err := requests.Bind(r, requests.UpdateOrderRequest{}, domain.Order{})
		if err != nil {
			log.Printf("OrderController: %s", err)
			InternalServerError(w, err)
			return
		}

		newOrder, err := c.orderService.Update(o, order)
		if err != nil {
			log.Printf("OrderController: %s", err)
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
			log.Printf("OrderController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}

func (c OrderController) FindByFarmUserId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(UserKey).(domain.User)
		pagination, err := requests.DecodePaginationQuery(r)
		if err != nil {
			log.Printf("OrderController: %s", err)
			InternalServerError(w, err)
			return
		}

		orders, err := c.orderService.FindByFarmUserId(u.Id, pagination)
		if err != nil {
			log.Printf("OrderController: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.OrderDto{}.DomainToDtoPaginatedCollection(orders))
	}
}

func (c OrderController) SetOrderStatusAsReceiver() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderStatus, err := requests.Bind(r, requests.OrderStatusRequest{}, domain.Order{})
		if err != nil {
			log.Printf("OrderController: %s", err)
			BadRequest(w, err)
			return
		}

		orderInstance := r.Context().Value(OrderKey).(domain.Order)
		if orderInstance.IsReceiverStatus(orderStatus.Status) {
			orderInstance.Status = orderStatus.Status
			order, err := c.orderService.NoRequestUpdate(orderInstance)
			if err != nil {
				log.Printf("OrderController: %s", err)
				InternalServerError(w, err)
				return
			}

			Success(w, resources.OrderDto{}.DomainToDto(order))
		} else {
			err = errors.New("status declined")
			log.Printf("OrderController: %s", err)
			BadRequest(w, err)
			return
		}

	}
}

func (c OrderController) SetOrderStatusAsFarmer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderStatus, err := requests.Bind(r, requests.OrderStatusRequest{}, domain.Order{})
		if err != nil {
			log.Printf("OrderController: %s", err)
			BadRequest(w, err)
			return
		}

		orderInstance := r.Context().Value(OrderKey).(domain.Order)
		if orderInstance.IsFarmerStatus(orderStatus.Status) {
			orderInstance.Status = orderStatus.Status
			order, err := c.orderService.NoRequestUpdate(orderInstance)
			if err != nil {
				log.Printf("OrderController: %s", err)
				InternalServerError(w, err)
				return
			}

			Success(w, resources.OrderDto{}.DomainToDto(order))
		} else {
			err = errors.New("status declined")
			log.Printf("OrderController: %s", err)
			BadRequest(w, err)
			return
		}
	}
}
