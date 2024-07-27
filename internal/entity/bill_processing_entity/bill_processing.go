package bill_processing_entity

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
)

type BillProcessing struct {
	Id        string
	Status    BillProcessingStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BillProcessingStatus uint8

const (
	Started BillProcessingStatus = iota + 1
	Success
	Error
	Timeout
)

func (s BillProcessingStatus) Name() string {
	return billProcessingStatusNames[s]
}

func (s BillProcessingStatus) IsFinished() bool {
	return s == Success || s == Error || s == Timeout
}

var billProcessingStatusNames = []string{
	"",
	"started",
	"success",
	"error",
	"timeout",
}

func GetBillProcessingStatusByName(name string) (BillProcessingStatus, *internal_error.InternalError) {
	for k, v := range billProcessingStatusNames {
		if v == name {
			return BillProcessingStatus(k), nil
		}
	}

	return BillProcessingStatus(0), internal_error.NewBadRequestError("invalid billProcessing status name")
}

func CreateBillProcessing(
	status string) (*BillProcessing, *internal_error.InternalError) {

	var billProcessingStatus BillProcessingStatus

	if status == "" {
		billProcessingStatus = Started
	} else {
		status, err := GetBillProcessingStatusByName(status)
		if err != nil {
			return nil, err
		}
		billProcessingStatus = status
	}

	billProcessing :=
		&BillProcessing{
			Id:        uuid.New().String(),
			Status:    billProcessingStatus,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

	if err := billProcessing.Validate(); err != nil {
		return nil, err
	}

	return billProcessing, nil
}

func (billProcessing *BillProcessing) Update(
	status string) *internal_error.InternalError {

	if status != "" {
		status, err := GetBillProcessingStatusByName(status)
		if err != nil {
			return err
		}
		billProcessing.Status = status
	}

	billProcessing.UpdatedAt = time.Now()

	if err := billProcessing.Validate(); err != nil {
		return err
	}

	return nil
}

func (billProcessing *BillProcessing) Validate() *internal_error.InternalError {

	return nil
}

type BillProcessingRepositoryInterface interface {
	CreateBillProcessing(
		ctx context.Context,
		billProcessingEntity *BillProcessing) *internal_error.InternalError
	UpdateBillProcessing(
		ctx context.Context,
		billProcessingEntity *BillProcessing) *internal_error.InternalError
	FindBillProcessingById(
		ctx context.Context, billProcessingId string) (*BillProcessing, *internal_error.InternalError)
}
