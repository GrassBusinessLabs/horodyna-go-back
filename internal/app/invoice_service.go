package app

import (
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/database"
	"log"
)

type InvoiceService interface {
	Save(invoice domain.Invoice) (domain.Invoice, error)
	Update(ref, req domain.Invoice) (domain.Invoice, error)
	Upsert(invoice domain.Invoice) (domain.Invoice, error)
	Find(uint64) (interface{}, error)
	FindAll() ([]domain.Invoice, error)
	FindAllUpdatedWithinOneDay() ([]domain.Invoice, error)
	Delete(invoiceId string) error
}

type invoiceService struct {
	invoiceRepository database.InvoiceRepository
}

func NewInvoiceService(ir database.InvoiceRepository) InvoiceService {
	return invoiceService{
		invoiceRepository: ir,
	}
}

func (s invoiceService) Save(invoice domain.Invoice) (domain.Invoice, error) {
	return s.invoiceRepository.Save(invoice)
}

func (s invoiceService) Update(ref, req domain.Invoice) (domain.Invoice, error) {
	if req.Status != "" {
		ref.Status = req.Status
	}
	if !req.CreatedDate.IsZero() {
		ref.CreatedDate = req.CreatedDate
	}
	if !req.UpdatedDate.IsZero() {
		ref.UpdatedDate = req.UpdatedDate
	}

	invoice, err := s.invoiceRepository.Update(ref)
	if err != nil {
		log.Printf("InvoiceService: %s", err)
		return domain.Invoice{}, err
	}

	return invoice, nil
}

func (s invoiceService) Upsert(invoice domain.Invoice) (domain.Invoice, error) {
	invoice, err := s.invoiceRepository.Upsert(invoice)
	if err != nil {
		log.Printf("InvoiceService: %s", err)
		return domain.Invoice{}, err
	}

	return invoice, nil
}

func (s invoiceService) Find(id uint64) (interface{}, error) {
	invoice, err := s.invoiceRepository.FindOne(string(rune(id)))
	if err != nil {
		log.Printf("InvoiceService: %s", err)
		return domain.Invoice{}, err
	}

	return invoice, nil
}

func (s invoiceService) FindAll() ([]domain.Invoice, error) {
	invoices, err := s.invoiceRepository.FindAll()
	if err != nil {
		log.Printf("InvoiceService: %s", err)
		return []domain.Invoice{}, err
	}

	return invoices, nil
}

func (s invoiceService) FindAllUpdatedWithinOneDay() ([]domain.Invoice, error) {
	invoices, err := s.invoiceRepository.FindAllUpdatedWithinOneDay()
	if err != nil {
		log.Printf("InvoiceService: %s", err)
		return []domain.Invoice{}, err
	}

	return invoices, nil
}

func (s invoiceService) Delete(invoiceId string) error {
	err := s.invoiceRepository.Delete(invoiceId)
	if err != nil {
		log.Printf("InvoiceService: %s", err)
		return err
	}

	return nil
}
