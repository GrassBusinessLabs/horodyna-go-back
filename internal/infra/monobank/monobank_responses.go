package monobank

import "time"

type ErrorResponse struct {
	ErrCode string `json:"errCode"`
	ErrText string `json:"errText"`
}

type CreateInvoiceResponse struct {
	InvoiceId string `json:"invoiceId"`
	PageUrl   string `json:"pageUrl"`
}

type GetInvoiceDataResponse struct {
	InvoiceId     string                   `json:"invoiceId"`
	Status        string                   `json:"status"`
	FailureReason *string                  `json:"failureReason,omitempty"`
	ErrCode       *string                  `json:"errCode,omitempty"`
	Amount        *int64                   `json:"amount,omitempty"`
	Ccy           int32                    `json:"ccy"`
	FinalAmount   *int64                   `json:"finalAmount,omitempty"`
	CreatedDate   *time.Time               `json:"createdDate,omitempty"`
	ModifiedDate  *time.Time               `json:"modifiedDate,omitempty"`
	Reference     *string                  `json:"reference,omitempty"`
	CancelList    []CancelListResponseItem `json:"cancelList,omitempty"`
}

type CancelListResponseItem struct {
	Status       CancelListItemStatus `json:"status"`
	Amount       *int64               `json:"amount,omitempty"` //сума у мінімальних одиницях валюти (1 грн = 100 коп)
	Ccy          *int32               `json:"ccy,omitempty"`    //ISO 4217 код валюти
	CreatedDate  time.Time            `json:"createdDate"`
	ModifiedDate time.Time            `json:"modifiedDate"`
	ApprovalCode *string              `json:"approvalCode,omitempty"` //Код авторизації
	Rrn          *string              `json:"rrn,omitempty"`          //Ідентифікатор транзакції в платіжній системі
	ExtRef       *string              `json:"extRef,omitempty"`       //Референс операції скасування, який було вказано продавцем
}

type CancelSuccessfulInvoiceResponse struct {
	Status       string    `json:"status"`
	CreatedDate  time.Time `json:"createdDate"`
	ModifiedDate time.Time `json:"modifiedDate"`
}
