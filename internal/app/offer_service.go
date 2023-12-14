package app

import (
	"boilerplate/internal/domain"
	"boilerplate/internal/filesystem"
	"boilerplate/internal/infra/database"
	"encoding/base64"
	"log"
)

type OfferService interface {
	Save(offer domain.Offer, fs filesystem.ImageStorageService) (domain.Offer, error)
	FindById(id uint64) (domain.Offer, error)
	Update(off domain.Offer, req domain.Offer, fs filesystem.ImageStorageService) (domain.Offer, error)
	Delete(offer domain.Offer, fs filesystem.ImageStorageService) error
	Find(uint64) (interface{}, error)
	FindAll(p domain.Pagination) (domain.Offers, error)
}

func NewOfferService(or database.OfferRepository) OfferService {
	return offerService{
		offerRepo: or,
	}
}

type offerService struct {
	offerRepo database.OfferRepository
}

func (s offerService) Find(id uint64) (interface{}, error) {
	f, err := s.offerRepo.FindById(id)
	if err != nil {
		log.Printf("OfferService -> Find: %s", err)
		return domain.Offer{}, err
	}
	return f, err
}

func (s offerService) Save(offer domain.Offer, fs filesystem.ImageStorageService) (domain.Offer, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(offer.Image.Data)

	if err != nil {
		log.Printf("OfferService: %s", err)
		return domain.Offer{}, err
	}

	offer.Cover = offer.Image.Name
	err = fs.SaveImage(offer.Cover, decodedBytes)

	if err != nil {
		log.Printf("OfferService: %s", err)
		return domain.Offer{}, err
	}

	offer, uerr := s.offerRepo.Update(offer)
	if uerr != nil {
		log.Printf("OfferService: %s", uerr)
		return domain.Offer{}, uerr
	}

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

func (s offerService) Update(off domain.Offer, req domain.Offer, fs filesystem.ImageStorageService) (domain.Offer, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(req.Image.Data)

	if err != nil {
		log.Printf("OfferService: %s", err)
		return domain.Offer{}, err
	}

	err = fs.RemoveImage(off.Cover)

	if err != nil {
		log.Printf("OfferService: %s", err)
	}

	req.Id = off.Id
	req.UserId = off.UserId

	err = fs.SaveImage(req.Image.Name, decodedBytes)

	if err != nil {
		log.Printf("OfferService: %s", err)
		return domain.Offer{}, err
	}

	req.Cover = req.Image.Name

	offer, err := s.offerRepo.Update(req)
	if err != nil {
		log.Printf("OfferService: %s", err)
		return domain.Offer{}, err
	}

	return offer, nil
}

func (s offerService) Delete(offer domain.Offer, fs filesystem.ImageStorageService) error {
	fs.RemoveImage(offer.Cover)
	err := s.offerRepo.Delete(offer.Id)
	if err != nil {
		log.Printf("OfferService: %s", err)
		return err
	}

	return nil
}

func (s offerService) FindAll(p domain.Pagination) (domain.Offers, error) {
	offers, err := s.offerRepo.FindAll(p)
	if err != nil {
		log.Printf("OfferService: %s", err)
		return domain.Offers{}, err
	}

	return offers, nil
}
