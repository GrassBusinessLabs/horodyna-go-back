package database

import (
	"boilerplate/internal/domain"
	"fmt"

	"github.com/upper/db/v4"
)

const ImageTableName = "images"

type image struct {
	Id       uint64 `db:"id,omitempty"`
	Name     string `db:"title"`
	Entity   string `db:"entity"`
	EntityId uint64 `db:"entity_id"`
}

type ImageRepository interface {
	FindAll(entity string, id uint64) ([]domain.Image, error)
	Save(user domain.Image) (domain.Image, error)
	FindById(id uint64) (domain.Image, error)
	Delete(id uint64) error
}

type imageRepository struct {
	coll db.Collection
}

func NewImageModelRepository(dbSession db.Session) ImageRepository {
	return imageRepository{
		coll: dbSession.Collection(ImageTableName),
	}
}

func (r imageRepository) Save(imageM domain.Image) (domain.Image, error) {
	i := r.mapDomainToModel(imageM)
	err := r.coll.InsertReturning(&i)
	if err != nil {
		return domain.Image{}, err
	}
	return mapImageModelToDomain(i), nil
}

func (r imageRepository) FindAll(entity string, id uint64) ([]domain.Image, error) {
	var imgs []image
	err := r.coll.Find(db.Cond{"entity": entity, "entity_id": id}).All(&imgs)
	if err != nil {
		return []domain.Image{}, err
	}

	domainImages := make([]domain.Image, len(imgs))
	for i, item := range imgs {
		domainImages[i] = mapImageModelToDomain(item)
	}
	return domainImages, nil
}

func (r imageRepository) FindById(id uint64) (domain.Image, error) {
	var im image
	err := r.coll.Find(db.Cond{"id": id}).One(&im)
	if err != nil {
		return domain.Image{}, err
	}
	return mapImageModelToDomain(im), nil

}

func (r imageRepository) Delete(id uint64) error {
	err := r.coll.Find(db.Cond{"id": id}).Delete()
	if err != nil {
		return fmt.Errorf("error delete: %v", err)
	}
	return nil
}

func (r imageRepository) mapDomainToModel(d domain.Image) image {
	return image{
		Id:       d.Id,
		Name:     d.Name,
		Entity:   d.Entity,
		EntityId: d.EntityId,
	}
}

func mapImageModelToDomain(m image) domain.Image {
	return domain.Image{
		Id:       m.Id,
		Name:     m.Name,
		Entity:   m.Entity,
		EntityId: m.EntityId,
	}
}

func mapImageModelToDomainList(images []image) []domain.Image {
	domainImages := make([]domain.Image, len(images))
	for i, item := range images {
		domainImages[i] = mapImageModelToDomain(item)
	}
	return domainImages
}
