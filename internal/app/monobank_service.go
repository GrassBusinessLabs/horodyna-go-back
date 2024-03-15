package app

import (
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/monobank"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	apiKeyHeader                  = "X-Token"
	contentType                   = "application/json"
	apiCreateInvoiceUrl           = "https://api.monobank.ua/api/merchant/invoice/create"
	apiGetInvoiceDataUrl          = "https://api.monobank.ua/api/merchant/invoice/status?invoiceId=%s"
	apiCancelSuccessfulInvoiceUrl = "https://api.monobank.ua/api/merchant/invoice/cancel"
)

type MonobankService interface {
	CreateInvoice(request monobank.CreateInvoiceRequest) (monobank.CreateInvoiceResponse, error)
	GetInvoiceData(invoiceId string) (monobank.GetInvoiceDataResponse, error)
	CancelSuccessfulInvoice(request monobank.CancelSuccessfulInvoiceRequest) (monobank.CancelSuccessfulInvoiceResponse, error)
}

type monobankService struct {
	privateKey     string
	invoiceService InvoiceService
}

func NewMonobankService(privateKey string, is InvoiceService) MonobankService {
	return monobankService{
		privateKey:     privateKey,
		invoiceService: is,
	}
}

func (s monobankService) CreateInvoice(request monobank.CreateInvoiceRequest) (monobank.CreateInvoiceResponse, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		log.Printf("json.Marshal(monobankService.CreateInvoice): %s", err)
		return monobank.CreateInvoiceResponse{}, err
	}

	resp, err := s.makeHttpRequest(http.MethodPost, apiCreateInvoiceUrl, requestBody)
	if err != nil {
		log.Printf("s.makeHttpRequest(monobankService.CreateInvoice): %s", err)
		return monobank.CreateInvoiceResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Printf("Body.Close(monobankService.CreateInvoice): %s", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		var errorResponse monobank.ErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			log.Printf("json.NewDecoder(monobankService.CreateInvoice): %s", err)
			return monobank.CreateInvoiceResponse{}, err
		}
		log.Printf("monobankService.CreateInvoice: %s", errorResponse.ErrCode+" : "+errorResponse.ErrText)
		return monobank.CreateInvoiceResponse{}, errors.New(errorResponse.ErrText)
	}

	var createInvoiceResponse monobank.CreateInvoiceResponse
	err = json.NewDecoder(resp.Body).Decode(&createInvoiceResponse)
	if err != nil {
		log.Printf("json.NewDecoder(monobankService.CreateInvoice): %s", err)
		return monobank.CreateInvoiceResponse{}, err
	}

	invoice := domain.Invoice{
		InvoiceId: createInvoiceResponse.InvoiceId,
	}

	_, err = s.invoiceService.Save(invoice)
	if err != nil {
		log.Printf("s.invoiceService.Save(monobankService.CreateInvoice): %s", err)
		return monobank.CreateInvoiceResponse{}, err
	}

	return createInvoiceResponse, nil
}

func (s monobankService) GetInvoiceData(invoiceId string) (monobank.GetInvoiceDataResponse, error) {
	resp, err := s.makeHttpRequest(http.MethodGet, fmt.Sprintf(apiGetInvoiceDataUrl, invoiceId), nil)
	if err != nil {
		log.Printf("s.makeHttpRequest(monobankService.GetInvoiceData): %s", err)
		return monobank.GetInvoiceDataResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Printf("Body.Close(monobankService.GetInvoiceData): %s", err)
		}
	}(resp.Body)

	var getInvoiceDataResponse monobank.GetInvoiceDataResponse
	err = json.NewDecoder(resp.Body).Decode(&getInvoiceDataResponse)
	if err != nil {
		log.Printf("json.NewDecoder(monobankService.GetInvoiceData): %s", err)
		return monobank.GetInvoiceDataResponse{}, err
	}

	return getInvoiceDataResponse, nil
}

func (s monobankService) CancelSuccessfulInvoice(request monobank.CancelSuccessfulInvoiceRequest) (monobank.CancelSuccessfulInvoiceResponse, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		log.Printf("json.Marshal(monobankService.CancelSuccessfulInvoice): %s", err)
		return monobank.CancelSuccessfulInvoiceResponse{}, err
	}

	resp, err := s.makeHttpRequest(http.MethodPost, apiCancelSuccessfulInvoiceUrl, requestBody)
	if err != nil {
		log.Printf("s.makeHttpRequest(monobankService.CancelSuccessfulInvoice): %s", err)
		return monobank.CancelSuccessfulInvoiceResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Printf("Body.Close(monobankService.CancelSuccessfulInvoice): %s", err)
		}
	}(resp.Body)

	var cancelSuccessfulInvoiceResponse monobank.CancelSuccessfulInvoiceResponse
	err = json.NewDecoder(resp.Body).Decode(&cancelSuccessfulInvoiceResponse)
	if err != nil {
		log.Printf("json.NewDecoder(monobankService.CancelSuccessfulInvoice): %s", err)
		return monobank.CancelSuccessfulInvoiceResponse{}, err
	}

	return cancelSuccessfulInvoiceResponse, nil
}

func (s monobankService) makeHttpRequest(method, url string, requestBody []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	log.Print(s.privateKey)
	req.Header.Set(apiKeyHeader, s.privateKey)
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}
