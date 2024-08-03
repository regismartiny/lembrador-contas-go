package bill_processing_usecase

import (
	"context"
	"time"

	"github.com/regismartiny/lembrador-contas-go/internal/entity/bill_processing_entity"
	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
)

type FindBillProcessingOutputDTO struct {
	Id        string    `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt" time_format:"2006-01-02 15:04:05"`
	UpdatedAt time.Time `json:"updatedAt" time_format:"2006-01-02 15:04:05"`
}

func (u *BillProcessingUseCase) FindBillProcessings(
	ctx context.Context,
	status bill_processing_entity.BillProcessingStatus) ([]*FindBillProcessingOutputDTO, *internal_error.InternalError) {

	billProcessings, err := u.billProcessingRepository.FindBillProcessings(ctx, status)
	if err != nil {
		return nil, err
	}

	billProcessingsOutput := make([]*FindBillProcessingOutputDTO, len(billProcessings))

	for i, billProcessing := range billProcessings {
		billProcessingsOutput[i] = &FindBillProcessingOutputDTO{
			Id:        billProcessing.Id,
			Status:    billProcessing.Status.Name(),
			CreatedAt: billProcessing.CreatedAt,
			UpdatedAt: billProcessing.UpdatedAt,
		}
	}

	return billProcessingsOutput, nil
}
