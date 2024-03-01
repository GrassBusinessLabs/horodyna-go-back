package app

import (
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/database"
	"errors"
	"log"
)

type FarmService interface {
	Save(farm domain.Farm) (domain.Farm, error)
	FindById(id uint64) (domain.Farm, error)
	Update(farm domain.Farm, req domain.Farm) (domain.Farm, error)
	Delete(id uint64) error
	Find(uint64) (interface{}, error)
	FindAll(p domain.Pagination) (domain.Farms, error)
	FindAllByCoords(points domain.Points, p domain.Pagination) (domain.Farms, error)
}

func NewFarmService(fr database.FarmRepository, or database.OfferRepository, orr database.OrderRepository) FarmService {
	return farmService{
		farmRepo:  fr,
		offerRepo: or,
		orderRepo: orr,
	}
}

type farmService struct {
	farmRepo  database.FarmRepository
	offerRepo database.OfferRepository
	orderRepo database.OrderRepository
}

func (s farmService) Find(id uint64) (interface{}, error) {
	f, err := s.farmRepo.FindById(id)
	if err != nil {
		log.Printf("FarmService -> Find: %s", err)
		return domain.Farm{}, err
	}
	return f, err
}

func (s farmService) Save(farm domain.Farm) (domain.Farm, error) {
	u, err := s.farmRepo.Save(farm)
	if err != nil {
		log.Printf("FarmService: %s", err)
		return domain.Farm{}, err
	}

	return u, err
}

func (s farmService) FindById(id uint64) (domain.Farm, error) {
	farm, err := s.farmRepo.FindById(id)

	if err != nil {
		log.Printf("FarmService: %s", err)
		return domain.Farm{}, err
	}

	return farm, err
}

func (s farmService) Update(farm domain.Farm, req domain.Farm) (domain.Farm, error) {
	farm, err := s.farmRepo.Update(farm)
	if err != nil {
		log.Printf("FarmService: %s", err)
		return domain.Farm{}, err
	}

	return farm, nil
}

func (s farmService) Delete(id uint64) error {
	activeOrders, err := s.orderRepo.GetActiveOrdersByFarmId(id)
	if err != nil {
		log.Printf("FarmService: %s", err)
		return err
	}

	if len(activeOrders) > 0 {
		return errors.New("you can`t delete farm if there is still acitive orders")
	}

	offers, err := s.offerRepo.FindOnlyOffersByFarmId(id)
	if err != nil {
		log.Printf("FarmService: %s", err)
		return err
	}

	for _, offer := range offers {
		err = s.offerRepo.Delete(offer.Id)
		if err != nil {
			log.Printf("FarmService: %s", err)
			return err
		}
	}

	err = s.farmRepo.Delete(id)
	if err != nil {
		log.Printf("FarmService: %s", err)
		return err
	}

	return nil
}

func (s farmService) FindAll(p domain.Pagination) (domain.Farms, error) {
	farms, err := s.farmRepo.FindAll(p)
	if err != nil {
		log.Printf("FarmService: %s", err)
		return domain.Farms{}, err
	}

	return farms, nil
}

func (s farmService) FindAllByCoords(points domain.Points, p domain.Pagination) (domain.Farms, error) {
	farms, err := s.farmRepo.FindAllByCoords(points, p)
	if err != nil {
		log.Printf("FarmService: %s", err)
		return domain.Farms{}, err
	}
	return farms, nil
}
