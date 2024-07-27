package table_value_source_usecase

import (
	"context"
	"time"

	"github.com/regismartiny/lembrador-contas-go/internal/entity/table_value_source_entity"
	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
)

func FindTableValueSourceUseCase(tableValueSourceRepository table_value_source_entity.TableValueSourceRepositoryInterface) TableValueSourceUseCaseInterface {
	return &TableValueSourceUseCase{
		tableValueSourceRepository,
	}
}

type TableValueSourceOutputDTO struct {
	Id        string                                           `json:"id"`
	Name      string                                           `json:"name"`
	Data      []table_value_source_entity.TableValueSourceData `json:"data"`
	Status    string                                           `json:"status"`
	CreatedAt time.Time                                        `json:"createdAt" time_format:"2006-01-02 15:04:05"`
	UpdatedAt time.Time                                        `json:"updatedAt" time_format:"2006-01-02 15:04:05"`
}

func (u *TableValueSourceUseCase) FindTableValueSourceById(
	ctx context.Context, id string) (*TableValueSourceOutputDTO, *internal_error.InternalError) {
	tableValueSourceEntity, err := u.tableValueSourceRepository.FindTableValueSourceById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &TableValueSourceOutputDTO{
		Id:        tableValueSourceEntity.Id,
		Name:      tableValueSourceEntity.Name,
		Data:      tableValueSourceEntity.Data,
		Status:    table_value_source_entity.TableValueSourceStatus(tableValueSourceEntity.Status).Name(),
		CreatedAt: tableValueSourceEntity.CreatedAt,
		UpdatedAt: tableValueSourceEntity.UpdatedAt,
	}, nil
}

func (u *TableValueSourceUseCase) FindTableValueSources(
	ctx context.Context,
	status table_value_source_entity.TableValueSourceStatus,
	name string) ([]TableValueSourceOutputDTO, *internal_error.InternalError) {
	tableValueSourceEntities, err := u.tableValueSourceRepository.FindTableValueSources(
		ctx, status, name)
	if err != nil {
		return nil, err
	}

	var tableValueSourceOutputs []TableValueSourceOutputDTO
	for _, value := range tableValueSourceEntities {
		tableValueSourceOutputs = append(tableValueSourceOutputs, TableValueSourceOutputDTO{
			Id:        value.Id,
			Name:      value.Name,
			Data:      value.Data,
			Status:    table_value_source_entity.TableValueSourceStatus(value.Status).Name(),
			CreatedAt: value.CreatedAt,
			UpdatedAt: value.UpdatedAt,
		})
	}

	return tableValueSourceOutputs, nil
}
