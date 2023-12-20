package database

import (
	"boilerplate/internal/domain"
	"time"

	"github.com/upper/db/v4"
)

const ImageTableName = "images"

type imageModel struct {
	Id          uint64     `db:"id,omitempty"`
	Title       string     `db:"title"`
	Entity      string     `db:"entity"`
	EntityId    uint64     `db:"entity_id"`
	CreatedDate time.Time  `db:"created_date,omitempty"`
	UpdatedDate time.Time  `db:"updated_date,omitempty"`
	DeletedDate *time.Time `db:"deleted_date,omitempty"`
}

type ImageModelRepository interface {
	Save(user domain.ImageModel) (domain.ImageModel, error)
	FindById(id uint64) (domain.ImageModel, error)
	Update(user domain.ImageModel) (domain.ImageModel, error)
	Delete(id uint64) error
}

type imageModelRepository struct {
	coll db.Collection
}

func NewImageModelRepository(dbSession db.Session) ImageModelRepository {
	return imageModelRepository{
		coll: dbSession.Collection(ImageTableName),
	}

}

func (r imageModelRepository) Save(imageM domain.ImageModel) (domain.ImageModel, error) {
	i := r.mapDomainToModel(imageM)
	i.CreatedDate, i.UpdatedDate = time.Now(), time.Now()
	err := r.coll.InsertReturning(&i)
	if err != nil {
		return domain.ImageModel{}, err
	}
	return r.mapModelToDomain(i), nil
}

func (r imageModelRepository) FindById(id uint64) (domain.ImageModel, error) {
	var im imageModel
	err := r.coll.Find(db.Cond{"id": id}).One(&im)
	if err != nil {
		return domain.ImageModel{}, err
	}
	return r.mapModelToDomain(im), nil

}

func (r imageModelRepository) Update(imageM domain.ImageModel) (domain.ImageModel, error) {
	i := r.mapDomainToModel(imageM)
	i.UpdatedDate = time.Now()
	err := r.coll.Find(db.Cond{"id": i.Id}).Update(&i)
	if err != nil {
		return domain.ImageModel{}, err
	}
	return r.mapModelToDomain(i), nil
}

func (r imageModelRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r imageModelRepository) mapDomainToModel(d domain.ImageModel) imageModel {
	return imageModel{
		Id:       d.Id,
		Title:    d.Title,
		Entity:   d.Entity,
		EntityId: d.EntityId,
	}
}

func (r imageModelRepository) mapModelToDomain(m imageModel) domain.ImageModel {
	return domain.ImageModel{
		Id:       m.Id,
		Title:    m.Title,
		Entity:   m.Entity,
		EntityId: m.EntityId,
	}
}
