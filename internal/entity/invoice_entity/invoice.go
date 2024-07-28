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
	DueDate   time.Time
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

	parsedDueDate, err := time.Parse("2006-01-02", dueDate)
	if err != nil {
		return nil, internal_error.NewBadRequestError("invalid due date")
	}

	invoice :=
		&Invoice{
			Id:        uuid.New().String(),
			BillId:    billId,
			DueDate:   parsedDueDate,
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
		return internal_error.NewBadRequestError("invalid invoice object")
	}

	return nil
}

type InvoiceRepositoryInterface interface {
	CreateInvoice(ctx context.Context, invoiceEntity *Invoice) *internal_error.InternalError
	FindInvoiceById(ctx context.Context, invoiceId string) (*Invoice, *internal_error.InternalError)
	FindInvoices(
		ctx context.Context,
		billId string,
		status InvoiceStatus) ([]Invoice, *internal_error.InternalError)
	DeleteInvoices(
		ctx context.Context,
		billId string,
		status InvoiceStatus,
		dueDate time.Time) (uint, *internal_error.InternalError)
}
