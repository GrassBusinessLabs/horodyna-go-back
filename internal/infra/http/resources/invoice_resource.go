package resources

import (
	"boilerplate/internal/domain"
	"time"
)

type CancelListItemDto struct {
	InvoiceId    string    `json:"invoice_id"`
	Status       string    `json:"status"`
	Amount       *float64  `json:"amount"`
	ApprovalCode *string   `json:"approval_code"`
	Rrn          *string   `json:"rrn"`
	CreatedDate  time.Time `json:"created_date"`
	UpdatedDate  time.Time `json:"updated_date"`
}

type InvoiceDto struct {
	InvoiceId       string              `json:"invoice_id"`
	Status          string              `json:"status"`
	Amount          *float64            `json:"amount"`
	FinalAmount     *float64            `json:"final_amount"`
	FailureReason   *string             `json:"failure_reason"`
	ErrCode         *string             `json:"err_code"`
	CancelListItems []CancelListItemDto `json:"cancel_list_items"`
	CreatedDate     time.Time           `json:"created_date"`
	UpdatedDate     time.Time           `json:"update_date"`
}

func (d InvoiceDto) DomainToDto(invoice domain.Invoice) InvoiceDto {
	return InvoiceDto{
		InvoiceId:       invoice.InvoiceId,
		Status:          string(invoice.Status),
		Amount:          invoice.Amount,
		FinalAmount:     invoice.FinalAmount,
		FailureReason:   invoice.FailureReason,
		ErrCode:         invoice.ErrCode,
		CancelListItems: CancelListItemDto{}.DomainToDtoPaginatedCollection(invoice.CancelListItems),
		CreatedDate:     invoice.CreatedDate,
		UpdatedDate:     invoice.UpdatedDate,
	}
}

func (d InvoiceDto) DomainToDtoPaginatedCollection(invoices []domain.Invoice) []InvoiceDto {
	result := make([]InvoiceDto, len(invoices))

	for i := range invoices {
		result[i] = d.DomainToDto(invoices[i])
	}

	return result
}

func (d CancelListItemDto) DomainToDto(cancelListItem domain.CancelListItem) CancelListItemDto {
	return CancelListItemDto{
		InvoiceId:    cancelListItem.InvoiceId,
		Status:       cancelListItem.Status,
		Amount:       cancelListItem.Amount,
		ApprovalCode: cancelListItem.ApprovalCode,
		Rrn:          cancelListItem.Rrn,
		CreatedDate:  cancelListItem.CreatedDate,
		UpdatedDate:  cancelListItem.UpdatedDate,
	}
}

func (d CancelListItemDto) DomainToDtoPaginatedCollection(CancelListItems []domain.CancelListItem) []CancelListItemDto {
	result := make([]CancelListItemDto, len(CancelListItems))

	for i := range CancelListItems {
		result[i] = d.DomainToDto(CancelListItems[i])
	}

	return result
}
