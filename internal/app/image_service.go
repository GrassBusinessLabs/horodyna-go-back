package app

import (
	"boilerplate/internal/domain"
	"boilerplate/internal/filesystem"
	"boilerplate/internal/infra/database"
	"encoding/base64"
	"log"
)

type ImageModelService interface {
	Find(id uint64) (interface{}, error)
	Save(imageM domain.Image) (domain.Image, error)
	FindById(id uint64) (domain.Image, error)
	Update(imageM domain.Image, im domain.Image) (domain.Image, error)
	Delete(id uint64) error
}

type imageModelService struct {
	imageMRepo database.ImageRepository
	imageServ  filesystem.ImageStorageService
}

func NewImageModelService(ir database.ImageRepository, is filesystem.ImageStorageService) ImageModelService {
	return imageModelService{
		imageMRepo: ir,
		imageServ:  is,
	}
}

func (s imageModelService) Save(image domain.Image) (domain.Image, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(image.Data)
	if err != nil {
		log.Printf("ImageModelService: %s", err)
		return domain.Image{}, err
	}

	err = s.imageServ.SaveImage(image.Name, decodedBytes)
	if err != nil {
		log.Printf("ImageModelService: %s", err)
		return domain.Image{}, err
	}

	savedim, err := s.imageMRepo.Save(image)
	if err != nil {
		log.Printf("ImageModelService: %s", err)
		return domain.Image{}, err
	}

	return savedim, nil

}

func (s imageModelService) Find(id uint64) (interface{}, error) {
	i, err := s.imageMRepo.FindById(id)
	if err != nil {
		log.Printf("imageModelService -> Find: %s", err)
		return domain.Image{}, err
	}
	return i, err
}

func (s imageModelService) FindById(id uint64) (domain.Image, error) {
	imageM, err := s.imageMRepo.FindById(id)
	if err != nil {
		log.Printf("ImageModelService: %s", err)
		return domain.Image{}, err
	}

	return imageM, err
}

func (s imageModelService) Update(imageM domain.Image, req domain.Image) (domain.Image, error) {
	imageM, err := s.imageMRepo.Update(imageM)
	if err != nil {
		log.Printf("ImageModelService: %s", err)
		return domain.Image{}, err
	}

	return imageM, nil
}

func (s imageModelService) Delete(id uint64) error {
	err := s.imageMRepo.Delete(id)
	if err != nil {
		log.Printf("ImageModelService: %s", err)
		return err
	}

	return nil
}
