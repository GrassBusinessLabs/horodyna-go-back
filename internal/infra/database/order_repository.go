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
	Id               uint64     `db:"id,omitempty"`
	Comment          string     `db:"comment"`
	UserId           uint64     `db:"user_id"`
	Address          *string    `db:"address"`
	ProductsPrice    float64    `db:"products_price"`
	ShippingPrice    float64    `db:"shipping_price"`
	TotalPrice       float64    `db:"total_price"`
	Status           string     `db:"status"`
	PostOffice       *string    `db:"post_office"`
	PostOfficeCity   *string    `db:"post_office_city"`
	Ttn              *string    `db:"ttn"`
	IsPercentagePaid *bool      `db:"is_percentage_paid"`
	CreatedDate      time.Time  `db:"created_date,omitempty"`
	UpdatedDate      time.Time  `db:"updated_date,omitempty"`
	DeletedDate      *time.Time `db:"deleted_date,omitempty"`
}

type OrderRepository interface {
	Save(ordr domain.Order) (domain.Order, error)
	FindById(id uint64) (domain.Order, error)
	Update(order domain.Order) (domain.Order, error)
	FindAllByUserId(userId uint64, p domain.Pagination) (domain.Orders, error)
	Delete(order domain.Order) error
	Recalculate(orderId uint64) error
	GetOrdersByFarmUserId(farmUserId uint64, p domain.Pagination) (domain.Orders, error)
	SplitOrderByFarms(order domain.Order) (map[uint64]domain.Order, error)
	SubmitSplitedOrder(order domain.Order, farmId uint64) (domain.Order, error)
	DeleteSplitedOrder(order domain.Order, farmId uint64) error
	GetActiveOrdersByFarmId(farmId uint64) ([]domain.Order, error)
	GetFarmerOrdersPercentage(farmUserId uint64) ([]domain.Order, float64, error)
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
	exists, err := r.coll.Find(db.Cond{"deleted_date": nil, "status": domain.DRAFT, "user_id": order.User.Id}).Exists()
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

func (r orderRepository) SplitOrderByFarms(order domain.Order) (map[uint64]domain.Order, error) {
	farmOrderItems := make(map[uint64][]domain.OrderItem)
	for _, orderItem := range order.OrderItems {
		_, keyExists := farmOrderItems[orderItem.Farm.Id]
		if keyExists {
			farmOrderItems[orderItem.Farm.Id] = append(farmOrderItems[orderItem.Farm.Id], orderItem)
		} else {
			farmOrderItems[orderItem.Farm.Id] = []domain.OrderItem{orderItem}
		}
	}

	splitedOrders := make(map[uint64]domain.Order)
	for farmId, orderItems := range farmOrderItems {
		orderItemsModel, productPrice, err := r.orderItemRepo.PrepareAllToSave(orderItems, order.User.Id)
		if err != nil {
			return make(map[uint64]domain.Order, 0), err
		}

		splitedOrder := domain.Order{
			Comment:         order.Comment,
			User:            order.User,
			Address:         order.Address,
			OrderItems:      orderItems,
			OrderItemsCount: uint64(len(orderItemsModel)),
			ProductsPrice:   productPrice,
			TotalPrice:      math.Round((productPrice+order.ShippingPrice)*100) / 100,
			ShippingPrice:   order.ShippingPrice,
			Status:          domain.DRAFT,
			PostOffice:      order.PostOffice,
			Ttn:             order.Ttn,
		}
		splitedOrders[farmId] = splitedOrder
	}

	return splitedOrders, nil
}

func (r orderRepository) SubmitSplitedOrder(order domain.Order, farmId uint64) (domain.Order, error) {
	splitedOrders, err := r.SplitOrderByFarms(order)
	if err != nil {
		return domain.Order{}, err
	}

	splitedOrder, exists := splitedOrders[farmId]
	if !exists {
		return domain.Order{}, errors.New("no such farm in splited orders")
	}

	splitedOrderModel := r.mapDomainToModel(splitedOrder)
	splitedOrderModel.Status = string(domain.SUBMITTED)
	splitedOrderModel.CreatedDate, splitedOrderModel.UpdatedDate = time.Now(), time.Now()
	err = r.coll.InsertReturning(&splitedOrderModel)
	if err != nil {
		return domain.Order{}, err
	}

	submittedSplitedOrder := r.mapModelToDomain(splitedOrderModel)
	for _, orderItem := range splitedOrder.OrderItems {
		orderItem.Order = submittedSplitedOrder
		orderItem, err = r.orderItemRepo.Update(orderItem)
		if err != nil {
			return domain.Order{}, err
		}
	}

	err = r.Recalculate(submittedSplitedOrder.Id)
	if err != nil {
		return domain.Order{}, err
	}

	err = r.Recalculate(order.Id)
	if err != nil {
		return domain.Order{}, err
	}

	submittedSplitedOrder.OrderItems, err = r.orderItemRepo.FindAllWithoutPagination(submittedSplitedOrder.Id)
	if err != nil {
		return domain.Order{}, err
	}

	return submittedSplitedOrder, nil
}

func (r orderRepository) DeleteSplitedOrder(order domain.Order, farmId uint64) error {
	splitedOrders, err := r.SplitOrderByFarms(order)
	if err != nil {
		return err
	}

	splitedOrder, exists := splitedOrders[farmId]
	if !exists {
		return errors.New("no such farm in splited orders")
	}

	for _, orderItem := range splitedOrder.OrderItems {
		err = r.orderItemRepo.Delete(orderItem.Id)
		if err != nil {
			return err
		}
	}

	err = r.Recalculate(order.Id)
	if err != nil {
		return err
	}

	return nil
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

func (r orderRepository) FindAllByUserId(userId uint64, p domain.Pagination) (domain.Orders, error) {
	var data []order
	query := r.coll.Find(db.Cond{"user_id": userId, "deleted_date": nil}).OrderBy("created_date")
	res := query.Paginate(uint(p.CountPerPage))
	err := res.Page(uint(p.Page)).All(&data)
	if err != nil {
		return domain.Orders{}, err
	}

	orders := r.mapModelToDomainPagination(data)
	if err != nil {
		return domain.Orders{}, err
	}

	for i, item := range orders.Items {
		orderItems, err := r.orderItemRepo.FindAllWithoutPagination(item.Id)
		if err != nil {
			return domain.Orders{}, err
		}
		orders.Items[i].OrderItemsCount = uint64(len(orderItems))
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

func (r orderRepository) GetOrdersByFarmUserId(farmUserId uint64, p domain.Pagination) (domain.Orders, error) {
	var orders []order
	query := r.coll.Session().SQL().
		Select("orders.*").
		From("orders").
		Join("order_items").On("order_items.order_id = orders.id").
		Join("offers").On("order_items.offer_id = offers.id").
		Join("farms").On("offers.farm_id = farms.id").
		Where(db.Cond{"farms.user_id": farmUserId, "orders.deleted_date": nil, "orders.status !=": "DRAFT"}).
		OrderBy("created_date").
		Distinct()

	err := query.All(&orders)
	if err != nil {
		return domain.Orders{}, err
	}

	paginatedOrders := r.mapModelToDomainPagination(orders)
	if err != nil {
		return domain.Orders{}, err
	}

	for i, item := range paginatedOrders.Items {
		orderItems, err := r.orderItemRepo.FindAllWithoutPagination(item.Id)
		if err != nil {
			return domain.Orders{}, err
		}
		paginatedOrders.Items[i].OrderItemsCount = uint64(len(orderItems))
		paginatedOrders.Items[i].OrderItems = orderItems
	}

	paginatedOrders.Total = uint64(len(orders))
	paginatedOrders.Pages = uint(math.Ceil(float64(paginatedOrders.Total) / float64(p.CountPerPage)))
	return paginatedOrders, nil
}

func (r orderRepository) GetActiveOrdersByFarmId(farmId uint64) ([]domain.Order, error) {
	activeStatusses := domain.GetActiveOrderStatuses()
	stringActiveStatusses := make([]string, len(activeStatusses))
	for i, status := range activeStatusses {
		stringActiveStatusses[i] = string(status)
	}

	var activeOrders []order
	query := r.coll.Session().SQL().Select("orders.*").
		From("orders").
		Join("order_items").On("order_items.order_id = orders.id").
		Join("offers").On("offers.id = order_items.offer_id").
		Where(db.Cond{"offers.farm_id": farmId, "orders.status IN": stringActiveStatusses, "orders.deleted_date": nil}).
		Distinct()

	err := query.All(&activeOrders)
	if err != nil {
		return []domain.Order{}, err
	}

	return r.mapModelToDomainCollection(activeOrders), nil
}

func (r orderRepository) GetFarmerOrdersPercentage(farmUserId uint64) ([]domain.Order, float64, error) {
	var orders []order
	query := r.coll.Session().SQL().
		Select("orders.*").
		From("orders").
		Join("order_items").On("order_items.order_id = orders.id").
		Join("offers").On("order_items.offer_id = offers.id").
		Join("farms").On("offers.farm_id = farms.id").
		Where(db.Cond{"farms.user_id": farmUserId, "orders.deleted_date": nil, "orders.status": "COMPLETED", "orders.is_percentage_paid": false}).
		Distinct()
	err := query.All(&orders)
	if err != nil {
		return []domain.Order{}, 0, err
	}

	domainOrders := r.mapModelToDomainCollection(orders)
	var total float64 = 0
	for i, domainOrder := range domainOrders {
		percentage := domainOrder.TotalPrice / 10
		domainOrders[i].Percentage = &percentage
		total += percentage
	}

	return domainOrders, total, nil
}

func (r orderRepository) mapDomainToModel(o domain.Order) order {

	return order{
		Id:               o.Id,
		Comment:          o.Comment,
		UserId:           o.User.Id,
		Address:          o.Address,
		ProductsPrice:    o.ProductsPrice,
		ShippingPrice:    o.ShippingPrice,
		TotalPrice:       o.TotalPrice,
		Status:           string(o.Status),
		PostOffice:       o.PostOffice,
		PostOfficeCity:   o.PostOfficeCity,
		Ttn:              o.Ttn,
		IsPercentagePaid: o.IsPercentagePaid,
		CreatedDate:      o.CreatedDate,
		UpdatedDate:      o.UpdatedDate,
		DeletedDate:      o.DeletedDate,
	}
}

func (r orderRepository) mapModelToDomain(o order) domain.Order {
	var user user
	err := r.coll.Session().SQL().Select("*").From("users").Where("id = ?", o.UserId).One(&user)
	if err != nil {
		return domain.Order{}
	}

	return domain.Order{
		Id:               o.Id,
		Comment:          o.Comment,
		User:             mapModelToDomainUser(user),
		Address:          o.Address,
		ProductsPrice:    o.ProductsPrice,
		ShippingPrice:    o.ShippingPrice,
		TotalPrice:       o.TotalPrice,
		Status:           domain.OrderStatus(o.Status),
		PostOffice:       o.PostOffice,
		PostOfficeCity:   o.PostOfficeCity,
		Ttn:              o.Ttn,
		OrderItems:       make([]domain.OrderItem, 0),
		IsPercentagePaid: o.IsPercentagePaid,
		CreatedDate:      o.CreatedDate,
		UpdatedDate:      o.UpdatedDate,
		DeletedDate:      o.DeletedDate,
	}
}

func MapModelToDomain(ord order) domain.Order {
	return domain.Order{
		Id:   ord.Id,
		User: domain.User{Id: ord.UserId},
	}
}

func (f orderRepository) mapModelToDomainCollection(orders []order) []domain.Order {
	newOrders := make([]domain.Order, len(orders))
	for i, order := range orders {
		newOrders[i] = f.mapModelToDomain(order)
	}
	return newOrders
}

func (f orderRepository) mapModelToDomainPagination(orders []order) domain.Orders {
	newOrders := make([]domain.Order, len(orders))
	for i, order := range orders {
		newOrders[i] = f.mapModelToDomain(order)
	}
	return domain.Orders{Items: newOrders}
}
