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

type offerWithUser struct {
	Offer           offer
	UserId          uint64 `db:"id_user"`
	UserName        string `db:"user_name"`
	UserEmail       string `db:"user_email"`
	UserPhoneNumber string `db:"user_phone_number"`
}

type OfferRepository interface {
	Save(offer domain.Offer) (domain.Offer, error)
	FindById(id uint64) (domain.Offer, error)
	Update(offer domain.Offer) (domain.Offer, error)
	FindAll(user domain.User, pag domain.Pagination) (domain.Offers, error)
	FindAllByFarmId(farmId uint64, p domain.Pagination) (domain.Offers, error)
	FindOnlyOffersByFarmId(farmId uint64) ([]domain.Offer, error)
	FindByCategory(category string) ([]domain.Offer, error)
	Delete(id uint64) error
}

type offerRepository struct {
	coll     db.Collection
	collUser db.Collection
}

func NewOfferRepository(dbSession db.Session) OfferRepository {
	return offerRepository{
		coll:     dbSession.Collection(OffersTableName),
		collUser: dbSession.Collection(UsersTableName),
	}
}

func (r offerRepository) Save(offer domain.Offer) (domain.Offer, error) {
	u := r.mapDomainToModel(offer)
	u.CreatedDate, u.UpdatedDate = time.Now(), time.Now()
	err := r.coll.InsertReturning(&u)
	if err != nil {
		return domain.Offer{}, err
	}
	user, err := r.GetUserForOffer(u.UserId)
	if err != nil {
		return domain.Offer{}, err
	}
	offer = r.mapModelToDomain(u)
	offer.User = user

	return offer, nil
}

func (r offerRepository) FindById(id uint64) (domain.Offer, error) {
	var o offer
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&o)
	if err != nil {
		return domain.Offer{}, err
	}

	user, err := r.GetUserForOffer(o.UserId)
	if err != nil {
		return domain.Offer{}, err
	}
	offer := r.mapModelToDomain(o)
	offer.User = user

	return offer, nil
}

func (or offerRepository) Update(offer domain.Offer) (domain.Offer, error) {
	o := or.mapDomainToModel(offer)
	o.UpdatedDate = time.Now()
	err := or.coll.Find(db.Cond{"id": o.Id}).Update(&o)
	if err != nil {
		return domain.Offer{}, err
	}
	user, err := or.GetUserForOffer(o.UserId)
	if err != nil {
		return domain.Offer{}, err
	}
	offer = or.mapModelToDomain(o)
	offer.User = user

	return offer, nil
}

func (r offerRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r offerRepository) FindAllByFarmId(farmId uint64, p domain.Pagination) (domain.Offers, error) {
	var data []offerWithUser
	query := r.coll.Session().SQL().Select("ofr.*", "u.id AS id_user", "u.name AS user_name", "u.email AS user_email", "u.phone_number AS user_phone_number").
		From("offers AS ofr").
		Where(" ofr.farm_id = ? AND ofr.deleted_date IS NULL", farmId).
		Join("users AS u").On("u.id = ofr.user_id").
		OrderBy("-ofr.status")
	res := query.Paginate(uint(p.CountPerPage))
	err := res.Page(uint(p.Page)).All(&data)
	if err != nil {
		return domain.Offers{}, err
	}
	offers := r.mapModelsToDomainsWithFarm(data)
	totalCount, err := res.TotalEntries()
	if err != nil {
		return domain.Offers{}, err
	}
	offers.Total = totalCount
	offers.Pages = uint(math.Ceil(float64(offers.Total) / float64(p.CountPerPage)))
	return offers, nil
}

func (r offerRepository) FindOnlyOffersByFarmId(farmId uint64) ([]domain.Offer, error) {
	var data []offer
	err := r.coll.Find(db.Cond{"farm_id": farmId}).All(&data)
	if err != nil {
		return []domain.Offer{}, err
	}

	offers := r.mapModelToDomainMass(data)
	if err != nil {
		return []domain.Offer{}, err
	}

	return offers, nil
}

func (r offerRepository) FindAll(user domain.User, p domain.Pagination) (domain.Offers, error) {
	var data []offerWithUser
	query := r.coll.Session().SQL().Select("ofr.*", "u.id AS id_user", "u.name AS user_name", "u.email AS user_email", "u.phone_number AS user_phone_number").
		From("offers AS ofr")
	if user.Id != 0 {
		query = query.Where("ofr.user_id = ? AND ofr.deleted_date IS NULL", user.Id)
	} else {
		query = query.Where("ofr.deleted_date IS NULL")
	}
	query = query.Join("users AS u").On("u.id = ofr.user_id")
	res := query.Paginate(uint(p.CountPerPage))
	err := res.Page(uint(p.Page)).All(&data)
	if err != nil {
		return domain.Offers{}, err
	}

	offers := r.mapModelsToDomainsWithFarm(data)
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
	query := r.coll.Find(db.Cond{"deleted_date": nil})
	if category != "" {
		query = query.And(db.Cond{"category": category})
	}
	err := query.All(&data)
	if err != nil {
		return []domain.Offer{}, err
	}
	return r.mapModelToDomainMass(data), nil
}

func (r offerRepository) GetUserForOffer(id uint64) (domain.User, error) {
	var user user
	err := r.collUser.Find(db.Cond{"id": id}).Select("id", "name", "email").One(&user)
	if err != nil {
		return domain.User{}, err
	}
	return mapModelToDomainUser(user), nil
}

func (r offerRepository) InsertUsersIntoArray(offers []domain.Offer) error {
	for i := range offers {
		user, err := r.GetUserForOffer(offers[i].User.Id)
		if err != nil {
			return err
		}
		offers[i].User = user
	}
	return nil
}

func mapModelToDomainUser(m user) domain.User {
	return domain.User{
		Id:          m.Id,
		Name:        m.Name,
		Email:       m.Email,
		PhoneNumber: m.PhoneNumber,
	}
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
		UserId:      d.User.Id,
		FarmId:      d.Farm.Id,
		CreatedDate: d.CreatedDate,
		UpdatedDate: d.UpdatedDate,
		DeletedDate: d.DeletedDate,
	}
}

func (r offerRepository) mapModelToDomain(o offer) domain.Offer {
	var additionalImages []image
	additionalImagesQuery := r.coll.Session().SQL().Select("*").From("images").Where("entity = ? AND entity_id = ?", "offers", o.Id)
	err := additionalImagesQuery.All(&additionalImages)
	if err != nil {
		return domain.Offer{}
	}

	return domain.Offer{
		Id:               o.Id,
		Title:            o.Title,
		Description:      o.Description,
		Category:         o.Category,
		Price:            o.Price,
		Unit:             o.Unit,
		Stock:            o.Stock,
		Cover:            domain.Image{Name: o.Cover},
		AdditionalImages: mapImageModelToDomainList(additionalImages),
		Status:           o.Status,
		User:             domain.User{Id: o.UserId},
		Farm:             domain.Farm{Id: o.FarmId},
		CreatedDate:      o.CreatedDate,
		UpdatedDate:      o.UpdatedDate,
		DeletedDate:      o.DeletedDate,
	}
}

func (f offerRepository) mapModelsToDomainsWithFarm(offers []offerWithUser) domain.Offers {
	domainOffers := make([]domain.Offer, len(offers))
	for i, item := range offers {
		domainOffers[i] = f.mapModelToDomain(item.Offer)
		domainOffers[i].User = mapModelToDomainUser(user{Id: item.UserId, Name: item.UserName, Email: item.UserEmail, PhoneNumber: &item.UserPhoneNumber})
	}

	return domain.Offers{Items: domainOffers}
}

func (f offerRepository) mapModelToDomainMass(offers []offer) []domain.Offer {
	newOffers := make([]domain.Offer, len(offers))
	for i, offer := range offers {
		newOffers[i] = f.mapModelToDomain(offer)
	}
	return newOffers
}
