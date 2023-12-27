package database

import (
	"boilerplate/internal/domain"
	"math"
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
	Status        string     `db:"status"`
	CreatedDate   time.Time  `db:"created_date,omitempty"`
	UpdatedDate   time.Time  `db:"updated_date,omitempty"`
	DeletedDate   *time.Time `db:"deleted_date,omitempty"`
}

type OrderRepository interface {
	Save(ordr domain.Order) (domain.Order, error)
	FindById(id uint64) (domain.Order, error)
	Update(order domain.Order) (domain.Order, error)
	FindAllByUserId(userId uint64, p domain.Pagination) (domain.Orders, error)
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
	o.Status = string(domain.DRAFT)
	o.CreatedDate, o.UpdatedDate = time.Now(), time.Now()
	ordrItmsModel, ProdPrice, err := r.orderItemRepo.PrepareAllToSave(order.OrderItems, o.UserId)
	if err != nil {
		return domain.Order{}, err
	}

	err = r.coll.InsertReturning(&o)
	if err != nil {
		return domain.Order{}, err
	}

	err = r.orderItemRepo.SaveAll(ordrItmsModel, o.Id)
	if err != nil {
		return domain.Order{}, err
	}

	order.ProductsPrice = ProdPrice
	order.TotalPrice = math.Round((ProdPrice+order.ShippingPrice)*100) / 100
	order.Id = o.Id
	o = r.mapDomainToModel(order)
	err = r.coll.Find(db.Cond{"id": o.Id}).Update(&o)
	if err != nil {
		return domain.Order{}, err
	}

	orderDomain := r.mapModelToDomain(o)
	orderItems, err := r.orderItemRepo.FindAllWithoutPagination(o.Id)
	if err != nil {
		return domain.Order{}, err
	}

	orderDomain.OrderItems = orderItems
	return orderDomain, nil
}

func (r orderRepository) FindById(id uint64) (domain.Order, error) {
	var o order
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&o)
	if err != nil {
		return domain.Order{}, err
	}
	orderDomain := r.mapModelToDomain(o)
	orderItems, err := r.orderItemRepo.FindAllWithoutPagination(o.Id)
	if err != nil {
		return domain.Order{}, err
	}
	orderDomain.OrderItems = orderItems

	return orderDomain, nil
}

func (or orderRepository) FindAllByUserId(userId uint64, p domain.Pagination) (domain.Orders, error) {
	var data []order
	query := or.coll.Find(db.Cond{"user_id": userId, "deleted_date": nil})
	res := query.Paginate(uint(p.CountPerPage))
	err := res.Page(uint(p.Page)).All(&data)
	if err != nil {
		return domain.Orders{}, err
	}

	orders := or.mapModelToDomainPagination(data)
	if err != nil {
		return domain.Orders{}, err
	}

	for i, item := range orders.Items {
		orderItems, err := or.orderItemRepo.FindAllWithoutPagination(item.Id)
		if err != nil {
			return domain.Orders{}, err
		}
		orders.Items[i].OrderItems = orderItems
	}

	totalCount, err := res.TotalEntries()
	if err != nil {
		return domain.Orders{}, err
	}

	orders.Total = totalCount
	orders.Pages = uint(math.Ceil(float64(orders.Total) / float64(p.CountPerPage)))
	return orders, nil
}

func (or orderRepository) Update(req domain.Order) (domain.Order, error) {
	var err error
	o := or.mapDomainToModel(req)
	o.UpdatedDate = time.Now()
	err = or.coll.Find(db.Cond{"id": o.Id}).Update(&o)
	if err != nil {
		return domain.Order{}, err
	}
	orderDomain := or.mapModelToDomain(o)
	orderItems, err := or.orderItemRepo.FindAllWithoutPagination(o.Id)
	if err != nil {
		return domain.Order{}, err
	}
	orderDomain.OrderItems = orderItems

	return orderDomain, nil
}

func (r orderRepository) Delete(order domain.Order) error {
	err := r.orderItemRepo.DeleteByOrder(order.Id)
	if err != nil {
		return err
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
	order.TotalPrice = math.Round((ProdPrice+order.ShippingPrice)*100) / 100
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
		Status:        string(o.Status),
		CreatedDate:   o.CreatedDate,
		UpdatedDate:   o.UpdatedDate,
		DeletedDate:   o.DeletedDate,
	}
}

func (r orderRepository) mapModelToDomain(o order) domain.Order {
	return domain.Order{
		Id:            o.Id,
		Comment:       o.Comment,
		UserId:        o.UserId,
		AddressId:     o.AddressId,
		ProductsPrice: o.ProductsPrice,
		ShippingPrice: o.ShippingPrice,
		TotalPrice:    o.TotalPrice,
		Status:        domain.OrderStatus(o.Status),
		OrderItems:    make([]domain.OrderItem, 0),
		CreatedDate:   o.CreatedDate,
		UpdatedDate:   o.UpdatedDate,
		DeletedDate:   o.DeletedDate,
	}
}

func MapModelToDomain(ord order) domain.Order {
	return domain.Order{
		Id:     ord.Id,
		UserId: ord.UserId,
	}
}

func (f orderRepository) mapModelToDomainPagination(orders []order) domain.Orders {
	newOrders := make([]domain.Order, len(orders))
	for i, order := range orders {
		newOrders[i] = f.mapModelToDomain(order)
	}
	return domain.Orders{Items: newOrders}
}
