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
	Name        *string    `db:"name"`
	City        string     `db:"city"`
	Address     string     `db:"address"`
	UserId      uint64     `db:"user_id"`
	Longitude   float64    `db:"longitude"`
	Latitude    float64    `db:"latitude"`
	CreatedDate time.Time  `db:"created_date,omitempty"`
	UpdatedDate time.Time  `db:"updated_date,omitempty"`
	DeletedDate *time.Time `db:"deleted_date,omitempty"`
}

type farmWithUser struct {
	Farm      farm
	UserId    uint64 `db:"id_user"`
	UserName  string `db:"user_name"`
	UserEmail string `db:"user_email"`
}

type FarmRepository interface {
	Save(farm domain.Farm) (domain.Farm, error)
	FindById(id uint64) (domain.Farm, error)
	Update(farm domain.Farm) (domain.Farm, error)
	FindAllByCoords(points domain.Points, p domain.Pagination) (domain.Farms, error)
	FindAll(pag domain.Pagination) (domain.Farms, error)
	Delete(id uint64) error
	mapModelToDomainWithoutUser(m farm) domain.Farm
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
	farmR, err := r.findFarmWithUser(farm.Id)
	if err != nil {
		return domain.Farm{}, err
	}
	return farmR, nil
}

func (r farmRepository) FindById(id uint64) (domain.Farm, error) {
	farm, err := r.findFarmWithUser(id)
	if err != nil {
		return domain.Farm{}, err
	}
	return farm, nil
}

func (r farmRepository) Update(farm domain.Farm) (domain.Farm, error) {
	u := r.mapDomainToModel(farm)
	u.UpdatedDate = time.Now()
	err := r.coll.Find(db.Cond{"id": u.Id}).Update(&u)
	if err != nil {
		return domain.Farm{}, err
	}
	farmR, err := r.findFarmWithUser(farm.Id)
	if err != nil {
		return domain.Farm{}, err
	}
	return farmR, nil
}

func (r farmRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r farmRepository) FindAll(p domain.Pagination) (domain.Farms, error) {
	farms, err := r.findFarmsWithUsers(db.Cond{"farms.deleted_date": nil}, p)
	if err != nil {
		return domain.Farms{}, err
	}

	return farms, nil
}

func (r farmRepository) FindAllByCoords(points domain.Points, p domain.Pagination) (domain.Farms, error) {
	offers, err := r.offerRepo.FindByCategory(points.Category)
	if err != nil {
		return domain.Farms{}, err
	}

	ids := make([]uint64, len(offers))
	for i, item := range offers {
		ids[i] = item.Farm.Id
	}

	farms, err := r.findFarmsWithUsers(db.Cond{"farms.deleted_date": nil,
		"farms.id IN":       ids,
		"farms.latitude <":  points.UpperLeftPoint.Lat,
		"farms.latitude >":  points.BottomRightPoint.Lat,
		"farms.longitude <": points.UpperLeftPoint.Lng,
		"farms.longitude >": points.BottomRightPoint.Lng}, p)
	if err != nil {
		return domain.Farms{}, err
	}

	return farms, nil
}

func (r farmRepository) findFarmsWithUsers(cond db.Cond, p domain.Pagination) (domain.Farms, error) {
	var farms []farmWithUser
	query := r.coll.Session().SQL().Select("farms.*", "u.id AS id_user", "u.name AS user_name", "u.email AS user_email").
		From("farms").
		Where(cond).
		Join("users as u").On("u.id = farms.user_id")
	res := query.Paginate(uint(p.CountPerPage))
	err := res.Page(uint(p.Page)).All(&farms)
	if err != nil {
		return domain.Farms{}, err
	}

	domainFarms := make([]domain.Farm, len(farms))
	for i, farm := range farms {
		domainFarms[i] = r.mapModelToDomain(farm.Farm, user{Id: farm.UserId, Name: farm.UserName, Email: farm.UserEmail})
	}
	farmsR := domain.Farms{Items: domainFarms}
	totalCount, err := res.TotalEntries()
	if err != nil {
		return domain.Farms{}, err
	}

	farmsR.Total = totalCount
	farmsR.Pages = uint(math.Ceil(float64(farmsR.Total) / float64(p.CountPerPage)))

	return farmsR, nil
}

func (r farmRepository) findFarmWithUser(farmId uint64) (domain.Farm, error) {
	var farm farmWithUser
	err := r.coll.Session().SQL().Select("farms.*", "u.id AS id_user", "u.name AS user_name", "u.email AS user_email").
		From("farms").
		Where(db.Cond{"farms.id": farmId, "farms.deleted_date": nil}).
		Join("users as u").On("u.id = farms.user_id").One(&farm)
	if err != nil {
		return domain.Farm{}, err
	}

	return r.mapModelToDomain(farm.Farm, user{Id: farm.UserId, Name: farm.UserName, Email: farm.UserEmail}), nil
}

func (r farmRepository) GetAllImages(farmId uint64) []domain.Image {
	var offers []offer
	var offersId []uint64
	var coverImages []image
	var additionalImages []image
	offersQuery := r.coll.Session().SQL().Select("*").From("offers").Where("farm_id = ?", farmId)
	err := offersQuery.All(&offers)
	if err != nil {
		return []domain.Image{}
	}

	for _, offer := range offers {
		offersId = append(offersId, offer.Id)
		coverImages = append(coverImages, image{Name: offer.Cover})
	}

	imagesQuery := r.coll.Session().SQL().Select("*").From("images").Where("entity = ? AND entity_id IN ?", "offers", offersId)
	err = imagesQuery.All(&additionalImages)
	if err != nil {
		return []domain.Image{}
	}

	return mapImageModelToDomainList(append(coverImages, additionalImages...))
}

func (r farmRepository) mapDomainToModel(m domain.Farm) farm {
	return farm{
		Id:          m.Id,
		Name:        m.Name,
		City:        m.City,
		Address:     m.Address,
		CreatedDate: m.CreatedDate,
		UserId:      m.User.Id,
		Latitude:    m.Latitude,
		Longitude:   m.Longitude,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}
}

func (r farmRepository) mapModelToDomain(m farm, u user) domain.Farm {
	return domain.Farm{
		Id:          m.Id,
		Name:        m.Name,
		City:        m.City,
		Address:     m.Address,
		CreatedDate: m.CreatedDate,
		User:        mapModelToDomainUser(u),
		Latitude:    m.Latitude,
		Longitude:   m.Longitude,
		AllImages:   r.GetAllImages(m.Id),
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}
}

func (r farmRepository) mapModelToDomainWithoutUser(m farm) domain.Farm {
	return domain.Farm{
		Id:          m.Id,
		Name:        m.Name,
		City:        m.City,
		Address:     m.Address,
		CreatedDate: m.CreatedDate,
		User:        domain.User{Id: m.UserId},
		Latitude:    m.Latitude,
		Longitude:   m.Longitude,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}
}
