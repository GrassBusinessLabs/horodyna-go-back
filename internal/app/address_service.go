package app

import (
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/database"
	"log"
)

type AddressService interface {
	Save(address domain.Address) (domain.Address, error)
	Find(uint64) (interface{}, error)
	FindAllByUserId(userId uint64) ([]domain.Address, error)
	Update(address domain.Address) (domain.Address, error)
	Delete(id uint64) error
}

func NewAddressService(ar database.AddressRepository) AddressService {
	return addressService{
		addressRepo: ar,
	}
}

type addressService struct {
	addressRepo database.AddressRepository
}

func (s addressService) Save(address domain.Address) (domain.Address, error) {
	address, err := s.addressRepo.Save(address)
	if err != nil {
		log.Printf("AddressService: %s", err)
		return domain.Address{}, err
	}

	return address, nil
}

func (s addressService) Find(id uint64) (interface{}, error) {
	f, err := s.addressRepo.FindById(id)
	if err != nil {
		log.Printf("AddressService -> Find: %s", err)
		return domain.Farm{}, err
	}
	return f, err
}

func (s addressService) FindAllByUserId(userId uint64) ([]domain.Address, error) {
	addresses, err := s.addressRepo.FindAllByUserId(userId)
	if err != nil {
		log.Printf("AddressService: %s", err)
		return []domain.Address{}, err
	}

	return addresses, nil
}

func (s addressService) Update(address domain.Address) (domain.Address, error) {
	address, err := s.addressRepo.Update(address)
	if err != nil {
		log.Printf("AddressService: %s", err)
		return domain.Address{}, err
	}

	return address, nil
}

func (s addressService) Delete(id uint64) error {
	err := s.addressRepo.Delete(id)
	if err != nil {
		log.Printf("AddressService: %s", err)
		return err
	}

	return nil
}
