package table_value_source_entity

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
)

type TableValueSource struct {
	Id        string
	Name      string
	Data      []TableValueSourceData
	Status    TableValueSourceStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TableValueSourceData struct {
	Period TableValueSourceDataPeriod
	Amount float64
}

type TableValueSourceDataPeriod struct {
	Month uint8
	Year  uint16
}

type TableValueSourceStatus uint8

const (
	Active TableValueSourceStatus = iota + 1
	Inactive
)

func (s TableValueSourceStatus) Name() string {
	return tableValueSourceStatusNames[s]
}

var tableValueSourceStatusNames = []string{
	"",
	"active",
	"inactive",
}

func GetTableValueSourceStatusByName(name string) (TableValueSourceStatus, *internal_error.InternalError) {
	for k, v := range tableValueSourceStatusNames {
		if v == name {
			return TableValueSourceStatus(k), nil
		}
	}

	return TableValueSourceStatus(0), internal_error.NewBadRequestError("invalid tableValueSource status name")
}

func CreateTableValueSource(
	name string,
	data []TableValueSourceData,
	status string) (*TableValueSource, *internal_error.InternalError) {

	var tableValueSourceStatus TableValueSourceStatus

	if status == "" {
		tableValueSourceStatus = Active
	} else {
		status, err := GetTableValueSourceStatusByName(status)
		if err != nil {
			return nil, err
		}
		tableValueSourceStatus = status
	}

	if data == nil {
		data = make([]TableValueSourceData, 0)
	}

	tableValueSource :=
		&TableValueSource{
			Id:        uuid.New().String(),
			Name:      name,
			Data:      data,
			Status:    tableValueSourceStatus,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

	if err := tableValueSource.Validate(); err != nil {
		return nil, err
	}

	return tableValueSource, nil
}

func (tableValueSource *TableValueSource) Update(
	name string,
	data []TableValueSourceData,
	status string) *internal_error.InternalError {

	if name != "" {
		tableValueSource.Name = name
	}

	if data != nil {
		tableValueSource.Data = data
	}

	if status != "" {
		status, err := GetTableValueSourceStatusByName(status)
		if err != nil {
			return err
		}
		tableValueSource.Status = status
	}

	tableValueSource.UpdatedAt = time.Now()

	if err := tableValueSource.Validate(); err != nil {
		return err
	}

	return nil
}

func (tableValueSource *TableValueSource) Validate() *internal_error.InternalError {
	if len(tableValueSource.Name) < 3 {
		return internal_error.NewBadRequestError("invalid tableValueSource object")
	}
	for _, v := range tableValueSource.Data {
		if v.Period.Month > 12 || v.Period.Year > 9999 {
			return internal_error.NewBadRequestError("invalid tableValueSource data")
		}
	}

	return nil
}

type TableValueSourceRepositoryInterface interface {
	CreateTableValueSource(ctx context.Context, tableValueSourceEntity *TableValueSource) *internal_error.InternalError
	FindTableValueSourceById(ctx context.Context, tableValueSourceId string) (*TableValueSource, *internal_error.InternalError)
	FindTableValueSources(
		ctx context.Context,
		status TableValueSourceStatus,
		name string) ([]*TableValueSource, *internal_error.InternalError)
	UpdateTableValueSource(ctx context.Context, tableValueSourceEntity *TableValueSource) *internal_error.InternalError
}
