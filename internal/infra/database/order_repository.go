package database

import (
	"boilerplate/internal/domain"
	"time"

	"github.com/upper/db/v4"
)

const OrdersTableName = "orders"

type order struct {
	Id            uint64     `db:"id,omitempty"`
	Comment       string     `db:"comment"`
	UserId        uint64     `db:"user_id"`
	AddressId     uint64     `db:"address_id"`
	ProductsPrice float64    `db:"products_price"`
	ShippingPrice float64    `db:"shipping_price"`
	TotalPrice    float64    `db:"total_price"`
	Status        bool       `db:"status"`
	CreatedDate   time.Time  `db:"created_date,omitempty"`
	UpdatedDate   time.Time  `db:"updated_date,omitempty"`
	DeletedDate   *time.Time `db:"deleted_date,omitempty"`
}

type OrderRepository interface {
	Save(ordr domain.Order) (domain.Order, error)
	FindById(id uint64) (domain.Order, error)
	Update(order domain.Order) (domain.Order, error)
	Delete(order domain.Order) error
	Recalculate(orderId uint64) error
}

type orderRepository struct {
	orderItemRepo OrderItemRepository
	coll          db.Collection
}

func NewOrderRepository(dbSession db.Session, order_item_repo OrderItemRepository) OrderRepository {
	return orderRepository{
		orderItemRepo: order_item_repo,
		coll:          dbSession.Collection(OrdersTableName),
	}
}

func (r orderRepository) Save(order domain.Order) (domain.Order, error) {
	o := r.mapDomainToModel(order)
	o.CreatedDate, o.UpdatedDate = time.Now(), time.Now()
	err := r.coll.InsertReturning(&o)
	if err != nil {
		return domain.Order{}, err
	}

	var ProdPrice float64
	for _, item := range order.OrderItems {
		created, err := r.orderItemRepo.Save(item, o.Id)
		if err != nil {
			return domain.Order{}, err
		}
		ProdPrice += created.TotalPrice
	}

	order.ProductsPrice = ProdPrice
	order.TotalPrice = ProdPrice + order.ShippingPrice
	order.Id = o.Id
	o = r.mapDomainToModel(order)
	err = r.coll.Find(db.Cond{"id": o.Id}).Update(&o)
	if err != nil {
		return domain.Order{}, err
	}

	return r.mapModelToDomain(o), nil
}

func (r orderRepository) FindById(id uint64) (domain.Order, error) {
	var o order
	err := r.coll.Find(db.Cond{"id": id}).One(&o)
	if err != nil {
		return domain.Order{}, err
	}
	return r.mapModelToDomain(o), nil
}

func (or orderRepository) Update(req domain.Order) (domain.Order, error) {
	var err error
	o := or.mapDomainToModel(req)
	o.UpdatedDate = time.Now()
	err = or.coll.Find(db.Cond{"id": o.Id}).Update(&o)
	if err != nil {
		return domain.Order{}, err
	}
	return or.mapModelToDomain(o), nil
}

func (r orderRepository) Delete(order domain.Order) error {
	for _, item := range order.OrderItems {
		err := r.orderItemRepo.Delete(item)
		if err != nil {
			return err
		}
	}
	return r.coll.Find(db.Cond{"id": order.Id, "deleted_date": nil}).Update(map[string]interface{}{"deleted_date": time.Now()})
}

func (r orderRepository) Recalculate(orderId uint64) error {
	order, err := r.FindById(orderId)
	if err != nil {
		return err
	}
	var ProdPrice float64
	for _, i := range order.OrderItems {
		ProdPrice += i.TotalPrice
	}

	order.ProductsPrice = ProdPrice
	order.TotalPrice = ProdPrice + order.ShippingPrice
	o := r.mapDomainToModel(order)
	err = r.coll.Find(db.Cond{"id": order.Id}).Update(&o)
	if err != nil {
		return err
	}
	return nil
}

func (r orderRepository) mapDomainToModel(o domain.Order) order {
	return order{
		Id:            o.Id,
		Comment:       o.Comment,
		UserId:        o.UserId,
		AddressId:     o.AddressId,
		ProductsPrice: o.ProductsPrice,
		ShippingPrice: o.ShippingPrice,
		TotalPrice:    o.TotalPrice,
		Status:        o.Status,
		CreatedDate:   o.CreatedDate,
		UpdatedDate:   o.UpdatedDate,
		DeletedDate:   o.DeletedDate,
	}
}

func (r orderRepository) mapModelToDomain(o order) domain.Order {
	order_items, err := r.orderItemRepo.FindAllWithoutPagination(o.Id)
	if err != nil {
		return domain.Order{}
	}

	return domain.Order{
		Id:            o.Id,
		Comment:       o.Comment,
		UserId:        o.UserId,
		AddressId:     o.AddressId,
		OrderItems:    order_items,
		ProductsPrice: o.ProductsPrice,
		ShippingPrice: o.ShippingPrice,
		TotalPrice:    o.TotalPrice,
		Status:        o.Status,
		CreatedDate:   o.CreatedDate,
		UpdatedDate:   o.UpdatedDate,
		DeletedDate:   o.DeletedDate,
	}
}
