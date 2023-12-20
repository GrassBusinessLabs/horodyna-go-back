package domain

import "time"

// Address представляє структуру для адреси.
type Address struct {
	ID          uint64
	UserID      uint64
	Street      string
	City        string
	State       string
	ZipCode     string
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
