package database

import (
	"boilerplate/internal/domain"
	"math"
	"time"

	"github.com/upper/db/v4"
)

const FarmsTableName = "farms"

type farm struct {
	Id          uint64     `db:"id,omitempty"`
	Name        string     `db:"name"`
	City        string     `db:"city"`
	Address     string     `db:"address"`
	UserId      uint64     `db:"user_id"`
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
	FindAllByCoords(points domain.Points, p domain.Pagination) (domain.Farms, error)
	FindAll(pag domain.Pagination) (domain.Farms, error)
	Delete(id uint64) error
	mapModelToDomain(m farm) domain.Farm
}

type farmRepository struct {
	coll      db.Collection
	offerRepo OfferRepository
}

func NewFarmRepository(dbSession db.Session, offerR OfferRepository) FarmRepository {
	return farmRepository{
		coll:      dbSession.Collection(FarmsTableName),
		offerRepo: offerR,
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
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&f)
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

func (r farmRepository) FindAll(p domain.Pagination) (domain.Farms, error) {
	var data []farm
	query := r.coll.Find(db.Cond{"deleted_date": nil})

	res := query.Paginate(uint(p.CountPerPage))
	err := res.Page(uint(p.Page)).All(&data)
	if err != nil {
		return domain.Farms{}, err
	}

	farms := r.mapModelToDomainPagination(data)

	totalCount, err := res.TotalEntries()
	if err != nil {
		return domain.Farms{}, err
	}

	farms.Total = totalCount
	farms.Pages = uint(math.Ceil(float64(farms.Total) / float64(p.CountPerPage)))

	return farms, nil
}

func (r farmRepository) FindAllByCoords(points domain.Points, p domain.Pagination) (domain.Farms, error) {
	var data []farm
	offers, err := r.offerRepo.FindByCategory(points.Category)
	if err != nil {
		return domain.Farms{}, err
	}

	ids := make([]uint64, len(offers))
	for i, item := range offers {
		ids[i] = item.Farm.Id
	}

	query := r.coll.Find(db.Cond{"deleted_date": nil,
		"id IN":       ids,
		"latitude <":  points.UpperLeftPoint.Lat,
		"latitude >":  points.BottomRightPoint.Lat,
		"longitude <": points.UpperLeftPoint.Lng,
		"longitude >": points.BottomRightPoint.Lng})
	res := query.Paginate(uint(p.CountPerPage))
	err = res.Page(uint(p.Page)).All(&data)
	if err != nil {
		return domain.Farms{}, err
	}

	farms := r.mapModelToDomainPagination(data)
	totalCount, err := res.TotalEntries()
	if err != nil {
		return domain.Farms{}, err
	}

	farms.Total = totalCount
	farms.Pages = uint(math.Ceil(float64(farms.Total) / float64(p.CountPerPage)))
	return farms, nil
}

func (r farmRepository) mapDomainToModel(m domain.Farm) farm {
	return farm{
		Id:          m.Id,
		Name:        m.Name,
		City:        m.City,
		Address:     m.Address,
		CreatedDate: m.CreatedDate,
		UserId:      m.UserId,
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
		UserId:      m.UserId,
		Latitude:    m.Latitude,
		Longitude:   m.Longitude,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}
}

func (f farmRepository) mapModelToDomainPagination(farms []farm) domain.Farms {
	new_farms := make([]domain.Farm, len(farms))
	for i, farm := range farms {
		new_farms[i] = f.mapModelToDomain(farm)
	}
	return domain.Farms{Items: new_farms}
}
