package domain

import (
	"time"
)

type Offer struct {
	Id          uint64
	Title       string
	Description string
	Category    string
	Price       float64
	Unit        string
	Stock       uint
	Status      bool
	User        User
	Farm        Farm
	Cover       Image
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

type Offers struct {
	Items []Offer
	Total uint64
	Pages uint
}

func (o Offer) GetUserId() uint64 {
	return o.User.Id
}
