package domain

import (
	"time"
)

type Address struct {
	Id          uint64
	User        User
	City        string
	Country     string
	Address     string
	Department  string
	Lat         float64
	Lon         float64
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

func (a Address) GetUserId() uint64 {
	return a.User.Id
}
