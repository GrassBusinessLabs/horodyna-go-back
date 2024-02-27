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
	FindByFarmUserId(farmUserId uint64, p domain.Pagination) (domain.Orders, error)
	SetOrderStatus(order domain.Order) (domain.Order, error)
}

func NewOrderService(or database.OrderRepository) OrderService {
	return orderService{
		orderRepo: or,
	}
}

type orderService struct {
	orderRepo database.OrderRepository
}

func (s orderService) Find(id uint64) (interface{}, error) {
	o, err := s.orderRepo.FindById(id)
	if err != nil {
		log.Printf("OrderService -> Find: %s", err)
		return domain.Order{}, err
	}

	return o, err
}

func (s orderService) Save(ord domain.Order) (domain.Order, error) {
	ord.Status = domain.DRAFT
	o, err := s.orderRepo.Save(ord)
	if err != nil {
		log.Printf("OrderService: %s", err)
		return domain.Order{}, err
	}

	return o, err
}

func (s orderService) FindById(id uint64) (domain.Order, error) {
	order, err := s.orderRepo.FindById(id)

	if err != nil {
		log.Printf("OrderService: %s", err)
		return domain.Order{}, err
	}

	return order, err
}

func (s orderService) FindAllByUserId(userId uint64, pag domain.Pagination) (domain.Orders, error) {
	orders, err := s.orderRepo.FindAllByUserId(userId, pag)
	if err != nil {
		log.Printf("OrderService: %s", err)
		return domain.Orders{}, err
	}

	return orders, nil
}

func (s orderService) Update(ord domain.Order, req domain.Order) (domain.Order, error) {
	ord.Address = req.Address
	ord.Comment = req.Comment
	if ord.ShippingPrice != req.ShippingPrice {
		ord.TotalPrice = req.ShippingPrice + ord.ProductsPrice
	}

	ord.ShippingPrice = req.ShippingPrice
	order, err := s.orderRepo.Update(ord)
	if err != nil {
		log.Printf("OrderService: %s", err)
		return domain.Order{}, err
	}

	return order, nil
}

func (s orderService) Delete(order domain.Order) error {
	err := s.orderRepo.Delete(order)
	if err != nil {
		log.Printf("OrderService: %s", err)
		return err
	}

	return nil
}

func (s orderService) FindByFarmUserId(farmUserId uint64, p domain.Pagination) (domain.Orders, error) {
	orders, err := s.orderRepo.FindByFarmUserId(farmUserId, p)
	if err != nil {
		log.Printf("OrderService: %s", err)
		return domain.Orders{}, err
	}

	return orders, nil
}

func (s orderService) SetOrderStatus(order domain.Order) (domain.Order, error) {
	order, err := s.orderRepo.SetOrderStatus(order)
	if err != nil {
		log.Printf("OrderService: %s", err)
		return domain.Order{}, err
	}

	return order, nil
}
