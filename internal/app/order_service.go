package app

import (
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/database"
	"log"
)

type OrderService interface {
	Save(o domain.Order) (domain.Order, error)
	FindById(id uint64) (domain.Order, error)
	Update(o domain.Order, req domain.Order) (domain.Order, error)
	FindAllByUserId(userId uint64, p domain.Pagination) (domain.Orders, error)
	Delete(o domain.Order) error
	Find(uint64) (interface{}, error)
}

func NewOrderService(or database.OrderRepository) OrderService {
	return orderService{
		order_repo: or,
	}
}

type orderService struct {
	order_repo database.OrderRepository
}

func (s orderService) Find(id uint64) (interface{}, error) {
	o, err := s.order_repo.FindById(id)
	if err != nil {
		log.Printf("OrderService -> Find: %s", err)
		return domain.Order{}, err
	}

	return o, err
}

func (s orderService) Save(ord domain.Order) (domain.Order, error) {
	o, err := s.order_repo.Save(ord)
	if err != nil {
		log.Printf("OrderService: %s", err)
		return domain.Order{}, err
	}
	return o, err
}

func (s orderService) FindById(id uint64) (domain.Order, error) {
	order, err := s.order_repo.FindById(id)

	if err != nil {
		log.Printf("OrderService: %s", err)
		return domain.Order{}, err
	}

	return order, err
}

func (s orderService) FindAllByUserId(userId uint64, pag domain.Pagination) (domain.Orders, error) {
	orders, err := s.order_repo.FindAllByUserId(userId, pag)
	if err != nil {
		log.Printf("OrderService: %s", err)
		return domain.Orders{}, err
	}

	return orders, nil
}

func (s orderService) Update(ord domain.Order, req domain.Order) (domain.Order, error) {
	ord.AddressId = req.AddressId
	ord.Comment = req.Comment
	if ord.ShippingPrice != req.ShippingPrice {
		ord.TotalPrice = req.ShippingPrice + ord.ProductsPrice
	}

	ord.ShippingPrice = req.ShippingPrice
	ord.Status = req.Status
	order, err := s.order_repo.Update(ord)
	if err != nil {
		log.Printf("OrderService: %s", err)
		return domain.Order{}, err
	}

	return order, nil
}

func (s orderService) Delete(order domain.Order) error {
	err := s.order_repo.Delete(order)
	if err != nil {
		log.Printf("OrderService: %s", err)
		return err
	}

	return nil
}
