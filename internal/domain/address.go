package domain

import "time"

// Address представляє структуру для адреси.
type Address struct {
	ID          uint64    `db:"id"`
	UserID      uint64    `db:"user_id"`
	Street      string    `db:"street"`
	City        string    `db:"city"`
	State       string    `db:"state"`
	ZipCode     string    `db:"zip_code"`
	CreatedDate time.Time `db:"created_date"`
	UpdatedDate time.Time `db:"updated_date"`
	DeletedDate time.Time `db:"deleted_date"`
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
