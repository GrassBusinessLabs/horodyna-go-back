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

type OrderItemRepository interface {
	Save(ords domain.OrderItem, orderId uint64) (domain.OrderItem, error)
	SaveAll(ords []orderItem, orderId uint64) error
	PrepareAllToSave(ords []domain.OrderItem, orderUserId uint64) ([]orderItem, float64, error)
	Update(ords domain.OrderItem) (domain.OrderItem, error)
	FindAllWithoutPagination(id uint64) ([]domain.OrderItem, error)
	FindById(id uint64) (domain.OrderItem, error)
	DeleteByOrder(orderId uint64) error
	Delete(oiId uint64) error
}

type orderItemRepository struct {
	offerRepo OfferRepository
	farmRepo  FarmRepository
	coll      db.Collection
	orderColl db.Collection
}

func NewOrderItemRepository(dbSession db.Session, offerR OfferRepository, farmR FarmRepository) OrderItemRepository {
	return orderItemRepository{
		offerRepo: offerR,
		farmRepo:  farmR,
		coll:      dbSession.Collection(OrderItemsTableName),
		orderColl: dbSession.Collection(OrdersTableName),
	}
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
	return order, nil
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
	var orderItems []orderItem
	err := r.coll.Find(db.Cond{"order_id": orderId, "deleted_date": nil}).All(&orderItems)
	if err != nil {
		return []domain.OrderItem{}, err
	}

	newOrderItems, err := r.mapModelToDomainMass(orderItems)
	if err != nil {
		return []domain.OrderItem{}, err
	}

	for i, order := range newOrderItems {
		offer, err := r.offerRepo.FindById(order.OfferId)
		if err != nil {
			return []domain.OrderItem{}, err
		}

		farm, err := r.farmRepo.FindById(offer.Farm.Id)
		if err != nil {
			return []domain.OrderItem{}, err
		}
		newOrderItems[i].Farm = farm
	}
	return newOrderItems, nil
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

func (o orderItemRepository) mapModelToDomainMass(orderItems []orderItem) ([]domain.OrderItem, error) {
	newOrderItems := make([]domain.OrderItem, len(orderItems))
	var err error
	for i, orderItem := range orderItems {
		newOrderItems[i], err = o.mapModelToDomain(orderItem)
		if err != nil {
			return []domain.OrderItem{}, err
		}
	}
	return newOrderItems, nil
}
