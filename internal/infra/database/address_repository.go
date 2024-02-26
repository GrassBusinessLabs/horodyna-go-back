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
	Lat         float64   `db:"lat"`
	Lon         float64   `db:"lon"`
	CreatedDate time.Time `db:"created_date,omitempty"`
	UpdatedDate time.Time `db:"updated_date,omitempty"`
	DeletedDate time.Time `db:"deleted_date,omitempty"`
}

type addressWithUser struct {
	Address   address
	UserId    uint64 `db:"id_user"`
	UserName  string `db:"user_name"`
	UserEmail string `db:"user_email"`
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
	address, err = r.findAddressWithUser(a.ID)
	if err != nil {
		return domain.Address{}, nil
	}
	return address, nil
}

func (r addressRepository) Read(id uint64) (domain.Address, error) {
	address, err := r.findAddressWithUser(id)
	if err != nil {
		return domain.Address{}, err
	}
	return address, nil
}

func (r addressRepository) Update(address domain.Address) (domain.Address, error) {
	a := r.mapDomainToModel(address)
	a.UpdatedDate = time.Now()
	err := r.coll.Find(db.Cond{"id": address.ID, "deleted_date": nil}).Update(a)
	if err != nil {
		return domain.Address{}, err
	}
	address, err = r.findAddressWithUser(a.ID)
	if err != nil {
		return domain.Address{}, nil
	}
	return address, nil
}

func (r addressRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r addressRepository) FindAll(p domain.Pagination) (domain.Addresses, error) {
	addresses, err := r.findAddressesWithUsers(db.Cond{"addresses.deleted_date": nil}, p)
	if err != nil {
		return domain.Addresses{}, err
	}

	return addresses, nil
}

func (r addressRepository) findAddressesWithUsers(cond db.Cond, p domain.Pagination) (domain.Addresses, error) {
	var addresses []addressWithUser
	query := r.coll.Session().SQL().Select("addresses.*", "u.id AS id_user", "u.name AS user_name", "u.email AS user_email").
		From("addresses").
		Where(cond).
		Join("users as u").On("u.id = addresses.user_id")
	res := query.Paginate(uint(p.CountPerPage))
	err := res.Page(uint(p.Page)).All(&addresses)
	if err != nil {
		return domain.Addresses{}, err
	}

	domainAddresses := make([]domain.Address, len(addresses))
	for i, address := range addresses {
		domainAddresses[i] = r.mapModelToDomain(address.Address, user{Id: address.UserId, Name: address.UserName, Email: address.UserEmail})
	}
	addressesR := domain.Addresses{Items: domainAddresses}
	totalCount, err := res.TotalEntries()
	if err != nil {
		return domain.Addresses{}, err
	}

	addressesR.Total = totalCount
	addressesR.Pages = uint(math.Ceil(float64(addressesR.Total) / float64(p.CountPerPage)))

	return addressesR, nil
}

func (r addressRepository) findAddressWithUser(AddressId uint64) (domain.Address, error) {
	var address addressWithUser
	err := r.coll.Session().SQL().Select("addresses.*", "u.id AS id_user", "u.name AS user_name", "u.email AS user_email").
		From("addresses").
		Where(db.Cond{"addresses.id": AddressId, "addresses.deleted_date": nil}).
		Join("users as u").On("u.id = addresses.user_id").One(&address)
	if err != nil {
		return domain.Address{}, err
	}

	return r.mapModelToDomain(address.Address, user{Id: address.UserId, Name: address.UserName, Email: address.UserEmail}), nil
}

func (r addressRepository) mapDomainToModel(m domain.Address) address {
	return address{
		ID:          m.ID,
		UserID:      m.User.Id,
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

func (r addressRepository) mapModelToDomain(m address, u user) domain.Address {
	return domain.Address{
		ID:          m.ID,
		User:        mapModelToDomainUser(u),
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

func mapModelToDomainAdress(m address) domain.Address {
	return domain.Address{
		ID:          m.ID,
		User:        domain.User{Id: m.UserID},
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
