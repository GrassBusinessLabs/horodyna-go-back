package app

import (
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/database"
	"log"
	"math"
)

type OrderItemsService interface {
	Save(o domain.OrderItem, orderId uint64) (domain.OrderItem, error)
	FindAll(orderId uint64) ([]domain.OrderItem, error)
	Update(o domain.OrderItem, req domain.OrderItem) (domain.OrderItem, error)
	Delete(o domain.OrderItem) error
	Find(uint64) (interface{}, error)
}

func NewOrderItemsService(or database.OrderItemRepository, order database.OrderRepository) orderItemsService {
	return orderItemsService{
		orderItemsRepo: or,
		orderRepo:      order,
	}
}

type orderItemsService struct {
	orderItemsRepo database.OrderItemRepository
	orderRepo      database.OrderRepository
}

func (s orderItemsService) FindAll(orderId uint64) ([]domain.OrderItem, error) {
	orderItems, err := s.orderItemsRepo.FindAllWithoutPagination(orderId)
	if err != nil {
		log.Printf("OrderItemService: %s", err)
		return []domain.OrderItem{}, err
	}

	return orderItems, nil
}

func (s orderItemsService) Find(id uint64) (interface{}, error) {
	o, err := s.orderItemsRepo.FindById(id)
	if err != nil {
		log.Printf("OrderItemService -> Find: %s", err)
		return domain.Order{}, err
	}

	return o, err
}

func (s orderItemsService) Save(ord domain.OrderItem, orderId uint64) (domain.OrderItem, error) {
	o, err := s.orderItemsRepo.Save(ord, orderId)
	if err != nil {
		log.Printf("OrderItemService: %s", err)
		return domain.OrderItem{}, err
	}
	err = s.orderRepo.Recalculate(orderId)
	return o, err
}

func (s orderItemsService) Update(ord domain.OrderItem, req domain.OrderItem) (domain.OrderItem, error) {
	ord.Amount = req.Amount
	ord.TotalPrice = math.Round(ord.Price*float64(req.Amount)*100) / 100

	order_item, err := s.orderItemsRepo.Update(ord)
	if err != nil {
		log.Printf("OrderItemService: %s", err)
		return domain.OrderItem{}, err
	}

	err = s.orderRepo.Recalculate(ord.Order.Id)
	if err != nil {
		log.Printf("OrderItemService: %s", err)
		return domain.OrderItem{}, err
	}

	return order_item, nil
}

func (s orderItemsService) Delete(order domain.OrderItem) error {
	err := s.orderItemsRepo.Delete(order.Id)
	if err != nil {
		log.Printf("OrderItemService: %s", err)
		return err
	}

	err = s.orderRepo.Recalculate(order.Order.Id)
	if err != nil {
		log.Printf("OrderItemService: %s", err)
		return err
	}

	return nil
}
