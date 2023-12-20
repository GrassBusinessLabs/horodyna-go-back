package app

import (
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/database"
	"log"
)

type OrderItemsService interface {
	Save(o domain.OrderItem, orderId uint64) (domain.OrderItem, error)
	Update(o domain.OrderItem, req domain.OrderItem) (domain.OrderItem, error)
	Delete(o domain.OrderItem) error
	Find(uint64) (interface{}, error)
}

func NewOrderItemsService(or database.OrderItemRepository, order database.OrderRepository) orderItemsService {
	return orderItemsService{
		order_items_repo: or,
		order_repo:       order,
	}
}

type orderItemsService struct {
	order_items_repo database.OrderItemRepository
	order_repo       database.OrderRepository
}

func (s orderItemsService) Find(id uint64) (interface{}, error) {
	o, err := s.order_items_repo.FindById(id)
	if err != nil {
		log.Printf("OrderItemService -> Find: %s", err)
		return domain.Order{}, err
	}

	return o, err
}

func (s orderItemsService) Save(ord domain.OrderItem, orderId uint64) (domain.OrderItem, error) {
	o, err := s.order_items_repo.Save(ord, orderId)
	if err != nil {
		log.Printf("OrderItemService: %s", err)
		return domain.OrderItem{}, err
	}
	err = s.order_repo.Recalculate(orderId)
	return o, err
}

func (s orderItemsService) Update(ord domain.OrderItem, req domain.OrderItem) (domain.OrderItem, error) {
	ord.Amount = req.Amount
	ord.TotalPrice = ord.Price * float64(req.Amount)

	order_item, err := s.order_items_repo.Update(ord)
	if err != nil {
		log.Printf("OrderItemService: %s", err)
		return domain.OrderItem{}, err
	}

	err = s.order_repo.Recalculate(ord.OrderId)
	if err != nil {
		log.Printf("OrderItemService: %s", err)
		return domain.OrderItem{}, err
	}

	return order_item, nil
}

func (s orderItemsService) Delete(order domain.OrderItem) error {
	err := s.order_items_repo.Delete(order)
	if err != nil {
		log.Printf("OrderItemService: %s", err)
		return err
	}

	err = s.order_repo.Recalculate(order.OrderId)
	if err != nil {
		log.Printf("OrderItemService: %s", err)
		return err
	}

	return nil
}
