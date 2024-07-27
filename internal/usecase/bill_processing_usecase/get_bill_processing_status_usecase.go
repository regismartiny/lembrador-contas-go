package bill_processing_usecase

import (
	"context"

	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
)

type GetBillProcessingStatusOutputDTO struct {
	Status string `json:"status"`
}

func (u *BillProcessingUseCase) GetBillProcessingStatus(
	ctx context.Context,
	id string) (GetBillProcessingStatusOutputDTO, *internal_error.InternalError) {

	billProcessing, err := u.billProcessingRepository.FindBillProcessingById(ctx, id)

	if err != nil {
		return GetBillProcessingStatusOutputDTO{}, err
	}

	return GetBillProcessingStatusOutputDTO{
		Status: billProcessing.Status.Name()}, nil
}
