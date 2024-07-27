package email_value_source_usecase

import (
	"context"
	"time"

	"github.com/regismartiny/lembrador-contas-go/internal/entity/email_value_source_entity"
	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
)

func FindEmailValueSourceUseCase(emailValueSourceRepository email_value_source_entity.EmailValueSourceRepositoryInterface) EmailValueSourceUseCaseInterface {
	return &EmailValueSourceUseCase{
		emailValueSourceRepository,
	}
}

type EmailValueSourceOutputDTO struct {
	Id            string    `json:"id"`
	Address       string    `json:"address"`
	Subject       string    `json:"subject"`
	DataExtractor string    `json:"dataExtractor"`
	CreatedAt     time.Time `json:"createdAt" time_format:"2006-01-02 15:04:05"`
	UpdatedAt     time.Time `json:"updatedAt" time_format:"2006-01-02 15:04:05"`
}

func (u *EmailValueSourceUseCase) FindEmailValueSourceById(
	ctx context.Context, id string) (*EmailValueSourceOutputDTO, *internal_error.InternalError) {
	emailValueSourceEntity, err := u.emailValueSourceRepository.FindEmailValueSourceById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &EmailValueSourceOutputDTO{
		Id:            emailValueSourceEntity.Id,
		Address:       emailValueSourceEntity.Address,
		Subject:       emailValueSourceEntity.Subject,
		DataExtractor: email_value_source_entity.EmailValueSourceDataExtractor(emailValueSourceEntity.DataExtractor).Name(),
		CreatedAt:     emailValueSourceEntity.CreatedAt,
		UpdatedAt:     emailValueSourceEntity.UpdatedAt,
	}, nil
}

func (u *EmailValueSourceUseCase) FindEmailValueSources(
	ctx context.Context,
	address, subject string) ([]EmailValueSourceOutputDTO, *internal_error.InternalError) {
	emailValueSourceEntities, err := u.emailValueSourceRepository.FindEmailValueSources(
		ctx, address, subject)
	if err != nil {
		return nil, err
	}

	var emailValueSourceOutputs []EmailValueSourceOutputDTO
	for _, value := range emailValueSourceEntities {
		emailValueSourceOutputs = append(emailValueSourceOutputs, EmailValueSourceOutputDTO{
			Id:            value.Id,
			Address:       value.Address,
			Subject:       value.Subject,
			DataExtractor: email_value_source_entity.EmailValueSourceDataExtractor(value.DataExtractor).Name(),
			CreatedAt:     value.CreatedAt,
			UpdatedAt:     value.UpdatedAt,
		})
	}

	return emailValueSourceOutputs, nil
}
