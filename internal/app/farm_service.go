package app

import (
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/database"
	"log"
)

type FarmService interface {
	Save(farm domain.Farm) (domain.Farm, error)
	FindById(id uint64) (domain.Farm, error)
	Update(farm domain.Farm, req domain.Farm) (domain.Farm, error)
	Delete(id uint64) error
	FindAll() (domain.Farms, error)
	Count() uint64
}

type farmService struct {
	farmRepo database.FarmRepository
}

func (f farmService) Count() uint64 {
	count, err := f.farmRepo.Count()

	if err != nil {
		log.Printf("FarmService: %s", err)
		return 0
	}

	return count
}

func NewFarmService(fr database.FarmRepository) FarmService {
	return farmService{
		farmRepo: fr,
	}
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
	err := s.farmRepo.Delete(id)
	if err != nil {
		log.Printf("FarmService: %s", err)
		return err
	}

	return nil
}

func (s farmService) FindAll() (domain.Farms, error) {
	farms, err := s.farmRepo.FindAll()
	if err != nil {
		log.Printf("FarmService: %s", err)
		return domain.Farms{}, err
	}

	return farms, nil
}
