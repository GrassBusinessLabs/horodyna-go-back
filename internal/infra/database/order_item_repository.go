package database

import (
	"boilerplate/internal/domain"
	"errors"
	"math"
	"time"

	"github.com/upper/db/v4"
)

const OrderItemsTableName = "order_items"

type orderItem struct {
	Id          uint64     `db:"id,omitempty"`
	Title       string     `db:"title"`
	Price       float64    `db:"price"`
	TotalPrice  float64    `db:"total_price"`
	Amount      uint32     `db:"amount"`
	OrderId     uint64     `db:"order_id"`
	OfferId     uint64     `db:"offer_id"`
	CreatedDate time.Time  `db:"created_date,omitempty"`
	UpdatedDate time.Time  `db:"updated_date,omitempty"`
	DeletedDate *time.Time `db:"deleted_date,omitempty"`
}

type orderItemWithFarm struct {
	OrderItem orderItem
	FarmId    uint64 `db:"farm_id"`
	Name      string `db:"farm_name"`
	City      string `db:"farm_city"`
	Address   string `db:"farm_address"`
	UserId    uint64 `db:"user_id"`
}

type OrderItemRepository interface {
	Save(ords domain.OrderItem, orderId uint64) (domain.OrderItem, error)
	Count(orderId uint64) (uint64, error)
	SaveAll(ords []orderItem, orderId uint64) error
	PrepareAllToSave(ords []domain.OrderItem, orderUserId uint64) ([]orderItem, float64, error)
	Update(ords domain.OrderItem) (domain.OrderItem, error)
	FindAllWithoutPagination(id uint64) ([]domain.OrderItem, error)
	GetTotalPriceByOrder(orderId uint64) (float64, error)
	FindById(id uint64) (domain.OrderItem, error)
	DeleteByOrder(orderId uint64) error
	Delete(oiId uint64) error
}

type orderItemRepository struct {
	offerRepo OfferRepository
	farmRepo  FarmRepository
	coll      db.Collection
	orderColl db.Collection
	sess      db.Session
}

func NewOrderItemRepository(dbSession db.Session, offerR OfferRepository, farmR FarmRepository) OrderItemRepository {
	return orderItemRepository{
		offerRepo: offerR,
		farmRepo:  farmR,
		coll:      dbSession.Collection(OrderItemsTableName),
		orderColl: dbSession.Collection(OrdersTableName),
		sess:      dbSession,
	}
}

func (r orderItemRepository) Count(orderId uint64) (uint64, error) {
	count, err := r.coll.Find(db.Cond{"order_id": orderId}).Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r orderItemRepository) FindById(id uint64) (domain.OrderItem, error) {
	var o orderItem
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&o)
	if err != nil {
		return domain.OrderItem{}, err
	}

	order, err := r.mapModelToDomain(o)
	if err != nil {
		return domain.OrderItem{}, err
	}

	return order, nil
}

func (r orderItemRepository) SaveAll(ords []orderItem, orderId uint64) error {
	for _, item := range ords {
		item.OrderId = orderId
		err := r.coll.InsertReturning(&item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r orderItemRepository) PrepareAllToSave(ords []domain.OrderItem, orderUserId uint64) ([]orderItem, float64, error) {
	modelItms := make([]orderItem, len(ords))
	var prodPrice float64
	for i, item := range ords {
		offer, err := r.offerRepo.FindById(item.OfferId)
		if err != nil {
			return []orderItem{}, 0, err
		}

		if offer.Stock < uint(item.Amount) {
			return []orderItem{}, 0, errors.New("The orderitem amount can`t be more than in offer.")
		}
		if offer.User.Id == orderUserId {
			return []orderItem{}, 0, errors.New("The owner of the offer can`t buy his products.")
		}

		item.Title = offer.Title
		item.OfferId = offer.Id
		item.Price = offer.Price
		item.TotalPrice = math.Round(offer.Price*float64(item.Amount)*100) / 100
		o := r.mapDomainToModel(item)
		o.CreatedDate, o.UpdatedDate = time.Now(), time.Now()
		if err != nil {
			return []orderItem{}, 0, err
		}
		prodPrice += o.TotalPrice
		modelItms[i] = o
	}

	return modelItms, prodPrice, nil
}

func (r orderItemRepository) Save(ords domain.OrderItem, orderId uint64) (domain.OrderItem, error) {
	ords.Order.Id = orderId
	offer, err := r.offerRepo.FindById(ords.OfferId)
	if err != nil {
		return domain.OrderItem{}, err
	}

	if offer.Stock < uint(ords.Amount) {
		return domain.OrderItem{}, errors.New("The orderitem amount can`t be more than in offer.")
	}
	exists, err := r.coll.Find(db.Cond{"order_id": orderId, "offer_id": ords.OfferId, "deleted_date": nil}).Exists()
	if err == nil && exists {
		return domain.OrderItem{}, errors.New("The order already have an offer with this id.")
	}

	ords.Title = offer.Title
	ords.Price = offer.Price
	ords.TotalPrice = math.Round(offer.Price*float64(ords.Amount)*100) / 100
	o := r.mapDomainToModel(ords)
	o.CreatedDate, o.UpdatedDate = time.Now(), time.Now()
	err = r.coll.InsertReturning(&o)
	if err != nil {
		return domain.OrderItem{}, err
	}
	order, err := r.mapModelToDomain(o)
	if err != nil {
		return domain.OrderItem{}, err
	}
	farm, err := r.GetFarmByOfferId(order.OfferId)
	if err != nil {
		return domain.OrderItem{}, err
	}
	order.Farm = farm

	return order, nil
}

func (r orderItemRepository) GetTotalPriceByOrder(orderId uint64) (float64, error) {
	var total float64
	row, err := r.sess.SQL().QueryRow("SELECT SUM(total_price) FROM order_items WHERE order_id = ? AND deleted_date IS NULL", orderId)
	if err != nil {
		return 0, err
	}
	err = row.Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r orderItemRepository) Update(ords domain.OrderItem) (domain.OrderItem, error) {
	offer, err := r.offerRepo.FindById(ords.OfferId)
	if err != nil {
		return domain.OrderItem{}, err
	}

	if offer.Stock < uint(ords.Amount) {
		return domain.OrderItem{}, errors.New("The orderitem amount can`t be more than in offer.")
	}

	o := r.mapDomainToModel(ords)
	o.UpdatedDate = time.Now()
	err = r.coll.Find(db.Cond{"id": o.Id}).Update(&o)
	if err != nil {
		return domain.OrderItem{}, err
	}

	order, err := r.mapModelToDomain(o)
	if err != nil {
		return domain.OrderItem{}, err
	}
	farm, err := r.GetFarmByOfferId(order.OfferId)
	if err != nil {
		return domain.OrderItem{}, err
	}
	order.Farm = farm

	return order, nil
}

func (r orderItemRepository) DeleteByOrder(orderId uint64) error {
	err := r.coll.Find(db.Cond{"order_id": orderId, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
	return err
}

func (r orderItemRepository) Delete(oiId uint64) error {
	err := r.coll.Find(db.Cond{"id": oiId, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
	if err != nil {
		return err
	}

	return err
}

func (r orderItemRepository) FindAllWithoutPagination(orderId uint64) ([]domain.OrderItem, error) {
	var orderItems []orderItemWithFarm
	err := r.sess.SQL().Select("oi", "*", "f.id AS farm_id", "f.name AS farm_name", "f.city AS farm_city", "f.address AS farm_address", "f.user_id").
		From("order_items AS oi").
		Where("oi.order_id = ? AND oi.deleted_date IS NULL", orderId).
		Join("offers AS o").On("o.id = oi.offer_id").
		Join("farms AS f").On("f.id = o.farm_id").All(&orderItems)
	if err != nil {
		return []domain.OrderItem{}, err
	}
	domainOrderItems := make([]domain.OrderItem, len(orderItems))
	for i, oi := range orderItems {
		orderItem := r.mapModelToDomainWithoutOrder(oi.OrderItem, farm{
			Id:      oi.FarmId,
			Name:    oi.Name,
			City:    oi.City,
			Address: oi.Address,
			UserId:  oi.UserId,
		})
		domainOrderItems[i] = orderItem
	}
	return domainOrderItems, nil
}

func (r orderItemRepository) GetFarmByOfferId(offerId uint64) (domain.Farm, error) {
	var farmModel farm
	err := r.sess.SQL().Select("farms", "*").From("offers").
		Where("offers.id = ?", offerId).
		Join("farms").On("offers.farm_id = farms.id").
		One(&farmModel)
	if err != nil {
		return domain.Farm{}, err
	}

	return r.farmRepo.mapModelToDomain(farmModel), nil
}

func (r orderItemRepository) FindOrderWithTwoFields(orderId uint64) (domain.Order, error) {
	var o order
	err := r.orderColl.Find(db.Cond{"id": orderId}).Select("id", "user_id").One(&o)
	if err != nil {
		return domain.Order{}, err
	}

	return MapModelToDomain(o), nil
}

func (r orderItemRepository) mapDomainToModel(m domain.OrderItem) orderItem {
	return orderItem{
		Id:          m.Id,
		Price:       m.Price,
		TotalPrice:  m.TotalPrice,
		Amount:      m.Amount,
		Title:       m.Title,
		OrderId:     m.Order.Id,
		OfferId:     m.OfferId,
		CreatedDate: m.CreatedDate,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}
}

func (r orderItemRepository) mapModelToDomain(m orderItem) (domain.OrderItem, error) {
	order, err := r.FindOrderWithTwoFields(m.OrderId)
	if err != nil {
		return domain.OrderItem{}, err
	}

	return domain.OrderItem{
		Id:          m.Id,
		Price:       m.Price,
		TotalPrice:  m.TotalPrice,
		Amount:      m.Amount,
		Title:       m.Title,
		Order:       order,
		OfferId:     m.OfferId,
		CreatedDate: m.CreatedDate,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}, nil
}

func (r orderItemRepository) mapModelToDomainWithoutOrder(m orderItem, f farm) domain.OrderItem {
	return domain.OrderItem{
		Id:          m.Id,
		Price:       m.Price,
		TotalPrice:  m.TotalPrice,
		Amount:      m.Amount,
		Title:       m.Title,
		Farm:        r.farmRepo.mapModelToDomain(f),
		OfferId:     m.OfferId,
		CreatedDate: m.CreatedDate,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}
}
