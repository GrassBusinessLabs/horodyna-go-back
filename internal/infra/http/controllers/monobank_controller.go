package controllers

import (
	"boilerplate/internal/app"
	"boilerplate/internal/infra/monobank"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type MonobankController struct {
	monobankService app.MonobankService
}

func NewMonobankController(ms app.MonobankService) MonobankController {
	return MonobankController{
		monobankService: ms,
	}
}

func (c MonobankController) CreateInvoice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validatorInstance := validator.New()
		var request monobank.CreateInvoiceRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Printf("MonobankController Decode: %s", err)
			BadRequest(w, err)
			return
		}

		err = validatorInstance.Struct(request)
		if err != nil {
			log.Printf("MonobankController Struct: %s", err)
			BadRequest(w, err)
			return
		}

		response, err := c.monobankService.CreateInvoice(request)
		if err != nil {
			log.Printf("MonobankController CreateInvoice: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, response)
	}
}

func (c MonobankController) GetInvoiceData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		invoiceId, err := strconv.ParseUint(chi.URLParam(r, "invoiceId"), 10, 64)
		if err != nil {
			log.Printf("MonobankController ParseUint: %s", err)
			BadRequest(w, err)
			return
		}

		response, err := c.monobankService.GetInvoiceData(string(rune(invoiceId)))
		if err != nil {
			log.Printf("MonobankController GetInvoiceData: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, response)
	}
}

func (c MonobankController) CancelSuccessfulInvoice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validatorInstance := validator.New()
		var request monobank.CancelSuccessfulInvoiceRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			log.Printf("MonobankController Decode: %s", err)
			BadRequest(w, err)
			return
		}

		err = validatorInstance.Struct(request)
		if err != nil {
			log.Printf("MonobankController Struct: %s", err)
			BadRequest(w, err)
			return
		}

		response, err := c.monobankService.CancelSuccessfulInvoice(request)
		if err != nil {
			log.Printf("MonobankController CancelSuccessfulInvoice: %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, response)
	}
}
