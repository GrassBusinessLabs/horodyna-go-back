package app

import (
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/database"
	"errors"
	"log"
)

type AddressService interface {
	Save(address domain.Address) (domain.Address, error)
	Find(uint64) (interface{}, error)
	FindByUserId(userId uint64) (domain.Address, error)
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
	_, err := s.addressRepo.FindByUserId(address.User.Id)
	if err == nil {
		err = errors.New("there is already address with such user id")
		log.Printf("AddressService: %s", err)
		return domain.Address{}, err
	}

	address, err = s.addressRepo.Save(address)
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

func (s addressService) FindByUserId(userId uint64) (domain.Address, error) {
	address, err := s.addressRepo.FindByUserId(userId)
	if err != nil {
		log.Printf("AddressService: %s", err)
		return domain.Address{}, err
	}

	return address, nil
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
