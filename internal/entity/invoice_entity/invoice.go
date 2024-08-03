package invoice_entity

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
)

type Invoice struct {
	Id        string
	BillId    string
	DueDate   string
	Amount    float64
	Status    InvoiceStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type InvoiceStatus uint8

const (
	Unpaid InvoiceStatus = iota + 1
	Paid
)

func (s InvoiceStatus) Name() string {
	return invoiceStatusNames[s]
}

var invoiceStatusNames = []string{
	"",
	"unpaid",
	"paid",
}

func GetInvoiceStatusByName(name string) (InvoiceStatus, *internal_error.InternalError) {
	for k, v := range invoiceStatusNames {
		if v == name {
			return InvoiceStatus(k), nil
		}
	}

	return InvoiceStatus(0), internal_error.NewBadRequestError("invalid invoice status name")
}

func CreateInvoice(
	billId string,
	dueDate string,
	amount float64,
	status string) (*Invoice, *internal_error.InternalError) {

	var invoiceStatus InvoiceStatus

	if status == "" {
		invoiceStatus = Unpaid
	} else {
		status, err := GetInvoiceStatusByName(status)
		if err != nil {
			return nil, err
		}
		invoiceStatus = status
	}

	invoice :=
		&Invoice{
			Id:        uuid.New().String(),
			BillId:    billId,
			DueDate:   dueDate,
			Amount:    amount,
			Status:    invoiceStatus,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

	if err := invoice.Validate(); err != nil {
		return nil, err
	}

	return invoice, nil
}

func (invoice *Invoice) Validate() *internal_error.InternalError {
	if uuid.Validate(invoice.BillId) == nil {
		return internal_error.NewBadRequestError("invalid invoice object. invalid bill id")
	}
	if _, err := time.Parse("2006-01-02", invoice.DueDate); err != nil {
		return internal_error.NewBadRequestError("invalid invoice object. invalid due date")
	}

	return nil
}

type InvoiceRepositoryInterface interface {
	CreateInvoice(ctx context.Context, invoiceEntity *Invoice) *internal_error.InternalError
	FindInvoiceById(ctx context.Context, invoiceId string) (*Invoice, *internal_error.InternalError)
	FindInvoices(
		ctx context.Context,
		billId string,
		status InvoiceStatus) ([]*Invoice, *internal_error.InternalError)
	DeleteInvoices(
		ctx context.Context,
		billId string,
		status InvoiceStatus,
		dueDate string) (uint, *internal_error.InternalError)
}
