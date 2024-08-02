package email_data_extractor

import (
	"time"

	"github.com/regismartiny/lembrador-contas-go/internal/email_service"
	"github.com/regismartiny/lembrador-contas-go/internal/entity/email_value_source_entity"
	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
)

type EmailDataExtractorInterface interface {
	Extract(request EmailDataExtractorRequest) (*EmailDataExtractorResponse, *internal_error.InternalError)
}

type EmailDataExtractorRequest struct {
	Subject   string
	Address   string
	StartDate time.Time
	EndDate   time.Time
}

type EmailDataExtractorResponse struct {
	Amount float64
}

func NewEmailDataExtractor(emailService email_service.EmailServiceInterface, dataExtractor email_value_source_entity.EmailValueSourceDataExtractor) EmailDataExtractorInterface {
	switch dataExtractor {
	case email_value_source_entity.CPFL_EMAIL_EXTRACTOR:
		return NewCpflEmailDataExtractor(emailService)

	default:
		return nil
	}
}
