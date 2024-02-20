package database

import (
	"boilerplate/internal/domain"
	"errors"
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
	FindByFarmUserId(farmUserId uint64, p domain.Pagination) (domain.Orders, error)
}

type orderRepository struct {
	orderItemRepo OrderItemRepository
	coll          db.Collection
}

func NewOrderRepository(dbSession db.Session, orderItemRepo OrderItemRepository) OrderRepository {
	return orderRepository{
		orderItemRepo: orderItemRepo,
		coll:          dbSession.Collection(OrdersTableName),
	}
}

func (r orderRepository) Save(order domain.Order) (domain.Order, error) {
	o := r.mapDomainToModel(order)
	exists, err := r.coll.Find(db.Cond{"deleted_date": nil, "status": domain.DRAFT, "user_id": order.UserId}).Exists()
	if err != nil || exists {
		return domain.Order{}, errors.New("user already have an order in DRAFT status")
	}
	ordrItmsModel, ProdPrice, err := r.orderItemRepo.PrepareAllToSave(order.OrderItems, o.UserId)
	if err != nil {
		return domain.Order{}, err
	}
	o.Status = string(domain.DRAFT)
	o.CreatedDate, o.UpdatedDate = time.Now(), time.Now()
	err = r.coll.InsertReturning(&o)
	if err != nil {
		return domain.Order{}, err
	}
	err = r.orderItemRepo.SaveAll(ordrItmsModel, o.Id)
	if err != nil {
		return domain.Order{}, err
	}

	o.ProductsPrice = ProdPrice
	o.TotalPrice = math.Round((ProdPrice+order.ShippingPrice)*100) / 100
	err = r.coll.Find(db.Cond{"id": o.Id}).Update(&o)
	if err != nil {
		return domain.Order{}, err
	}

	orderDomain := r.mapModelToDomain(o)
	orderDomain.OrderItemsCount = uint64(len(ordrItmsModel))
	if err != nil {
		return domain.Order{}, err
	}
	return orderDomain, nil
}

func (r orderRepository) FindById(id uint64) (domain.Order, error) {
	var o order
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&o)
	if err != nil {
		return domain.Order{}, err
	}
	orderDomain := r.mapModelToDomain(o)
	count, err := r.orderItemRepo.Count(o.Id)
	if err != nil {
		return domain.Order{}, err
	}
	orderDomain.OrderItemsCount = count

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
		count, err := or.orderItemRepo.Count(item.Id)
		if err != nil {
			return domain.Orders{}, err
		}
		orders.Items[i].OrderItemsCount = count
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
	count, err := or.orderItemRepo.Count(o.Id)
	if err != nil {
		return domain.Order{}, err
	}
	orderDomain.OrderItemsCount = count

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
	var order order
	result := r.coll.Find(db.Cond{"id": orderId, "deleted_date": nil})
	err := result.One(&order)
	if err != nil {
		return err
	}
	totalPrice, err := r.orderItemRepo.GetTotalPriceByOrder(orderId)
	if err != nil {
		return err
	}
	order.ProductsPrice = totalPrice
	order.TotalPrice = math.Round((totalPrice+order.ShippingPrice)*100) / 100
	err = result.Update(&order)
	if err != nil {
		return err
	}

	return nil
}

func (r orderRepository) FindByFarmUserId(farmUserId uint64, p domain.Pagination) (domain.Orders, error) {
	var farmsId []uint64
	var offersId []uint64
	var ordersId []uint64
	var orders []order
	err := r.coll.Session().SQL().Select("id").From("farms").Where("user_id = ?", farmUserId).All(&farmsId)
	if err != nil {
		return domain.Orders{}, err
	}

	err = r.coll.Session().SQL().Select("id").From("offers").Where("farm_id IN ?", farmsId).All(&offersId)
	if err != nil {
		return domain.Orders{}, err
	}

	err = r.coll.Session().SQL().Select("order_id").From("order_items").Where("offer_id IN ?", offersId).All(&ordersId)
	if err != nil {
		return domain.Orders{}, err
	}

	err = r.coll.Session().SQL().Select("*").From("orders").Where("id IN ?", ordersId).All(&orders)
	if err != nil {
		return domain.Orders{}, err
	}

	paginatedOrders := r.mapModelToDomainPagination(orders)
	if err != nil {
		return domain.Orders{}, err
	}

	for i, item := range paginatedOrders.Items {
		count, err := r.orderItemRepo.Count(item.Id)
		if err != nil {
			return domain.Orders{}, err
		}
		paginatedOrders.Items[i].OrderItemsCount = count
	}

	paginatedOrders.Total = uint64(len(orders))
	paginatedOrders.Pages = uint(math.Ceil(float64(paginatedOrders.Total) / float64(p.CountPerPage)))
	return paginatedOrders, nil
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
