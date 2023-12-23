package domain

import "time"

// Address представляє структуру для адреси.
type Address struct {
	ID          uint64
	UserID      uint64
	Title       string
	City        string
	Country     string
	Address     string
	Lat         float64
	Lon         float64
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate time.Time
}

// Addresses представляє список адрес.

type Addresses struct {
	Items []Address
	Total uint64
	Pages uint
}

func (a Address) GetUserId() uint64 {
	return a.UserID
}
