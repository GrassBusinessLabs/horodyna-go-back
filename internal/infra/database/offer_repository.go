package database

import (
	"boilerplate/internal/domain"
	"math"
	"time"

	"github.com/upper/db/v4"
)

const OffersTableName = "offers"

type offer struct {
	Id          uint64     `db:"id,omitempty"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Category    string     `db:"category"`
	Price       float64    `db:"price"`
	Unit        string     `db:"unit"`
	Stock       uint       `db:"stock"`
	Cover       string     `db:"cover"`
	Status      bool       `db:"status"`
	FarmId      uint64     `db:"farm_id"`
	UserId      uint64     `db:"user_id"`
	CreatedDate time.Time  `db:"created_date,omitempty"`
	UpdatedDate time.Time  `db:"updated_date,omitempty"`
	DeletedDate *time.Time `db:"deleted_date,omitempty"`
}

type OfferRepository interface {
	Save(offer domain.Offer) (domain.Offer, error)
	FindById(id uint64) (domain.Offer, error)
	Update(offer domain.Offer) (domain.Offer, error)
	FindAll(user domain.User, pag domain.Pagination) (domain.Offers, error)
	FindAllByFarmId(farmId uint64, p domain.Pagination) (domain.Offers, error)
	FindByCategory(category string) ([]domain.Offer, error)
	Delete(id uint64) error
}

type offerRepository struct {
	coll db.Collection
}

func NewOfferRepository(dbSession db.Session) OfferRepository {
	return offerRepository{
		coll: dbSession.Collection(OffersTableName),
	}
}

func (r offerRepository) Save(offer domain.Offer) (domain.Offer, error) {
	u := r.mapDomainToModel(offer)
	u.CreatedDate, u.UpdatedDate = time.Now(), time.Now()
	err := r.coll.InsertReturning(&u)
	if err != nil {
		return domain.Offer{}, err
	}
	return r.mapModelToDomain(u), nil
}

func (r offerRepository) FindById(id uint64) (domain.Offer, error) {
	var o offer
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&o)
	if err != nil {
		return domain.Offer{}, err
	}
	return r.mapModelToDomain(o), nil
}

func (or offerRepository) Update(offer domain.Offer) (domain.Offer, error) {
	o := or.mapDomainToModel(offer)
	o.UpdatedDate = time.Now()
	err := or.coll.Find(db.Cond{"id": o.Id}).Update(&o)
	if err != nil {
		return domain.Offer{}, err
	}
	return or.mapModelToDomain(o), nil
}

func (r offerRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r offerRepository) FindAllByFarmId(farmId uint64, p domain.Pagination) (domain.Offers, error) {
	var data []offer
	query := r.coll.Find(db.Cond{"status": true, "deleted_date": nil, "farm_id": farmId})
	res := query.Paginate(uint(p.CountPerPage))
	err := res.Page(uint(p.Page)).All(&data)
	if err != nil {
		return domain.Offers{}, err
	}

	offers := r.mapModelToDomainPagination(data)
	totalCount, err := res.TotalEntries()
	if err != nil {
		return domain.Offers{}, err
	}

	offers.Total = totalCount
	offers.Pages = uint(math.Ceil(float64(offers.Total) / float64(p.CountPerPage)))

	return offers, nil
}

func (r offerRepository) FindAll(user domain.User, p domain.Pagination) (domain.Offers, error) {
	var data []offer
	query := r.coll.Find(db.Cond{"deleted_date": nil})
	if user.Id != 0 {
		query = query.And(db.Cond{"user_id": user.Id})
	}

	res := query.Paginate(uint(p.CountPerPage))
	err := res.Page(uint(p.Page)).OrderBy("-status").All(&data)
	if err != nil {
		return domain.Offers{}, err
	}

	offers := r.mapModelToDomainPagination(data)
	totalCount, err := res.TotalEntries()
	if err != nil {
		return domain.Offers{}, err
	}

	offers.Total = totalCount
	offers.Pages = uint(math.Ceil(float64(offers.Total) / float64(p.CountPerPage)))

	return offers, nil
}

func (r offerRepository) FindByCategory(category string) ([]domain.Offer, error) {
	var data []offer
	query := r.coll.Find(db.Cond{"status": true, "deleted_date": nil})
	if category != "" {
		query = query.And(db.Cond{"category": category})
	}
	err := query.All(&data)
	if err != nil {
		return []domain.Offer{}, err
	}
	return r.mapModelToDomainMass(data), nil
}

func (r offerRepository) mapDomainToModel(d domain.Offer) offer {
	return offer{
		Id:          d.Id,
		Title:       d.Title,
		Description: d.Description,
		Category:    d.Category,
		Price:       d.Price,
		Unit:        d.Unit,
		Stock:       d.Stock,
		Cover:       d.Cover.Name,
		Status:      d.Status,
		FarmId:      d.Farm.Id,
		UserId:      d.UserId,
		CreatedDate: d.CreatedDate,
		UpdatedDate: d.UpdatedDate,
		DeletedDate: d.DeletedDate,
	}
}

func (r offerRepository) mapModelToDomain(o offer) domain.Offer {

	return domain.Offer{
		Id:          o.Id,
		Title:       o.Title,
		Description: o.Description,
		Category:    o.Category,
		Price:       o.Price,
		Unit:        o.Unit,
		Stock:       o.Stock,
		Cover:       domain.Image{Name: o.Cover},
		Status:      o.Status,
		Farm:        domain.Farm{Id: o.FarmId},
		UserId:      o.UserId,
		CreatedDate: o.CreatedDate,
		UpdatedDate: o.UpdatedDate,
		DeletedDate: o.DeletedDate,
	}
}

func (f offerRepository) mapModelToDomainMass(offers []offer) []domain.Offer {
	newOffers := make([]domain.Offer, len(offers))
	for i, offer := range offers {
		newOffers[i] = f.mapModelToDomain(offer)
	}

	return newOffers
}

func (f offerRepository) mapModelToDomainPagination(offers []offer) domain.Offers {
	return domain.Offers{Items: f.mapModelToDomainMass(offers)}
}
