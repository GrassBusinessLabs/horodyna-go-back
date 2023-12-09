package database

import (
	"boilerplate/internal/domain"
	"time"

	"github.com/upper/db/v4"
)

const FarmsTableName = "farms"

type farm struct {
	Id          uint64     `db:"id,omitempty"`
	Name        string     `db:"name"`
	City        string     `db:"city"`
	Address     string     `db:"address"`
	User_id     uint64     `db:"user_id"`
	Longitude   float64    `db:"longitude"`
	Latitude    float64    `db:"latitude"`
	CreatedDate time.Time  `db:"created_date,omitempty"`
	UpdatedDate time.Time  `db:"updated_date,omitempty"`
	DeletedDate *time.Time `db:"deleted_date,omitempty"`
}

type FarmRepository interface {
	Save(farm domain.Farm) (domain.Farm, error)
	FindById(id uint64) (domain.Farm, error)
	Update(farm domain.Farm) (domain.Farm, error)
	FindAll() (domain.Farms, error)
	Delete(id uint64) error
	Count() (uint64, error)
}

type farmRepository struct {
	coll db.Collection
}

func NewFarmRepository(dbSession db.Session) FarmRepository {
	return farmRepository{
		coll: dbSession.Collection(FarmsTableName),
	}
}

func (r farmRepository) Save(farm domain.Farm) (domain.Farm, error) {
	u := r.mapDomainToModel(farm)
	u.CreatedDate, u.UpdatedDate = time.Now(), time.Now()
	err := r.coll.InsertReturning(&u)
	if err != nil {
		return domain.Farm{}, err
	}
	return r.mapModelToDomain(u), nil
}

func (r farmRepository) FindById(id uint64) (domain.Farm, error) {
	var f farm
	err := r.coll.Find(db.Cond{"id": id}).One(&f)
	if err != nil {
		return domain.Farm{}, err
	}
	return r.mapModelToDomain(f), nil
}

func (r farmRepository) Update(farm domain.Farm) (domain.Farm, error) {
	u := r.mapDomainToModel(farm)
	u.UpdatedDate = time.Now()
	err := r.coll.Find(db.Cond{"id": u.Id}).Update(&u)
	if err != nil {
		return domain.Farm{}, err
	}
	return r.mapModelToDomain(u), nil
}

func (r farmRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r farmRepository) Count() (uint64, error) {
	return r.coll.Count()
}

func (r farmRepository) FindAll() (domain.Farms, error) {
	var farms []farm
	err := r.coll.Find(db.Cond{}).All(&farms)
	if err != nil {
		return domain.Farms{}, err
	}

	return domain.Farms{Items: r.mapModelListToDomainList(farms)}, nil
}

func (r farmRepository) mapModelListToDomainList(m []farm) []domain.Farm {
	domainList := []domain.Farm{}

	for i := range m {
		domainList = append(domainList, r.mapModelToDomain(m[i]))
	}

	return domainList
}

func (r farmRepository) mapDomainToModel(m domain.Farm) farm {
	return farm{
		Id:          m.Id,
		Name:        m.Name,
		City:        m.City,
		Address:     m.Address,
		CreatedDate: m.CreatedDate,
		User_id:     m.User_id,
		Latitude:    m.Latitude,
		Longitude:   m.Longitude,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}
}

func (r farmRepository) mapModelToDomain(m farm) domain.Farm {
	return domain.Farm{
		Id:          m.Id,
		Name:        m.Name,
		City:        m.City,
		Address:     m.Address,
		CreatedDate: m.CreatedDate,
		User_id:     m.User_id,
		Latitude:    m.Latitude,
		Longitude:   m.Longitude,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}
}
