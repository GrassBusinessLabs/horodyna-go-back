package app

import (
	"boilerplate/internal/domain"
	"boilerplate/internal/filesystem"
	"boilerplate/internal/infra/database"
	"encoding/base64"
	"log"
)

type ImageModelService interface {
	Save(imageM domain.ImageModel) (domain.ImageModel, error)
	FindById(id uint64) (domain.ImageModel, error)
	Update(imageM domain.ImageModel, im domain.ImageModel) (domain.ImageModel, error)
	Delete(id uint64) error
}

type imageModelService struct {
	imageMRepo database.ImageModelRepository
	imageServ  filesystem.ImageStorageService
}

func NewImageModelService(ir database.ImageModelRepository, is filesystem.ImageStorageService) ImageModelService {
	return imageModelService{
		imageMRepo: ir,
		imageServ:  is,
	}
}

func (s imageModelService) Save(image domain.Image) (domain.ImageModel, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(image.Data)
	if err != nil {
		log.Printf("ImageModelService: %s", err)
		return domain.ImageModel{}, err
	}

	savedImageM, err := s.imageMRepo.Save(imageM)
	if err != nil {
		log.Printf("ImageModelService: %s", err)
		return domain.ImageModel{}, err
	}

	return savedImageM, err
}

func (s imageModelService) FindById(id uint64) (domain.ImageModel, error) {
	imageM, err := s.imageMRepo.FindById(id)
	if err != nil {
		log.Printf("ImageModelService: %s", err)
		return domain.ImageModel{}, err
	}

	return imageM, err
}

func (s imageModelService) Update(imageM domain.ImageModel, req domain.ImageModel) (domain.ImageModel, error) {
	imageM, err := s.imageMRepo.Update(imageM)
	if err != nil {
		log.Printf("ImageModelService: %s", err)
		return domain.ImageModel{}, err
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
