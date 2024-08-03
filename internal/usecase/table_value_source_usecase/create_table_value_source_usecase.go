package table_value_source_usecase

import (
	"context"

	"github.com/regismartiny/lembrador-contas-go/internal/entity/table_value_source_entity"
	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
)

type TableValueSourceInputDTO struct {
	Name   string                    `json:"name" binding:"required,min=3"`
	Data   []TableValueSourceDataDTO `json:"data"`
	Status string                    `json:"status"`
}

type TableValueSourceDataDTO struct {
	Period TableValueSourceDataPeriodDTO `json:"period"`
	Amount float64                       `json:"amount"`
}

type TableValueSourceDataPeriodDTO struct {
	Month uint8  `json:"month"`
	Year  uint16 `json:"year"`
}

type TableValueSourceUseCaseInterface interface {
	CreateTableValueSource(
		ctx context.Context,
		tableValueSourceInput TableValueSourceInputDTO) *internal_error.InternalError
	FindTableValueSourceById(
		ctx context.Context,
		id string) (*TableValueSourceOutputDTO, *internal_error.InternalError)
	FindTableValueSources(
		ctx context.Context,
		status table_value_source_entity.TableValueSourceStatus,
		name string) ([]*TableValueSourceOutputDTO, *internal_error.InternalError)
	UpdateTableValueSource(
		ctx context.Context,
		id string,
		tableValueSourceInput UpdateTableValueSourceInputDTO) *internal_error.InternalError
}

type TableValueSourceUseCase struct {
	tableValueSourceRepository table_value_source_entity.TableValueSourceRepositoryInterface
}

func NewTableValueSourceUseCase(
	tableValueSourceRepository table_value_source_entity.TableValueSourceRepositoryInterface) TableValueSourceUseCaseInterface {
	return &TableValueSourceUseCase{
		tableValueSourceRepository: tableValueSourceRepository,
	}
}

func (u *TableValueSourceUseCase) CreateTableValueSource(
	ctx context.Context,
	tableValueSourceInput TableValueSourceInputDTO) *internal_error.InternalError {

	data := make([]table_value_source_entity.TableValueSourceData, 0)
	for _, v := range tableValueSourceInput.Data {
		data = append(data, table_value_source_entity.TableValueSourceData{
			Period: table_value_source_entity.TableValueSourceDataPeriod{
				Month: v.Period.Month,
				Year:  v.Period.Year,
			},
			Amount: v.Amount,
		})
	}

	tableValueSource, err := table_value_source_entity.CreateTableValueSource(tableValueSourceInput.Name, data, tableValueSourceInput.Status)
	if err != nil {
		return err
	}

	if err := u.tableValueSourceRepository.CreateTableValueSource(ctx, tableValueSource); err != nil {
		return err
	}

	return nil
}
