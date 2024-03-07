package database

import (
	"boilerplate/internal/domain"
	"time"

	"github.com/upper/db/v4"
)

const AddressesTableName = "addresses"

type address struct {
	Id          uint64     `db:"id,omitempty"`
	UserId      uint64     `db:"user_id,omitempty"`
	City        string     `db:"city"`
	Country     string     `db:"country"`
	Address     string     `db:"address"`
	Department  string     `db:"department"`
	Lat         float64    `db:"lat"`
	Lon         float64    `db:"lon"`
	CityRef     *string    `db:"city_ref"`
	CreatedDate time.Time  `db:"created_date,omitempty"`
	UpdatedDate time.Time  `db:"updated_date,omitempty"`
	DeletedDate *time.Time `db:"deleted_date,omitempty"`
}

type AddressRepository interface {
	Save(address domain.Address) (domain.Address, error)
	FindById(id uint64) (domain.Address, error)
	FindByUserId(userId uint64) (domain.Address, error)
	Update(address domain.Address) (domain.Address, error)
	Delete(id uint64) error
}

type addressRepository struct {
	coll db.Collection
}

func NewAddressRepository(dbSession db.Session) AddressRepository {
	return addressRepository{
		coll: dbSession.Collection(AddressesTableName),
	}
}

func (r addressRepository) Save(address domain.Address) (domain.Address, error) {
	addressModel := r.mapDomainToModel(address)
	addressModel.CreatedDate, addressModel.UpdatedDate = time.Now(), time.Now()
	err := r.coll.InsertReturning(&addressModel)
	if err != nil {
		return domain.Address{}, err
	}

	var userModel user
	err = r.coll.Session().SQL().Select("*").From("users").Where(db.Cond{"id": addressModel.UserId}).One(&userModel)
	if err != nil {
		return domain.Address{}, err
	}

	address = r.mapModelToDomain(addressModel, userModel)
	return address, nil
}

func (r addressRepository) Update(address domain.Address) (domain.Address, error) {
	addressModel := r.mapDomainToModel(address)
	address.UpdatedDate = time.Now()
	err := r.coll.Find(db.Cond{"id": addressModel.Id}).Update(&addressModel)
	if err != nil {
		return domain.Address{}, err
	}

	var userModel user
	err = r.coll.Session().SQL().Select("*").From("users").Where(db.Cond{"id": addressModel.UserId}).One(&userModel)
	if err != nil {
		return domain.Address{}, err
	}

	address = r.mapModelToDomain(addressModel, userModel)
	return address, nil
}

func (r addressRepository) FindById(id uint64) (domain.Address, error) {
	var addressModel address
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&addressModel)
	if err != nil {
		return domain.Address{}, err
	}

	var userModel user
	err = r.coll.Session().SQL().Select("*").From("users").Where(db.Cond{"id": addressModel.UserId}).One(&userModel)
	if err != nil {
		return domain.Address{}, err
	}

	address := r.mapModelToDomain(addressModel, userModel)
	return address, nil
}

func (r addressRepository) FindByUserId(userId uint64) (domain.Address, error) {
	var addressModel address
	err := r.coll.Find(db.Cond{"user_id": userId, "deleted_date": nil}).One(&addressModel)
	if err != nil {
		return domain.Address{}, err
	}

	var userModel user
	err = r.coll.Session().SQL().Select("*").From("users").Where(db.Cond{"id": userId}).One(&userModel)
	if err != nil {
		return domain.Address{}, err
	}

	return r.mapModelToDomain(addressModel, userModel), nil
}

func (r addressRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r addressRepository) mapDomainToModel(d domain.Address) address {
	return address{
		Id:          d.Id,
		UserId:      d.User.Id,
		City:        d.City,
		Address:     d.Address,
		Department:  d.Department,
		Lon:         d.Lon,
		Lat:         d.Lat,
		CityRef:     d.CityRef,
		CreatedDate: d.CreatedDate,
		UpdatedDate: d.UpdatedDate,
		DeletedDate: d.DeletedDate,
	}
}

func (r addressRepository) mapModelToDomain(m address, userModel user) domain.Address {
	return domain.Address{
		Id:          m.Id,
		User:        mapModelToDomainUser(userModel),
		City:        m.City,
		Address:     m.Address,
		Department:  m.Department,
		Lon:         m.Lon,
		Lat:         m.Lat,
		CityRef:     m.CityRef,
		CreatedDate: m.CreatedDate,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}
}
