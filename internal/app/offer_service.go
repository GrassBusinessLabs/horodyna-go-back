package app

import (
	"boilerplate/internal/domain"
	"boilerplate/internal/filesystem"
	"boilerplate/internal/infra/database"
	"encoding/base64"
	"log"
)

type OfferService interface {
	Save(offer domain.Offer) (domain.Offer, error)
	FindById(id uint64) (domain.Offer, error)
	Update(off domain.Offer, req domain.Offer) (domain.Offer, error)
	Delete(offer domain.Offer) error
	Find(uint64) (interface{}, error)
	FindAll(user domain.User, p domain.Pagination) (domain.Offers, error)
	FindAllByFarmId(farmId uint64, p domain.Pagination) (domain.Offers, error)
}

func NewOfferService(or database.OfferRepository, fs filesystem.ImageStorageService) OfferService {
	return offerService{
		offerRepo:    or,
		imageService: fs,
	}
}

type offerService struct {
	offerRepo    database.OfferRepository
	imageService filesystem.ImageStorageService
}

func (s offerService) Find(id uint64) (interface{}, error) {
	f, err := s.offerRepo.FindById(id)
	if err != nil {
		log.Printf("OfferService -> Find: %s", err)
		return domain.Offer{}, err
	}
	return f, err
}

func (s offerService) Save(offer domain.Offer) (domain.Offer, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(offer.Cover.Data)
	if err != nil {
		log.Printf("OfferService: %s", err)
		return domain.Offer{}, err
	}

	name, err := s.imageService.SaveImage(offer.Cover.Name, decodedBytes)
	if err != nil {
		log.Printf("OfferService: %s", err)
		return domain.Offer{}, err
	}

	offer.Cover.Name = name
	o, err := s.offerRepo.Save(offer)
	if err != nil {
		log.Printf("OfferService: %s", err)
		return domain.Offer{}, err
	}

	return o, err
}

func (os offerService) FindById(id uint64) (domain.Offer, error) {
	offer, err := os.offerRepo.FindById(id)
	if err != nil {
		log.Printf("OfferService: %s", err)
		return domain.Offer{}, err
	}

	return offer, err
}

func (s offerService) FindAllByFarmId(farmId uint64, p domain.Pagination) (domain.Offers, error) {
	offers, err := s.offerRepo.FindAllByFarmId(farmId, p)
	if err != nil {
		log.Printf("OfferService: %s", err)
		return domain.Offers{}, err
	}

	return offers, nil
}

func (s offerService) Update(off domain.Offer, req domain.Offer) (domain.Offer, error) {
	if req.Cover.Name != "" {
		decodedBytes, err := base64.StdEncoding.DecodeString(req.Cover.Data)
		if err != nil {
			log.Printf("OfferService: %s", err)
			return domain.Offer{}, err
		}

		name, err := s.imageService.UpdateImage(off.Cover.Name, req.Cover.Name, decodedBytes)
		if err != nil {
			log.Printf("OfferService: %s", err)
			return domain.Offer{}, err
		}
		req.Cover.Name = name
	} else {
		req.Cover.Name = off.Cover.Name
	}

	req.Id = off.Id
	req.User.Id = off.User.Id
	offer, err := s.offerRepo.Update(req)
	if err != nil {
		log.Printf("OfferService: %s", err)
		return domain.Offer{}, err
	}

	return offer, nil
}

func (s offerService) Delete(offer domain.Offer) error {
	err := s.imageService.RemoveImage(offer.Cover.Name)
	if err != nil {
		log.Printf("OfferService: %s", err)
	}
	err = s.offerRepo.Delete(offer.Id)
	if err != nil {
		log.Printf("OfferService: %s", err)
		return err
	}

	return nil
}

func (s offerService) FindAll(user domain.User, p domain.Pagination) (domain.Offers, error) {
	offers, err := s.offerRepo.FindAll(user, p)
	if err != nil {
		log.Printf("OfferService: %s", err)
		return domain.Offers{}, err
	}

	return offers, nil
}
