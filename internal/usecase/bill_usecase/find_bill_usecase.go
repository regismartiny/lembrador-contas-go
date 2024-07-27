package bill_usecase

import (
	"context"
	"time"

	"github.com/regismartiny/lembrador-contas-go/internal/entity/bill_entity"
	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
)

func FindBillUseCase(billRepository bill_entity.BillRepositoryInterface) BillUseCaseInterface {
	return &BillUseCase{
		billRepository,
	}
}

type BillOutputDTO struct {
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

func (u *BillUseCase) FindBillById(
	ctx context.Context, id string) (*BillOutputDTO, *internal_error.InternalError) {
	billEntity, err := u.billRepository.FindBillById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &BillOutputDTO{
		Id:              billEntity.Id,
		Name:            billEntity.Name,
		Company:         billEntity.Company,
		ValueSourceType: bill_entity.ValueSourceType(billEntity.ValueSourceType).Name(),
		ValueSourceId:   billEntity.ValueSourceId,
		DueDay:          billEntity.DueDay,
		Status:          bill_entity.BillStatus(billEntity.Status).Name(),
		CreatedAt:       billEntity.CreatedAt,
		UpdatedAt:       billEntity.UpdatedAt,
	}, nil
}

func (u *BillUseCase) FindBills(
	ctx context.Context,
	status bill_entity.BillStatus,
	name, company string) ([]BillOutputDTO, *internal_error.InternalError) {
	billEntities, err := u.billRepository.FindBills(
		ctx, status, name, company)
	if err != nil {
		return nil, err
	}

	var billOutputs []BillOutputDTO
	for _, value := range billEntities {
		billOutputs = append(billOutputs, BillOutputDTO{
			Id:              value.Id,
			Name:            value.Name,
			Company:         value.Company,
			ValueSourceType: bill_entity.ValueSourceType(value.ValueSourceType).Name(),
			ValueSourceId:   value.ValueSourceId,
			DueDay:          value.DueDay,
			Status:          bill_entity.BillStatus(value.Status).Name(),
			CreatedAt:       value.CreatedAt,
			UpdatedAt:       value.UpdatedAt,
		})
	}

	return billOutputs, nil
}
