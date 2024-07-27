package bill_usecase

import (
	"context"
	"time"

	"github.com/regismartiny/lembrador-contas-go/internal/entity/bill_entity"
	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
)

type BillInputDTO struct {
	Name            string `json:"name" binding:"required,min=3"`
	Company         string `json:"company" binding:"required,min=3"`
	ValueSourceType string `json:"valueSourceType" binding:"required"`
	ValueSourceId   string `json:"valueSourceId" binding:"required"`
	DueDay          uint8  `json:"dueDay" binding:"required"`
	Status          string `json:"status"`
}

type CreateBillOutputDTO struct {
	Id              string    `json:"id"`
	Name            string    `json:"name"`
	Company         string    `json:"company"`
	ValueSourceType string    `json:"valueSourceType"`
	ValueSourceId   string    `json:"valueSourceId"`
	DueDay          uint8     `json:"dueDay"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"createdAt" time_format:"2006-01-02 15:04:05"`
	UpdatedAt       time.Time `json:"updatedAt" time_format:"2006-01-02 15:04:05"`
}

type BillUseCaseInterface interface {
	CreateBill(
		ctx context.Context,
		auctionInput BillInputDTO) *internal_error.InternalError
	FindBillById(
		ctx context.Context,
		id string) (*BillOutputDTO, *internal_error.InternalError)
	FindBills(
		ctx context.Context,
		status bill_entity.BillStatus,
		name, email string) ([]BillOutputDTO, *internal_error.InternalError)
}

type BillUseCase struct {
	billRepository bill_entity.BillRepositoryInterface
}

func NewBillUseCase(
	billRepository bill_entity.BillRepositoryInterface) BillUseCaseInterface {
	return &BillUseCase{
		billRepository: billRepository,
	}
}

func (u *BillUseCase) CreateBill(
	ctx context.Context,
	billInput BillInputDTO) *internal_error.InternalError {

	bill, err := bill_entity.CreateBill(billInput.Name, billInput.Company, billInput.ValueSourceType, billInput.ValueSourceId, billInput.DueDay, billInput.Status)
	if err != nil {
		return err
	}

	if err := u.billRepository.CreateBill(ctx, bill); err != nil {
		return err
	}

	return nil
}
