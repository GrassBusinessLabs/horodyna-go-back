package app

import (
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/database"
	"log"
)

type AddressService interface {
	Find(uint64) (interface{}, error)
	Create(address domain.Address) (domain.Address, error)
	Read(id uint64) (domain.Address, error)
	Update(address domain.Address) (domain.Address, error)
	Delete(id uint64) error
	FindAll(domain.Pagination) (domain.Addresses, error)
}

type AddressServiceImpl struct {
	addressRepo database.AddressRepository
}

func NewAddressService(addressRepo database.AddressRepository) AddressService {
	return &AddressServiceImpl{
		addressRepo: addressRepo,
	}
}

// NewAddressService створює новий екземпляр AddressService.

// Create створює нову адресу.
func (as AddressServiceImpl) Create(address domain.Address) (domain.Address, error) {
	newAddress, err := as.addressRepo.Create(address)
	if err != nil {
		log.Printf("AddressService: %s", err)
		return domain.Address{}, err
	}
	return newAddress, nil
}

// Read повертає інформацію про конкретну адресу за ідентифікатором.
func (as AddressServiceImpl) Read(id uint64) (domain.Address, error) {
	address, err := as.addressRepo.Read(id)
	if err != nil {
		log.Printf("AddressService: %s", err)
		return domain.Address{}, err
	}
	return address, nil
}

// Update оновлює інформацію про адресу.
func (as AddressServiceImpl) Update(address domain.Address) (domain.Address, error) {
	updatedAddress, err := as.addressRepo.Update(address)
	if err != nil {
		log.Printf("AddressService: %s", err)
		return domain.Address{}, err
	}
	return updatedAddress, nil
}

// Delete видаляє адресу за ідентифікатором.
func (as AddressServiceImpl) Delete(id uint64) error {
	err := as.addressRepo.Delete(id)
	if err != nil {
		log.Printf("AddressService: %s", err)
		return err
	}
	return nil
}

func (as AddressServiceImpl) FindAll(p domain.Pagination) (domain.Addresses, error) {
	addresses, err := as.addressRepo.FindAll(p)
	if err != nil {
		log.Printf("AddressService: %s", err)
		return domain.Addresses{}, err
	}

	return addresses, nil
}

func (s AddressServiceImpl) Find(id uint64) (interface{}, error) {
	f, err := s.addressRepo.Read(id)
	if err != nil {
		log.Printf("FarmService -> Find: %s", err)
		return domain.Farm{}, err
	}
	return f, err
}
