package monobank

type MerchantPaymInfoItem struct {
	Reference      *string  `json:"reference"`      // Номер чека, замовлення, тощо; визначається мерчантом
	Destination    *string  `json:"destination"`    // Призначення платежу
	Comment        *string  `json:"comment"`        // Службове інформаційне поле
	CustomerEmails []string `json:"customerEmails"` // Масив пошт, на які потрібно відправити фіскальний чек, якщо у мерчанта активна звʼязка з checkbox
}

type CreateInvoiceRequest struct {
	Amount           int64                 `json:"amount" validation:"required"` // Сума оплати у мінімальних одиницях (копійки для гривні)
	Ccy              *int32                `json:"ccy"`                          // ISO 4217 код валюти, за замовчуванням 980 (гривня)
	MerchantPaymInfo *MerchantPaymInfoItem `json:"merchantPaymInfo"`             // Інформаційні дані замовлення, яке буде оплачуватсь. Обовʼязково вказувати при активній звʼязці з ПРРО
	RedirectUrl      *string               `json:"redirectUrl"`                  // URL, на який буде перенаправлено після оплати
	WebHookUrl       *string               `json:"webHookUrl"`                   // Адреса для CallBack (POST) – на цю адресу буде надіслано дані про стан платежу при кожній зміні статусу.
	Validity         *int64                `json:"validity"`                     // Термін дії рахунку в секундах. За замовчуванням 86400 (24 години).
	PaymentType      *string               `json:"paymentType"`                  // Тип операції. Default: "debit". Possible values: "debit", "hold". Для значення hold термін складає 9 днів
}

type CancelSuccessfulInvoiceRequest struct {
	InvoiceId string  `json:"invoiceId" validation:"required"` // Ідентифікатор рахунку
	ExtRef    *string `json:"extRef"`                          // Референс операції скасування, який визначається продавцем
	Amount    *int64  `json:"amount"`                          // Сума скасування у мінімальних одиницях валюти (копійки для гривні)
}
