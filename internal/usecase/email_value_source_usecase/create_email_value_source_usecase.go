package email_value_source_usecase

import (
	"context"

	"github.com/regismartiny/lembrador-contas-go/internal/entity/email_value_source_entity"
	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
)

type EmailValueSourceInputDTO struct {
	Address       string `json:"address" binding:"required,min=5"`
	Subject       string `json:"subject" binding:"required,min=3"`
	DataExtractor string `json:"dataExtractor" binding:"required"`
}

type EmailValueSourceUseCaseInterface interface {
	CreateEmailValueSource(
		ctx context.Context,
		emailValueSourceInput EmailValueSourceInputDTO) *internal_error.InternalError
	FindEmailValueSourceById(
		ctx context.Context,
		id string) (*EmailValueSourceOutputDTO, *internal_error.InternalError)
	FindEmailValueSources(
		ctx context.Context,
		address string,
		subject string) ([]*EmailValueSourceOutputDTO, *internal_error.InternalError)
	UpdateEmailValueSource(
		ctx context.Context,
		id string,
		emailValueSourceInput UpdateEmailValueSourceInputDTO) *internal_error.InternalError
}

type EmailValueSourceUseCase struct {
	emailValueSourceRepository email_value_source_entity.EmailValueSourceRepositoryInterface
}

func NewEmailValueSourceUseCase(
	emailValueSourceRepository email_value_source_entity.EmailValueSourceRepositoryInterface) EmailValueSourceUseCaseInterface {
	return &EmailValueSourceUseCase{
		emailValueSourceRepository: emailValueSourceRepository,
	}
}

func (u *EmailValueSourceUseCase) CreateEmailValueSource(
	ctx context.Context,
	emailValueSourceInput EmailValueSourceInputDTO) *internal_error.InternalError {

	emailValueSource, err := email_value_source_entity.
		CreateEmailValueSource(emailValueSourceInput.Address, emailValueSourceInput.Subject, emailValueSourceInput.DataExtractor)
	if err != nil {
		return err
	}

	if err := u.emailValueSourceRepository.CreateEmailValueSource(ctx, emailValueSource); err != nil {
		return err
	}

	return nil
}
