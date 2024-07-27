package invoice_usecase

import (
	"context"
	"time"

	"github.com/regismartiny/lembrador-contas-go/internal/entity/invoice_entity"
	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
)

type InvoiceInputDTO struct {
	Name    string  `json:"name" binding:"required,min=3"`
	DueDate string  `json:"dueDate" binding:"required"`
	Amount  float64 `json:"amount" binding:"required"`
	Status  string  `json:"status"`
}

type CreateInvoiceOutputDTO struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	DueDate   time.Time `json:"dueDate"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt" time_format:"2006-01-02 15:04:05"`
	UpdatedAt time.Time `json:"updatedAt" time_format:"2006-01-02 15:04:05"`
}

type InvoiceUseCaseInterface interface {
	CreateInvoice(
		ctx context.Context,
		auctionInput InvoiceInputDTO) *internal_error.InternalError
	FindInvoiceById(
		ctx context.Context,
		id string) (*InvoiceOutputDTO, *internal_error.InternalError)
	FindInvoices(
		ctx context.Context,
		status invoice_entity.InvoiceStatus,
		name string) ([]InvoiceOutputDTO, *internal_error.InternalError)
}

type InvoiceUseCase struct {
	invoiceRepository invoice_entity.InvoiceRepositoryInterface
}

func NewInvoiceUseCase(
	invoiceRepository invoice_entity.InvoiceRepositoryInterface) InvoiceUseCaseInterface {
	return &InvoiceUseCase{
		invoiceRepository: invoiceRepository,
	}
}

func (u *InvoiceUseCase) CreateInvoice(
	ctx context.Context,
	invoiceInput InvoiceInputDTO) *internal_error.InternalError {

	invoice, err := invoice_entity.CreateInvoice(invoiceInput.Name, invoiceInput.DueDate, invoiceInput.Amount, invoiceInput.Status)
	if err != nil {
		return err
	}

	if err := u.invoiceRepository.CreateInvoice(ctx, invoice); err != nil {
		return err
	}

	return nil
}
