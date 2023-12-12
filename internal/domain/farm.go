package domain

import (
	"time"
)

type Farm struct {
	Id          uint64
	Name        string
	City        string
	Address     string
	UserId      uint64
	Longitude   float64
	Latitude    float64
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

type Farms struct {
	Items []Farm
	Total uint64
	Pages uint
}

func (f Farm) GetUserId() uint64 {
	return f.UserId
}
