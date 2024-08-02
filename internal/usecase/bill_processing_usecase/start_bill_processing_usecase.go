package bill_processing_usecase

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/regismartiny/lembrador-contas-go/internal/data_extractor/email_data_extractor"
	"github.com/regismartiny/lembrador-contas-go/internal/email_service"
	"github.com/regismartiny/lembrador-contas-go/internal/entity/bill_entity"
	"github.com/regismartiny/lembrador-contas-go/internal/entity/bill_processing_entity"
	"github.com/regismartiny/lembrador-contas-go/internal/entity/email_value_source_entity"
	"github.com/regismartiny/lembrador-contas-go/internal/entity/invoice_entity"
	"github.com/regismartiny/lembrador-contas-go/internal/entity/table_value_source_entity"
	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
)

const (
	PROCESSING_TIMEOUT_DURATION = "PROCESSING_TIMEOUT_DURATION"
)

type BillProcessingInputDTO struct {
}

type StartBillProcessingOutputDTO struct {
	BillProcessingId string `json:"id"`
}

type BillProcessingUseCaseInterface interface {
	StartBillProcessing(
		ctx context.Context,
		billProcessingInput BillProcessingInputDTO) (StartBillProcessingOutputDTO, *internal_error.InternalError)
	GetBillProcessingStatus(
		ctx context.Context,
		billProcessingId string) (GetBillProcessingStatusOutputDTO, *internal_error.InternalError)
}

type BillProcessingUseCase struct {
	billProcessingRepository   bill_processing_entity.BillProcessingRepositoryInterface
	billRepository             bill_entity.BillRepositoryInterface
	tableValueSourceRepository table_value_source_entity.TableValueSourceRepositoryInterface
	emailValueSourceRepository email_value_source_entity.EmailValueSourceRepositoryInterface
	invoiceRepository          invoice_entity.InvoiceRepositoryInterface
	emailService               email_service.EmailServiceInterface
}

func NewBillProcessingUseCase(
	billProcessingRepository bill_processing_entity.BillProcessingRepositoryInterface,
	billRepository bill_entity.BillRepositoryInterface,
	tableValueSourceRepository table_value_source_entity.TableValueSourceRepositoryInterface,
	emailValueSourceRepository email_value_source_entity.EmailValueSourceRepositoryInterface,
	invoiceRepository invoice_entity.InvoiceRepositoryInterface,
	emailService email_service.EmailServiceInterface) BillProcessingUseCaseInterface {

	return &BillProcessingUseCase{
		billProcessingRepository:   billProcessingRepository,
		billRepository:             billRepository,
		tableValueSourceRepository: tableValueSourceRepository,
		emailValueSourceRepository: emailValueSourceRepository,
		invoiceRepository:          invoiceRepository,
		emailService:               emailService,
	}
}

func (u *BillProcessingUseCase) StartBillProcessing(
	ctx context.Context,
	billProcessingInput BillProcessingInputDTO) (StartBillProcessingOutputDTO, *internal_error.InternalError) {

	// if err := u.verifyNoProcessingInProgress(ctx); err != nil {
	// 	log.Println("Error trying to start bill processing", err)
	// 	return StartBillProcessingOutputDTO{}, err
	// }

	billProcessing, err := bill_processing_entity.CreateBillProcessing("")
	if err != nil {
		return StartBillProcessingOutputDTO{}, err
	}

	if err := u.billProcessingRepository.CreateBillProcessing(ctx, billProcessing); err != nil {
		return StartBillProcessingOutputDTO{}, err
	}

	go u.startProcessing(ctx, billProcessing)
	go u.manageProcessingTimeot(ctx, billProcessing)

	return StartBillProcessingOutputDTO{
		BillProcessingId: billProcessing.Id}, nil
}

func (u *BillProcessingUseCase) verifyNoProcessingInProgress(ctx context.Context) *internal_error.InternalError {
	count, err := u.billProcessingRepository.GetProcessingsInProgressCount(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		return internal_error.NewBadRequestError("There are already bill processings in progress")
	}
	return nil
}

func (u *BillProcessingUseCase) manageProcessingTimeot(ctx context.Context, billProcessing *bill_processing_entity.BillProcessing) {
	processingTimeout := os.Getenv(PROCESSING_TIMEOUT_DURATION)
	processingTimeoutDuration, _ := time.ParseDuration(processingTimeout)
	time.Sleep(processingTimeoutDuration)

	billProcessing, err := u.billProcessingRepository.FindBillProcessingById(ctx, billProcessing.Id)
	if err != nil {
		log.Println("Error trying to find billProcessing", err)
		return
	}

	if !billProcessing.Status.IsFinished() {
		log.Println("Bill processing timeout")
		billProcessing.Status = bill_processing_entity.Timeout

		if err := u.billProcessingRepository.UpdateBillProcessing(ctx, billProcessing); err != nil {
			log.Println("Error trying to update billProcessing", err)
		}
	}
}

func (u *BillProcessingUseCase) startProcessing(ctx context.Context, billProcessing *bill_processing_entity.BillProcessing) {
	log.Println("Bill processing started")

	activeBills, err := u.billRepository.FindBills(ctx, bill_entity.Active, "", "")
	if err != nil {
		log.Println("Error trying to find active bills", err)
		return
	}

	log.Printf("activeBills: %v", activeBills)

	for _, bill := range activeBills {

		if err := u.processBill(ctx, bill); err != nil {
			log.Println("Error trying to process bill", err)
			billProcessing.Status = bill_processing_entity.Error
			u.billProcessingRepository.UpdateBillProcessing(ctx, billProcessing)
			return
		}
	}

	log.Println("Bill processing finished successfully")

	billProcessing.Status = bill_processing_entity.Success
	u.billProcessingRepository.UpdateBillProcessing(ctx, billProcessing)
}

func (u *BillProcessingUseCase) processBill(ctx context.Context, bill bill_entity.Bill) *internal_error.InternalError {
	log.Printf("Processing bill: %v", bill)

	// Deleting all unpaid invoices of current bill
	u.invoiceRepository.DeleteInvoices(ctx, bill.Id, invoice_entity.Unpaid, time.Time{})

	valueSourceId := bill.ValueSourceId
	valueSourceType := bill.ValueSourceType

	switch valueSourceType {
	case bill_entity.Table:
		tableValueSource, err := u.tableValueSourceRepository.FindTableValueSourceById(ctx, valueSourceId)
		if err != nil {
			log.Println("Error trying to find table value source:", err)
			return err
		}
		return u.processTableValueSource(ctx, bill, tableValueSource)
	case bill_entity.Email:
		emailValueSource, err := u.emailValueSourceRepository.FindEmailValueSourceById(ctx, valueSourceId)
		if err != nil {
			log.Println("Error trying to find email value source:", err)
			return err
		}
		return u.processEmailValueSource(ctx, bill, emailValueSource)
	}

	return nil
}

func (u *BillProcessingUseCase) processTableValueSource(ctx context.Context, bill bill_entity.Bill,
	tableValueSource *table_value_source_entity.TableValueSource) *internal_error.InternalError {

	currentDate := time.Now().UTC()

	log.Println("Today is", currentDate)
	log.Println("Processing Bill", bill.Name, "for period", currentDate.Month(), "/", currentDate.Year())

	amount := 0.0

	for _, v := range tableValueSource.Data {
		if v.Period.Month == uint8(currentDate.Month()) && v.Period.Year == uint16(currentDate.Year()) {
			log.Printf("Found invoice for current month. Value: %2.f\n", v.Amount)
			amount = v.Amount
			break
		}
	}

	if amount == 0.0 {
		log.Println("No invoice found for current period")
		return nil
	}

	return u.createInvoice(ctx, bill, amount)
}

func (u *BillProcessingUseCase) createInvoice(ctx context.Context, bill bill_entity.Bill, amount float64) *internal_error.InternalError {
	log.Println("Creating invoice")

	currentDate := time.Now().UTC()

	invoiceDueDate := time.Date(currentDate.Year(), currentDate.Month(), int(bill.DueDay), 0, 0, 0, 0, time.UTC)

	invoice, err := invoice_entity.CreateInvoice(
		bill.Name,
		invoiceDueDate.Format("2006-01-02"),
		amount,
		"",
	)
	if err != nil {
		return err
	}

	return u.invoiceRepository.CreateInvoice(ctx, invoice)
}

func (u *BillProcessingUseCase) processEmailValueSource(ctx context.Context, bill bill_entity.Bill,
	emailValueSource *email_value_source_entity.EmailValueSource) *internal_error.InternalError {

	log.Println("Processing email value source. Address:", emailValueSource.Address, "Subject:", emailValueSource.Subject)

	dataExtractor := email_data_extractor.NewEmailDataExtractor(u.emailService, emailValueSource.DataExtractor)

	today := time.Now().UTC()
	year := today.Year()
	month := today.Month() - 1 // previous month //TODO: should be configurable
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, month, 31, 0, 0, 0, 0, time.UTC) //TODO: get last day of month

	dataExtractorResponse, err := dataExtractor.Extract(email_data_extractor.EmailDataExtractorRequest{
		Subject:   emailValueSource.Subject,
		Address:   emailValueSource.Address,
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		return err
	}

	u.createInvoice(ctx, bill, dataExtractorResponse.Amount)

	return nil
}
