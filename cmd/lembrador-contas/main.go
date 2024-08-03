package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/regismartiny/lembrador-contas-go/configuration/database/mongodb"
	"github.com/regismartiny/lembrador-contas-go/internal/email_service"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/api/web/controller/bill_controller"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/api/web/controller/bill_processing_controller"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/api/web/controller/email_value_source_controller"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/api/web/controller/invoice_controller"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/api/web/controller/table_value_source_controller"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/api/web/controller/user_controller"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/database/bill"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/database/bill_processing"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/database/email_value_source"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/database/invoice"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/database/table_value_source"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/database/user"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/gmail_service"
	"github.com/regismartiny/lembrador-contas-go/internal/usecase/bill_processing_usecase"
	"github.com/regismartiny/lembrador-contas-go/internal/usecase/bill_usecase"
	"github.com/regismartiny/lembrador-contas-go/internal/usecase/email_value_source_usecase"
	"github.com/regismartiny/lembrador-contas-go/internal/usecase/invoice_usecase"
	"github.com/regismartiny/lembrador-contas-go/internal/usecase/table_value_source_usecase"
	"github.com/regismartiny/lembrador-contas-go/internal/usecase/user_usecase"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/api/gmail/v1"
)

func main() {

	ctx := context.Background()

	if err := godotenv.Load("cmd/lembrador-contas/.env"); err != nil {
		log.Fatal("Error trying to load env variables")
		return
	}

	log.Println("Establising connection with database...")
	databaseConnection, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	log.Println("Creating Gmail service...")
	gmailService, err := gmail_service.NewGmailService()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	router := gin.Default()

	deps := initDependencies(ctx, databaseConnection, gmailService)

	router.GET("/user", deps.userController.FindUsers)
	router.GET("/user/:id", deps.userController.FindUserById)
	router.POST("/user", deps.userController.CreateUser)
	router.GET("/bill", deps.billController.FindBills)
	router.GET("/bill/:id", deps.billController.FindBillById)
	router.POST("/bill", deps.billController.CreateBill)
	router.GET("/invoice", deps.invoiceControler.FindInvoices)
	router.GET("/invoice/:id", deps.invoiceControler.FindInvoiceById)
	router.POST("/invoice", deps.invoiceControler.CreateInvoice)
	router.GET("/table-value-source", deps.tableValueSourceController.FindTableValueSources)
	router.GET("/table-value-source/:id", deps.tableValueSourceController.FindTableValueSourceById)
	router.POST("/table-value-source", deps.tableValueSourceController.CreateTableValueSource)
	router.PUT("/table-value-source/:id", deps.tableValueSourceController.UpdateTableValueSource)
	router.GET("/email-value-source", deps.emailValueSourceController.FindEmailValueSources)
	router.GET("/email-value-source/:id", deps.emailValueSourceController.FindEmailValueSourceById)
	router.POST("/email-value-source", deps.emailValueSourceController.CreateEmailValueSource)
	router.PUT("/email-value-source/:id", deps.emailValueSourceController.UpdateEmailValueSource)
	router.POST("/bill-processing/start", deps.billProcessingController.StartBillProcessing)
	router.GET("/bill-processing/status/:id", deps.billProcessingController.GetBillProcessingStatus)
	router.GET("/bill-processing", deps.billProcessingController.FindBillProcessings)

	router.Run(":8080")
}

func initDependencies(ctx context.Context, database *mongo.Database, gmailService *gmail.Service) *Dependencies {

	userRepository := user.NewUserRepository(ctx, database)
	userUseCase := user_usecase.NewUserUseCase(userRepository)
	userController := user_controller.NewUserController(userUseCase)

	billRepository := bill.NewBillRepository(ctx, database)
	billUseCase := bill_usecase.NewBillUseCase(billRepository)
	billController := bill_controller.NewBillController(billUseCase)

	invoiceRepository := invoice.NewInvoiceRepository(ctx, database)
	invoiceUseCase := invoice_usecase.NewInvoiceUseCase(invoiceRepository)
	invoiceController := invoice_controller.NewInvoiceController(invoiceUseCase)

	tableValueSourceRepository := table_value_source.NewTableValueSourceRepository(ctx, database)
	tableValueSourceUseCase := table_value_source_usecase.NewTableValueSourceUseCase(tableValueSourceRepository)
	tableValueSourceController := table_value_source_controller.NewTableValueSourceController(tableValueSourceUseCase)

	emailValueSourceRepository := email_value_source.NewEmailValueSourceRepository(ctx, database)
	emailValueSourceUseCase := email_value_source_usecase.NewEmailValueSourceUseCase(emailValueSourceRepository)
	emailValueSourceController := email_value_source_controller.NewEmailValueSourceController(emailValueSourceUseCase)

	billProcessingRepository := bill_processing.NewBillProcessingRepository(ctx, database)
	emailService := email_service.NewGmailEmailService(gmailService)
	billProcessingUseCase := bill_processing_usecase.NewBillProcessingUseCase(billProcessingRepository, billRepository, tableValueSourceRepository,
		emailValueSourceRepository, invoiceRepository, emailService)
	billProcessingController := bill_processing_controller.NewBillProcessingController(billProcessingUseCase)

	return &Dependencies{
		userController, billController, invoiceController, tableValueSourceController, emailValueSourceController, billProcessingController,
	}
}

type Dependencies struct {
	userController             *user_controller.UserController
	billController             *bill_controller.BillController
	invoiceControler           *invoice_controller.InvoiceController
	tableValueSourceController *table_value_source_controller.TableValueSourceController
	emailValueSourceController *email_value_source_controller.EmailValueSourceController
	billProcessingController   *bill_processing_controller.BillProcessingController
}
