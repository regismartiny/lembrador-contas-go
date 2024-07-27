package invoice_usecase

import (
	"context"
	"time"

	"github.com/regismartiny/lembrador-contas-go/internal/entity/invoice_entity"
	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
)

func FindInvoiceUseCase(invoiceRepository invoice_entity.InvoiceRepositoryInterface) InvoiceUseCaseInterface {
	return &InvoiceUseCase{
		invoiceRepository,
	}
}

type InvoiceOutputDTO struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	DueDate   time.Time `json:"dueDate"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt" time_format:"2006-01-02 15:04:05"`
	UpdatedAt time.Time `json:"updatedAt" time_format:"2006-01-02 15:04:05"`
}

func (u *InvoiceUseCase) FindInvoiceById(
	ctx context.Context, id string) (*InvoiceOutputDTO, *internal_error.InternalError) {
	invoiceEntity, err := u.invoiceRepository.FindInvoiceById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &InvoiceOutputDTO{
		Id:        invoiceEntity.Id,
		Name:      invoiceEntity.Name,
		DueDate:   invoiceEntity.DueDate,
		Amount:    invoiceEntity.Amount,
		Status:    invoice_entity.InvoiceStatus(invoiceEntity.Status).Name(),
		CreatedAt: invoiceEntity.CreatedAt,
		UpdatedAt: invoiceEntity.UpdatedAt,
	}, nil
}

func (u *InvoiceUseCase) FindInvoices(
	ctx context.Context,
	status invoice_entity.InvoiceStatus,
	name string) ([]InvoiceOutputDTO, *internal_error.InternalError) {
	invoiceEntities, err := u.invoiceRepository.FindInvoices(
		ctx, status, name)
	if err != nil {
		return nil, err
	}

	var invoiceOutputs []InvoiceOutputDTO
	for _, value := range invoiceEntities {
		invoiceOutputs = append(invoiceOutputs, InvoiceOutputDTO{
			Id:        value.Id,
			Name:      value.Name,
			DueDate:   value.DueDate,
			Amount:    value.Amount,
			Status:    invoice_entity.InvoiceStatus(value.Status).Name(),
			CreatedAt: value.CreatedAt,
			UpdatedAt: value.UpdatedAt,
		})
	}

	return invoiceOutputs, nil
}
