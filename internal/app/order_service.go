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
	NoRequestUpdate(o domain.Order) (domain.Order, error)
	FindAllByUserId(userId uint64, p domain.Pagination) (domain.Orders, error)
	Delete(o domain.Order) error
	Find(uint64) (interface{}, error)
	FindByFarmerId(farmUserId uint64, p domain.Pagination) (domain.Orders, error)
	SplitOrderByFarms(order domain.Order) (map[uint64]domain.Order, error)
	SubmitSplitedOrder(order domain.Order, farmId uint64) (domain.Order, error)
	DeleteSplitedOrder(order domain.Order, farmId uint64) error
	GetFarmerOdersPercentage(farmUserId uint64) (map[uint64]float64, float64, error)
}

func NewOrderService(or database.OrderRepository, oir database.OrderItemRepository) OrderService {
	return orderService{
		orderRepo:     or,
		orderItemRepo: oir,
	}
}

type orderService struct {
	orderRepo     database.OrderRepository
	orderItemRepo database.OrderItemRepository
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
	ord.PostOffice = req.PostOffice
	ord.Ttn = req.Ttn
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

func (s orderService) NoRequestUpdate(order domain.Order) (domain.Order, error) {
	order, err := s.orderRepo.Update(order)
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

func (s orderService) FindByFarmerId(farmUserId uint64, p domain.Pagination) (domain.Orders, error) {
	orders, err := s.orderRepo.GetOrdersByFarmUserId(farmUserId, p)
	if err != nil {
		log.Printf("OrderService: %s", err)
		return domain.Orders{}, err
	}

	return orders, nil
}

func (s orderService) SplitOrderByFarms(order domain.Order) (map[uint64]domain.Order, error) {
	orderItems, err := s.orderItemRepo.FindAllWithoutPagination(order.Id)
	if err != nil {
		log.Printf("OrderService: %s", err)
		return make(map[uint64]domain.Order, 0), err
	}

	order.OrderItems = orderItems
	orders, err := s.orderRepo.SplitOrderByFarms(order)
	if err != nil {
		log.Printf("OrderService: %s", err)
		return make(map[uint64]domain.Order, 0), err
	}

	return orders, nil
}

func (s orderService) SubmitSplitedOrder(order domain.Order, farmId uint64) (domain.Order, error) {
	splitedOrder, err := s.orderRepo.SubmitSplitedOrder(order, farmId)
	if err != nil {
		log.Printf("OrderService: %s", err)
		return domain.Order{}, err
	}

	return splitedOrder, nil
}

func (s orderService) DeleteSplitedOrder(order domain.Order, farmId uint64) error {
	err := s.orderRepo.DeleteSplitedOrder(order, farmId)
	if err != nil {
		log.Printf("OrderService: %s", err)
		return err
	}

	return nil
}

func (s orderService) GetFarmerOdersPercentage(farmUserId uint64) (map[uint64]float64, float64, error) {
	ordersPercentage, total, err := s.orderRepo.GetFarmerOdersPercentage(farmUserId)
	if err != nil {
		log.Printf("OrderService: %s", err)
		return map[uint64]float64{}, 0, err
	}

	return ordersPercentage, total, nil
}
