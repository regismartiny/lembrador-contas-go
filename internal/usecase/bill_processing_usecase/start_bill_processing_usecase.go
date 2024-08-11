package bill_processing_usecase

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
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

type StartBillProcessingOutputDTO struct {
	BillProcessingId string `json:"id"`
}

type BillProcessingUseCaseInterface interface {
	StartBillProcessing(
		ctx context.Context,
		period string) (StartBillProcessingOutputDTO, *internal_error.InternalError)
	GetBillProcessingStatus(
		ctx context.Context,
		billProcessingId string) (GetBillProcessingStatusOutputDTO, *internal_error.InternalError)
	FindBillProcessings(
		ctx context.Context,
		status bill_processing_entity.BillProcessingStatus) ([]*FindBillProcessingOutputDTO, *internal_error.InternalError)
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
	period string) (StartBillProcessingOutputDTO, *internal_error.InternalError) {

	if err := u.verifyNoProcessingInProgress(ctx); err != nil {
		log.Println("Error trying to start bill processing", err)
		return StartBillProcessingOutputDTO{}, err
	}

	billProcessing, err := bill_processing_entity.CreateBillProcessing("", period)
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
	log.Println("Bill processing started. Period:", billProcessing.Period)

	//Find all active bills
	activeBills, err := u.billRepository.FindBills(ctx, bill_entity.Active, "", "", "")
	if err != nil {
		log.Println("Error trying to find active bills", err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(activeBills))

	billProcessingErrors := make([]string, 0)

	for _, bill := range activeBills {

		go func() {
			defer wg.Done()

			if err := u.processBill(ctx, bill, billProcessing); err != nil {
				errString := fmt.Sprintf("Error trying to process bill %v: %v", bill.Id, err)
				log.Println(errString)
				billProcessingErrors = append(billProcessingErrors, errString)
			}
		}()
	}

	wg.Wait()

	if len(billProcessingErrors) > 0 {
		log.Println("Bill processing finished with errors")
		billProcessing.Status = bill_processing_entity.Error
	} else {
		log.Println("Bill processing finished successfully")
		billProcessing.Status = bill_processing_entity.Success
	}

	u.billProcessingRepository.UpdateBillProcessing(ctx, billProcessing)
}

func (u *BillProcessingUseCase) processBill(ctx context.Context, bill *bill_entity.Bill, billProcessing *bill_processing_entity.BillProcessing) *internal_error.InternalError {
	log.Printf("Processing bill: %v", bill)

	u.deleteUnpaidInvoices(ctx, billProcessing, bill)

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
	case bill_entity.API:
		return internal_error.NewInternalServerError("valueSourceType API not implemented yet")
	}

	return nil
}

func (u *BillProcessingUseCase) deleteUnpaidInvoices(ctx context.Context, billProcessing *bill_processing_entity.BillProcessing, bill *bill_entity.Bill) {
	log.Printf("Deleting unpaid invoices for period: %v", billProcessing.Period)

	processingPeriod, _ := time.Parse("2006-01", billProcessing.Period)
	billDueDate := time.Date(processingPeriod.Year(), processingPeriod.Month()+1, int(bill.DueDay), 0, 0, 0, 0, time.Local).Format("2006-01-02")

	log.Printf("Deleting invoices with due date: %v", billDueDate)
	invoicesDeleted, _ := u.invoiceRepository.DeleteInvoices(ctx, bill.Id, invoice_entity.Unpaid, billDueDate)

	log.Printf("Invoices deleted: %v", invoicesDeleted)
}

func (u *BillProcessingUseCase) processTableValueSource(ctx context.Context, bill *bill_entity.Bill,
	tableValueSource *table_value_source_entity.TableValueSource) *internal_error.InternalError {

	now := time.Now()
	dueDate := time.Date(now.Year(), now.Month(), int(bill.DueDay), 0, 0, 0, 0, time.Local)
	processingDate := dueDate.AddDate(0, -1, 0) // always process previous month

	log.Println("Today is", now.Format("02/01/2006"))
	log.Println("Processing Bill", bill.Name, "for period", processingDate.Month(), "/", now.Year())

	amount := 0.0

	for _, v := range tableValueSource.Data {
		if v.Period.Month == uint8(processingDate.Month()) && v.Period.Year == uint16(processingDate.Year()) {
			log.Printf("Found data for current period. Value: %2.f\n", v.Amount)
			amount = v.Amount
			break
		}
	}

	if amount == 0.0 {
		log.Println("No invoice found for current period")
		return nil
	}

	return u.createInvoice(ctx, bill, dueDate, amount)
}

func (u *BillProcessingUseCase) createInvoice(ctx context.Context, bill *bill_entity.Bill, dueDate time.Time, amount float64) *internal_error.InternalError {
	log.Println("Creating invoice")

	invoice, err := invoice_entity.CreateInvoice(
		bill.Id,
		dueDate.Format("2006-01-02"),
		amount,
		"",
	)
	if err != nil {
		return err
	}

	return u.invoiceRepository.CreateInvoice(ctx, invoice)
}

func (u *BillProcessingUseCase) processEmailValueSource(ctx context.Context, bill *bill_entity.Bill,
	emailValueSource *email_value_source_entity.EmailValueSource) *internal_error.InternalError {

	log.Println("Processing email value source. Address:", emailValueSource.Address, "Subject:", emailValueSource.Subject)

	dataExtractor := email_data_extractor.NewEmailDataExtractor(u.emailService, emailValueSource.DataExtractor)

	now := time.Now()
	dueDate := time.Date(now.Year(), now.Month(), int(bill.DueDay), 0, 0, 0, 0, time.Local)
	processingDate := dueDate.AddDate(0, -1, 0) // previous month

	startDate := time.Date(processingDate.Year(), processingDate.Month(), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, -1) //last day of month

	dataExtractorResponse, err := dataExtractor.Extract(email_data_extractor.EmailDataExtractorRequest{
		Subject:   emailValueSource.Subject,
		Address:   emailValueSource.Address,
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		return err
	}

	u.createInvoice(ctx, bill, dueDate, dataExtractorResponse.Amount)

	return nil
}
