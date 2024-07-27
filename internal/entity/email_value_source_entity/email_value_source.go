package email_value_source_entity

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
)

type EmailValueSource struct {
	Id            string
	Address       string
	Subject       string
	DataExtractor EmailValueSourceDataExtractor
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type EmailValueSourceDataExtractor uint8

const (
	CPFL_EMAIL_EXTRACTOR EmailValueSourceDataExtractor = iota + 1
)

func (s EmailValueSourceDataExtractor) Name() string {
	return emailValueSourceDataExtractorNames[s]
}

var emailValueSourceDataExtractorNames = []string{
	"",
	"CPFL_EMAIL_EXTRACTOR",
}

func GetEmailValueSourceDataExtractorByName(name string) (EmailValueSourceDataExtractor, *internal_error.InternalError) {
	for k, v := range emailValueSourceDataExtractorNames {
		if v == name {
			return EmailValueSourceDataExtractor(k), nil
		}
	}

	return EmailValueSourceDataExtractor(0), internal_error.NewBadRequestError("invalid emailValueSource dataExtractor name")
}

func CreateEmailValueSource(
	address string,
	subject string,
	dataExtractor string) (*EmailValueSource, *internal_error.InternalError) {

	valueSourceDataExtractor, err := GetEmailValueSourceDataExtractorByName(dataExtractor)
	if err != nil {
		return nil, err
	}

	emailValueSource :=
		&EmailValueSource{
			Id:            uuid.New().String(),
			Address:       address,
			Subject:       subject,
			DataExtractor: valueSourceDataExtractor,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

	if err := emailValueSource.Validate(); err != nil {
		return nil, err
	}

	return emailValueSource, nil
}

func (emailValueSource *EmailValueSource) Update(
	address string,
	subject string,
	dataExtractor string) *internal_error.InternalError {

	if address != "" {
		emailValueSource.Address = address
	}

	if subject != "" {
		emailValueSource.Subject = subject
	}

	if dataExtractor != "" {
		dataExtractor, err := GetEmailValueSourceDataExtractorByName(dataExtractor)
		if err != nil {
			return err
		}
		emailValueSource.DataExtractor = dataExtractor
	}

	emailValueSource.UpdatedAt = time.Now()

	if err := emailValueSource.Validate(); err != nil {
		return err
	}

	return nil
}

func (emailValueSource *EmailValueSource) Validate() *internal_error.InternalError {
	if len(emailValueSource.Address) < 5 {
		return internal_error.NewBadRequestError("invalid emailValueSource object")
	}

	return nil
}

type EmailValueSourceRepositoryInterface interface {
	CreateEmailValueSource(ctx context.Context, emailValueSourceEntity *EmailValueSource) *internal_error.InternalError
	FindEmailValueSourceById(ctx context.Context, emailValueSourceId string) (*EmailValueSource, *internal_error.InternalError)
	FindEmailValueSources(
		ctx context.Context,
		address string,
		subject string) ([]EmailValueSource, *internal_error.InternalError)
	UpdateEmailValueSource(ctx context.Context, emailValueSourceEntity *EmailValueSource) *internal_error.InternalError
}
