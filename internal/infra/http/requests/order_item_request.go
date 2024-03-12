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
		Amount: m.Amount,
		Offer:  domain.Offer{Id: m.OfferId},
	}, nil
}

func (m OrderItemRequest) ToDomainModelNotInterface() (domain.OrderItem, error) {
	return domain.OrderItem{
		Amount: m.Amount,
		Offer:  domain.Offer{Id: m.OfferId},
	}, nil
}

func (m OrderItemRequest) ToDomainModelArray(arr []OrderItemRequest) ([]domain.OrderItem, error) {
	var err error
	new := make([]domain.OrderItem, len(arr))
	for i, item := range arr {
		new[i], err = item.ToDomainModelNotInterface()
		if err != nil {
			return []domain.OrderItem{}, err
		}
	}

	return new, nil
}
