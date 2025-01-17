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
	BillId    string    `json:"billId"`
	DueDate   string    `json:"dueDate"`
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
		BillId:    invoiceEntity.BillId,
		DueDate:   invoiceEntity.DueDate,
		Amount:    invoiceEntity.Amount,
		Status:    invoice_entity.InvoiceStatus(invoiceEntity.Status).Name(),
		CreatedAt: invoiceEntity.CreatedAt,
		UpdatedAt: invoiceEntity.UpdatedAt,
	}, nil
}

func (u *InvoiceUseCase) FindInvoices(
	ctx context.Context,
	billId string,
	status invoice_entity.InvoiceStatus) ([]*InvoiceOutputDTO, *internal_error.InternalError) {
	invoiceEntities, err := u.invoiceRepository.FindInvoices(
		ctx, billId, status)
	if err != nil {
		return nil, err
	}

	invoiceOutputs := make([]*InvoiceOutputDTO, len(invoiceEntities))
	for i, value := range invoiceEntities {
		invoiceOutputs[i] = &InvoiceOutputDTO{
			Id:        value.Id,
			BillId:    value.BillId,
			DueDate:   value.DueDate,
			Amount:    value.Amount,
			Status:    invoice_entity.InvoiceStatus(value.Status).Name(),
			CreatedAt: value.CreatedAt,
			UpdatedAt: value.UpdatedAt,
		}
	}

	return invoiceOutputs, nil
}
