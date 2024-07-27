package bill_entity

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
)

type Bill struct {
	Id              string
	Name            string
	Company         string
	ValueSourceType ValueSourceType
	ValueSourceId   string
	DueDay          uint8
	Status          BillStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type BillStatus uint8

const (
	Active BillStatus = iota
	Inactive
)

type ValueSourceType uint8

const (
	Table ValueSourceType = iota
	Email
	API
)

func (s BillStatus) Name() string {
	return billStatusNames[s]
}

var billStatusNames = []string{
	"active",
	"inactive",
}

func GetBillStatusByName(name string) (BillStatus, *internal_error.InternalError) {
	for k, v := range billStatusNames {
		if v == name {
			return BillStatus(k), nil
		}
	}

	return BillStatus(0), internal_error.NewBadRequestError("invalid bill status name")
}

func (s ValueSourceType) Name() string {
	return valueSourceTypeNames[s]
}

var valueSourceTypeNames = []string{
	"table",
	"email",
	"api",
}

func GetValueSourceTypeByName(name string) (ValueSourceType, *internal_error.InternalError) {
	for k, v := range valueSourceTypeNames {
		if v == name {
			return ValueSourceType(k), nil
		}
	}

	return ValueSourceType(0), internal_error.NewBadRequestError("invalid bill valueSourceType name")
}

func CreateBill(
	name string,
	company string,
	valueSourceType string,
	valueSourceId string,
	dueDay uint8,
	status string) (*Bill, *internal_error.InternalError) {

	var billStatus BillStatus

	if status == "" {
		billStatus = Active
	} else {
		status, err := GetBillStatusByName(status)
		if err != nil {
			return nil, err
		}
		billStatus = status
	}

	var billValueSourceType ValueSourceType

	if valueSourceType == "" {
		return nil, internal_error.NewBadRequestError("invalid bill value source type")
	} else {
		valueSourceType, err := GetValueSourceTypeByName(valueSourceType)
		if err != nil {
			return nil, err
		}
		billValueSourceType = valueSourceType
	}

	bill :=
		&Bill{
			Id:              uuid.New().String(),
			Name:            name,
			Company:         company,
			ValueSourceType: billValueSourceType,
			ValueSourceId:   valueSourceId,
			DueDay:          dueDay,
			Status:          billStatus,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}

	if err := bill.Validate(); err != nil {
		return nil, err
	}

	return bill, nil
}

func (bill *Bill) Validate() *internal_error.InternalError {
	if len(bill.Name) <= 3 ||
		len(bill.Company) <= 3 {
		return internal_error.NewBadRequestError("invalid bill object")
	}

	return nil
}

type BillRepositoryInterface interface {
	CreateBill(ctx context.Context, billEntity *Bill) *internal_error.InternalError
	FindBillById(ctx context.Context, billId string) (*Bill, *internal_error.InternalError)
	FindBills(
		ctx context.Context,
		status BillStatus,
		name, company string) ([]Bill, *internal_error.InternalError)
}
