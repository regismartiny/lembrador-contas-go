package email_value_source_usecase

import (
	"context"

	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
)

type UpdateEmailValueSourceInputDTO struct {
	Address       string `json:"address" binding:"min=5"`
	Subject       string `json:"subject" binding:"min=3"`
	DataExtractor string `json:"dataExtractor"`
}

func (u *EmailValueSourceUseCase) UpdateEmailValueSource(
	ctx context.Context,
	id string,
	emailValueSourceInput UpdateEmailValueSourceInputDTO) *internal_error.InternalError {

	emailValueSourceEntity, err := u.emailValueSourceRepository.FindEmailValueSourceById(ctx, id)
	if err != nil {
		return err
	}

	if err := emailValueSourceEntity.
		Update(emailValueSourceInput.Address, emailValueSourceInput.Subject, emailValueSourceInput.DataExtractor); err != nil {
		return err
	}

	if err := u.emailValueSourceRepository.UpdateEmailValueSource(ctx, emailValueSourceEntity); err != nil {
		return err
	}

	return nil
}
