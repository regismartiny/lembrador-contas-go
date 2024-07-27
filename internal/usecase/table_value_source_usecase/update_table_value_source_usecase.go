package table_value_source_usecase

import (
	"context"

	"github.com/regismartiny/lembrador-contas-go/internal/entity/table_value_source_entity"
	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
)

type UpdateTableValueSourceInputDTO struct {
	Name   string                    `json:"name" binding:"min=3"`
	Data   []TableValueSourceDataDTO `json:"data"`
	Status string                    `json:"status"`
}

func (u *TableValueSourceUseCase) UpdateTableValueSource(
	ctx context.Context,
	id string,
	tableValueSourceInput UpdateTableValueSourceInputDTO) *internal_error.InternalError {

	tableValueSourceEntity, err := u.tableValueSourceRepository.FindTableValueSourceById(ctx, id)
	if err != nil {
		return err
	}

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

	if err := tableValueSourceEntity.Update(tableValueSourceInput.Name, data, tableValueSourceInput.Status); err != nil {
		return err
	}

	if err := u.tableValueSourceRepository.UpdateTableValueSource(ctx, tableValueSourceEntity); err != nil {
		return err
	}

	return nil
}
