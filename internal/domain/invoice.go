package domain

import "time"

type Invoice struct {
	InvoiceId       string
	Status          InvoiceStatus
	Amount          *float64
	FinalAmount     *float64
	FailureReason   *string
	ErrCode         *string
	CreatedDate     time.Time
	UpdatedDate     time.Time
	CancelListItems []CancelListItem
}

type CancelListItem struct {
	InvoiceId    string
	Status       string
	Amount       *float64
	ApprovalCode *string
	Rrn          *string
	CreatedDate  time.Time
	UpdatedDate  time.Time
}

type InvoiceStatus string

var (
	INVOICE_STATUS_CREATED    InvoiceStatus = "created"    //рахунок створено успішно, очікується оплата
	INVOICE_STATUS_PROCESSING InvoiceStatus = "processing" //платіж обробляється
	INVOICE_STATUS_HOLD       InvoiceStatus = "hold"       //сума заблокована
	INVOICE_STATUS_SUCCESS    InvoiceStatus = "success"    //успішна оплата
	INVOICE_STATUS_FAILURE    InvoiceStatus = "failure"    //неуспішна оплата
	INVOICE_STATUS_REVERSED   InvoiceStatus = "reversed"   //оплата повернена після успіху
	INVOICE_STATUS_EXPIRED    InvoiceStatus = "expired"    //час дії вичерпано
)
