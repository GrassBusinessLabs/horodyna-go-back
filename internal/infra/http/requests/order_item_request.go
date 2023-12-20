package requests

import (
	"boilerplate/internal/domain"
)

type OrderItemRequest struct {
	OfferId uint64 `json:"offer_id" validate:"required"`
	Amount  uint32 `json:"amount" validate:"required"`
}

type OrderItemUpdateRequest struct {
	Amount uint32 `json:"amount" validate:"required"`
}

func (m OrderItemUpdateRequest) ToDomainModel() (interface{}, error) {
	return domain.OrderItem{
		Amount: m.Amount,
	}, nil
}

func (m OrderItemRequest) ToDomainModel() (interface{}, error) {
	return domain.OrderItem{
		Amount:  m.Amount,
		OfferId: m.OfferId,
	}, nil
}

func (m OrderItemRequest) ToDomainModelNotInterface() (domain.OrderItem, error) {
	return domain.OrderItem{
		Amount:  m.Amount,
		OfferId: m.OfferId,
	}, nil
}
