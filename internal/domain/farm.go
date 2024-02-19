package domain

import (
	"time"
)

type Point struct {
	Lat float64
	Lng float64
}

type Points struct {
	UpperLeftPoint   Point
	BottomRightPoint Point
	Category         string
}

type Farm struct {
	Id          uint64
	Name        *string
	City        string
	Address     string
	User        User
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
	return f.User.Id
}
