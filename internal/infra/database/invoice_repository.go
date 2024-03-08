package database

import (
	"boilerplate/internal/domain"
	"time"

	"github.com/upper/db/v4"
)

const invoicesTableName = "invoices"

type invoice struct {
	InvoiceId       string               `db:"invoice_id,omitempty"`
	Status          domain.InvoiceStatus `db:"status,omitempty"`
	Amount          *float64             `db:"amount,omitempty"`
	FinalAmount     *float64             `db:"final_amount,omitempty"`
	FailureReason   *string              `db:"failure_reason,omitempty"`
	ErrCode         *string              `db:"err_code,omitempty"`
	CreatedDate     time.Time            `db:"created_date,omitempty"`
	UpdatedDate     time.Time            `db:"updated_date,omitempty"`
	CancelListItems []cancelListItem
}

type cancelListItem struct {
	InvoiceId    string    `db:"invoice_id,omitempty"`
	Status       string    `db:"status,omitempty"`
	Amount       *float64  `db:"amount,omitempty"`
	ApprovalCode *string   `db:"approval_code,omitempty"`
	Rrn          *string   `db:"rrn,omitempty"`
	CreatedDate  time.Time `db:"created_date,omitempty"`
	UpdatedDate  time.Time `db:"modified_date,omitempty"`
}

type InvoiceRepository interface {
	Save(invoice domain.Invoice) (domain.Invoice, error)
	Update(invoice domain.Invoice) (domain.Invoice, error)
	Upsert(invoice domain.Invoice) (domain.Invoice, error)
	FindOne(invoiceId string) (domain.Invoice, error)
	FindAll() ([]domain.Invoice, error)
	FindAllUpdatedWithinOneDay() ([]domain.Invoice, error)
	Delete(invoiceId string) error
}

type invoiceRepository struct {
	coll db.Collection
	sess db.Session
}

func NewInvoiceRepository(dbSession db.Session) InvoiceRepository {
	return invoiceRepository{
		coll: dbSession.Collection(invoicesTableName),
		sess: dbSession,
	}
}

func (r invoiceRepository) Save(invoice domain.Invoice) (domain.Invoice, error) {
	invoiceModel := r.mapDomainToModel(invoice)

	err := r.coll.InsertReturning(&invoiceModel)
	if err != nil {
		return domain.Invoice{}, err
	}

	return r.mapModelToDomain(invoiceModel), nil
}

func (r invoiceRepository) Upsert(invoice domain.Invoice) (domain.Invoice, error) {
	invoiceModel := r.mapDomainToModel(invoice)

	err := r.sess.Tx(func(tx db.Session) error {
		query, err := r.sess.SQL().
			QueryRow(`INSERT INTO invoices (invoice_id, status, amount, final_amount, failure_reason, err_code, created_date, updated_date) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT (invoice_id) 
		DO UPDATE SET status = ?, amount = ?, final_amount = ?, failure_reason = ?, err_code = ?,updated_date = ? RETURNING *`,
				invoiceModel.InvoiceId, invoiceModel.Status, invoiceModel.Amount, invoiceModel.FinalAmount, invoiceModel.FailureReason, invoiceModel.ErrCode, invoiceModel.CreatedDate, invoiceModel.UpdatedDate,
				invoiceModel.Status, invoiceModel.Amount, invoiceModel.FinalAmount, invoiceModel.FailureReason, invoiceModel.ErrCode, invoiceModel.UpdatedDate)
		if err != nil {
			return err
		}

		err = query.Scan(&invoiceModel.InvoiceId, &invoiceModel.Status, &invoiceModel.CreatedDate, &invoiceModel.UpdatedDate, &invoiceModel.FailureReason, &invoiceModel.ErrCode, &invoiceModel.Amount, &invoiceModel.FinalAmount)

		if err != nil {
			return err
		}

		for i, v := range invoiceModel.CancelListItems {
			query, err = r.sess.SQL().
				QueryRow(`INSERT INTO invoice_cancellations (invoice_id, status, amount, approval_code,rrn, created_date, updated_date)
			VALUES (?, ?, ?, ?, ?, ?, ?) ON CONFLICT (invoice_id)
			DO UPDATE SET status = ?, amount = ?, approval_code = ?, rrn = ?, updated_date = ? RETURNING *`,
					v.InvoiceId, v.Status, v.Amount, v.ApprovalCode, v.Rrn, v.CreatedDate, v.UpdatedDate,
					v.Status, v.Amount, v.ApprovalCode, v.Rrn, v.UpdatedDate)
			if err != nil {
				return err
			}

			err = query.Scan(
				&invoiceModel.CancelListItems[i].InvoiceId,
				&invoiceModel.CancelListItems[i].Status,
				&invoiceModel.CancelListItems[i].Amount,
				&invoiceModel.CancelListItems[i].ApprovalCode,
				&invoiceModel.CancelListItems[i].Rrn,
				&invoiceModel.CancelListItems[i].CreatedDate,
				&invoiceModel.CancelListItems[i].UpdatedDate)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return domain.Invoice{}, err
	}

	return r.mapModelToDomain(invoiceModel), nil
}

func (r invoiceRepository) Update(invoice domain.Invoice) (domain.Invoice, error) {
	invoiceModel := r.mapDomainToModel(invoice)

	err := r.coll.Find(db.Cond{"invoice_id": invoiceModel.InvoiceId}).Update(&invoiceModel)
	if err != nil {
		return domain.Invoice{}, err
	}

	return r.mapModelToDomain(invoiceModel), nil
}

func (r invoiceRepository) FindOne(invoiceId string) (domain.Invoice, error) {
	var invoiceModel invoice

	err := r.coll.Find(db.Cond{"invoice_id": invoiceId}).One(&invoiceModel)
	if err != nil {
		return domain.Invoice{}, err
	}

	return r.mapModelToDomain(invoiceModel), nil
}

func (r invoiceRepository) FindAll() ([]domain.Invoice, error) {
	var invoiceModels []invoice

	err := r.coll.Find().All(&invoiceModels)
	if err != nil {
		return nil, err
	}

	return r.mapModelToDomainCollection(invoiceModels), nil
}

func (r invoiceRepository) FindAllUpdatedWithinOneDay() ([]domain.Invoice, error) {
	var invoiceModels []invoice

	oneDayAgo := time.Now().Add(-24 * time.Hour)

	err := r.coll.Find(db.Cond{"updated_date >=": oneDayAgo}).All(&invoiceModels)
	if err != nil {
		return nil, err
	}

	return r.mapModelToDomainCollection(invoiceModels), nil
}

func (r invoiceRepository) Delete(invoiceId string) error {
	return r.coll.Find(db.Cond{"invoice_id": invoiceId}).Delete()
}

func (r invoiceRepository) mapDomainToModel(d domain.Invoice) invoice {
	m := invoice{
		InvoiceId:     d.InvoiceId,
		Status:        d.Status,
		Amount:        d.Amount,
		FinalAmount:   d.FinalAmount,
		ErrCode:       d.ErrCode,
		FailureReason: d.FailureReason,
		CreatedDate:   d.CreatedDate,
		UpdatedDate:   d.UpdatedDate,
	}

	m.CancelListItems = make([]cancelListItem, len(d.CancelListItems))

	for i, v := range d.CancelListItems {
		m.CancelListItems[i] = cancelListItem{
			InvoiceId:    v.InvoiceId,
			Status:       v.Status,
			Amount:       v.Amount,
			ApprovalCode: v.ApprovalCode,
			Rrn:          v.Rrn,
			CreatedDate:  v.CreatedDate,
			UpdatedDate:  v.UpdatedDate,
		}
	}

	return m
}

func (r invoiceRepository) mapModelToDomain(m invoice) domain.Invoice {
	d := domain.Invoice{
		InvoiceId:     m.InvoiceId,
		Status:        m.Status,
		Amount:        m.Amount,
		FinalAmount:   m.FinalAmount,
		ErrCode:       m.ErrCode,
		FailureReason: m.FailureReason,
		CreatedDate:   m.CreatedDate,
		UpdatedDate:   m.UpdatedDate,
	}

	d.CancelListItems = make([]domain.CancelListItem, len(m.CancelListItems))

	for i, v := range m.CancelListItems {
		d.CancelListItems[i] = domain.CancelListItem{
			InvoiceId:    v.InvoiceId,
			Status:       v.Status,
			Amount:       v.Amount,
			ApprovalCode: v.ApprovalCode,
			Rrn:          v.Rrn,
			CreatedDate:  v.CreatedDate,
			UpdatedDate:  v.UpdatedDate,
		}
	}

	return d
}

func (r invoiceRepository) mapModelToDomainCollection(m []invoice) []domain.Invoice {
	var d []domain.Invoice

	for _, v := range m {
		d = append(d, r.mapModelToDomain(v))
	}

	return d
}
