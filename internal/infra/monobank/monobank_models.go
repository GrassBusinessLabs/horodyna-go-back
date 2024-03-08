package monobank

import (
	"boilerplate/internal/domain"
	"time"
)

type Invoice struct {
	InvoiceId     string
	Status        domain.InvoiceStatus
	FailureReason string
	ErrCode       string
	Amount        int64  //сума у мінімальних одиницях валюти (1 грн = 100 коп)
	Ccy           *int32 //валюта
	FinalAmount   int64  //підсумкова сума у мінімальних одиницях валюти, змінюється після оплати та повернень
	CreatedDate   time.Time
	ModifiedDate  *time.Time
	Reference     *string //Референс платежу, який визначається продавцем
	CancelList    *[]CancelListItem
}

type CancelListItem struct {
	Status       CancelListItemStatus
	Amount       int64 //сума у мінімальних одиницях валюти (1 грн = 100 коп)
	Ccy          int32 //ISO 4217 код валюти
	CreatedDate  time.Time
	ModifiedDate *time.Time
	ApprovalCode string //Код авторизації
	Rrn          string //Ідентифікатор транзакції в платіжній системі
	ExtRef       string //Референс операції скасування, який було вказано продавцем
}

type CancelListItemStatus string

var (
	CANCEL_LIST_ITEM_STATUS_PROCESSING CancelListItemStatus = "processing" //заява на скасування знаходиться в обробці
	CANCEL_LIST_ITEM_STATUS_SUCCESS    CancelListItemStatus = "success"    //заяву на скасування виконано успішно
	CANCEL_LIST_ITEM_STATUS_FAILURE    CancelListItemStatus = "failure"    //неуспішне скасування
)
