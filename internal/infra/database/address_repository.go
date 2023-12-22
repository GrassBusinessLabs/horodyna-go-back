package database

import (
	"boilerplate/internal/domain"
	"math"
	"time"

	"github.com/upper/db/v4"
)

const AddressTableName = "addresses"

type address struct {
	ID          uint64    `db:"id,omitempty"`
	UserID      uint64    `db:"user_id"`
	Title       string    `db:"title"`
	City        string    `db:"city"`
	Country     string    `db:"country"`
	Address     string    `db:"address"`
	Lat         string    `db:"lat"`
	Lon         string    `db:"lon"`
	CreatedDate time.Time `db:"created_date,omitempty"`
	UpdatedDate time.Time `db:"updated_date,omitempty"`
	DeletedDate time.Time `db:"deleted_date,omitempty"`
}

type AddressRepository interface {
	Create(address domain.Address) (domain.Address, error)
	Read(id uint64) (domain.Address, error)
	Update(address domain.Address) (domain.Address, error)
	Delete(id uint64) error
	FindAll(domain.Pagination) (domain.Addresses, error)
}

type addressRepository struct {
	coll db.Collection
}

func NewAddressepository(dbSession db.Session) AddressRepository {
	return addressRepository{
		coll: dbSession.Collection(AddressTableName),
	}
}

func (r addressRepository) Create(address domain.Address) (domain.Address, error) {
	a := r.mapDomainToModel(address)
	a.CreatedDate, a.UpdatedDate = time.Now(), time.Now()
	err := r.coll.InsertReturning(&a)
	if err != nil {
		return domain.Address{}, err
	}
	return r.mapModelToDomain(a), nil
}

func (r addressRepository) Read(id uint64) (domain.Address, error) {
	var a address
	err := r.coll.Find(db.Cond{"id": id}).One(&a)
	if err != nil {
		return domain.Address{}, err
	}
	return r.mapModelToDomain(a), nil
}

func (r addressRepository) Update(address domain.Address) (domain.Address, error) {
	a := r.mapDomainToModel(address)
	a.UpdatedDate = time.Now()
	err := r.coll.Find(db.Cond{"id": address.ID}).Update(a)
	if err != nil {
		return domain.Address{}, err
	}
	return r.mapModelToDomain(a), nil
}

func (r addressRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r addressRepository) FindAll(p domain.Pagination) (domain.Addresses, error) {
	var data []address
	query := r.coll.Find(db.Cond{})

	res := query.Paginate(uint(p.CountPerPage))
	err := res.Page(uint(p.Page)).All(&data)
	if err != nil {
		return domain.Addresses{}, err
	}

	addresses := r.mapModelSliceToDomainAddress(data)

	totalCount, err := res.TotalEntries()
	if err != nil {
		return domain.Addresses{}, err
	}

	addresses.Total = totalCount
	addresses.Pages = uint(math.Ceil(float64(addresses.Total) / float64(p.CountPerPage)))

	return addresses, nil
}

func (r addressRepository) mapDomainToModel(m domain.Address) address {
	return address{
		ID:          m.ID,
		UserID:      m.UserID,
		Title:       m.Title,
		City:        m.City,
		Country:     m.Country,
		Address:     m.Address,
		Lat:         m.Lat,
		Lon:         m.Lon,
		CreatedDate: m.CreatedDate,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}
}

func (r addressRepository) mapModelToDomain(m address) domain.Address {
	return domain.Address{
		ID:          m.ID,
		UserID:      m.UserID,
		Title:       m.Title,
		City:        m.City,
		Country:     m.Country,
		Address:     m.Address,
		Lat:         m.Lat,
		Lon:         m.Lon,
		CreatedDate: m.CreatedDate,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}
}

func (r addressRepository) mapModelSliceToDomainAddress(addresses []address) domain.Addresses {
	newAddresses := make([]domain.Address, len(addresses))
	for i, addr := range addresses {
		newAddresses[i] = r.mapModelToDomain(addr)
	}
	return domain.Addresses{Items: newAddresses}
}
