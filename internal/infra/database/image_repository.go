package database

import (
	"boilerplate/internal/domain"
	"time"

	"github.com/upper/db/v4"
)

const ImageTableName = "images"

type image struct {
	Id          uint64     `db:"id,omitempty"`
	Name        string     `db:"title"`
	Entity      string     `db:"entity"`
	EntityId    uint64     `db:"entity_id"`
	CreatedDate time.Time  `db:"created_date,omitempty"`
	UpdatedDate time.Time  `db:"updated_date,omitempty"`
	DeletedDate *time.Time `db:"deleted_date,omitempty"`
}

type ImageRepository interface {
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
	i.CreatedDate, i.UpdatedDate = time.Now(), time.Now()
	err := r.coll.InsertReturning(&i)
	if err != nil {
		return domain.Image{}, err
	}
	return r.mapModelToDomain(i), nil
}

func (r imageRepository) FindById(id uint64) (domain.Image, error) {
	var im image
	err := r.coll.Find(db.Cond{"id": id}).One(&im)
	if err != nil {
		return domain.Image{}, err
	}
	return r.mapModelToDomain(im), nil

}

func (r imageRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r imageRepository) mapDomainToModel(d domain.Image) image {
	return image{
		Id:       d.Id,
		Name:     d.Name,
		Entity:   d.Entity,
		EntityId: d.EntityId,
	}
}

func (r imageRepository) mapModelToDomain(m image) domain.Image {
	return domain.Image{
		Id:       m.Id,
		Name:     m.Name,
		Entity:   m.Entity,
		EntityId: m.EntityId,
	}
}
