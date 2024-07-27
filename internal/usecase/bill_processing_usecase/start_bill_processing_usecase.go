package bill_processing_usecase

import (
	"context"
	"log"
	"time"

	"github.com/regismartiny/lembrador-contas-go/internal/entity/bill_processing_entity"
	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
)

type BillProcessingInputDTO struct {
}

type StartBillProcessingOutputDTO struct {
	BillProcessingId string `json:"id"`
}

type BillProcessingUseCaseInterface interface {
	StartBillProcessing(
		ctx context.Context,
		billProcessingInput BillProcessingInputDTO) (StartBillProcessingOutputDTO, *internal_error.InternalError)
	GetBillProcessingStatus(
		ctx context.Context,
		billProcessingId string) (GetBillProcessingStatusOutputDTO, *internal_error.InternalError)
}

type BillProcessingUseCase struct {
	billProcessingRepository bill_processing_entity.BillProcessingRepositoryInterface
}

func NewBillProcessingUseCase(
	billProcessingRepository bill_processing_entity.BillProcessingRepositoryInterface) BillProcessingUseCaseInterface {
	return &BillProcessingUseCase{
		billProcessingRepository: billProcessingRepository,
	}
}

func (u *BillProcessingUseCase) StartBillProcessing(
	ctx context.Context,
	billProcessingInput BillProcessingInputDTO) (StartBillProcessingOutputDTO, *internal_error.InternalError) {

	billProcessing, err := bill_processing_entity.CreateBillProcessing("")
	if err != nil {
		return StartBillProcessingOutputDTO{}, err
	}

	if err := u.billProcessingRepository.CreateBillProcessing(ctx, billProcessing); err != nil {
		return StartBillProcessingOutputDTO{}, err
	}

	go func() {

		go startProcessing()

		time.Sleep(30 * time.Second)

		billProcessing, err := u.billProcessingRepository.FindBillProcessingById(ctx, billProcessing.Id)
		if err != nil {
			log.Println("Error trying to find billProcessing", err)
			return
		}

		if !billProcessing.Status.IsFinished() {
			log.Println("Bill processing timeout")
			billProcessing.Status = bill_processing_entity.Timeout

			if err := u.billProcessingRepository.UpdateBillProcessing(ctx, billProcessing); err != nil {
				log.Println("Error trying to update billProcessing", err)
			}
		}

	}()

	return StartBillProcessingOutputDTO{
		BillProcessingId: billProcessing.Id}, nil
}

func startProcessing() {
	log.Println("Bill processing started")

}
