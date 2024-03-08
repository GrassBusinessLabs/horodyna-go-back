package controllers

import (
	"boilerplate/internal/app"
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/http/resources"
	"log"
	"net/http"
)

type InvoiceController struct {
	invoiceService app.InvoiceService
}

func NewInvoiceController(is app.InvoiceService) InvoiceController {
	return InvoiceController{
		invoiceService: is,
	}
}

func (c InvoiceController) FindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		invoices, err := c.invoiceService.FindAll()
		if err != nil {
			log.Printf("InvoiceController: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.InvoiceDto{}.DomainToDtoPaginatedCollection(invoices))
	}
}

func (c InvoiceController) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		invoice := r.Context().Value(InvoiceKey).(domain.Invoice)
		Success(w, resources.InvoiceDto{}.DomainToDto(invoice))
	}
}

func (c InvoiceController) FindAllUpdatedWithinOneDay() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		invoices, err := c.invoiceService.FindAllUpdatedWithinOneDay()
		if err != nil {
			log.Printf("InvoiceController: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, resources.InvoiceDto{}.DomainToDtoPaginatedCollection(invoices))
	}
}

func (c InvoiceController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		invoice := r.Context().Value(InvoiceKey).(domain.Invoice)
		err := c.invoiceService.Delete(invoice.InvoiceId)
		if err != nil {
			log.Printf("InvoiceController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}
