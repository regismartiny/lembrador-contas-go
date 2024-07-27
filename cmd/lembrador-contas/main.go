package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/regismartiny/lembrador-contas-go/configuration/database/mongodb"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/api/web/controller/bill_controller"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/api/web/controller/invoice_controller"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/api/web/controller/table_value_source_controller"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/api/web/controller/user_controller"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/database/bill"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/database/invoice"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/database/table_value_source"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/database/user"
	"github.com/regismartiny/lembrador-contas-go/internal/usecase/bill_usecase"
	"github.com/regismartiny/lembrador-contas-go/internal/usecase/invoice_usecase"
	"github.com/regismartiny/lembrador-contas-go/internal/usecase/table_value_source_usecase"
	"github.com/regismartiny/lembrador-contas-go/internal/usecase/user_usecase"
	"go.mongodb.org/mongo-driver/mongo"
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

	router := gin.Default()

	userController, billController, invoiceControler, tableValueSourceController := initDependencies(ctx, databaseConnection)

	router.GET("/user", userController.FindUsers)
	router.GET("/user/:id", userController.FindUserById)
	router.POST("/user", userController.CreateUser)
	router.GET("/bill", billController.FindBills)
	router.GET("/bill/:id", billController.FindBillById)
	router.POST("/bill", billController.CreateBill)
	router.GET("/invoice", invoiceControler.FindInvoices)
	router.GET("/invoice/:id", invoiceControler.FindInvoiceById)
	router.POST("/invoice", invoiceControler.CreateInvoice)
	router.GET("/table_value_source", tableValueSourceController.FindTableValueSources)
	router.GET("/table_value_source/:id", tableValueSourceController.FindTableValueSourceById)
	router.POST("/table_value_source", tableValueSourceController.CreateTableValueSource)
	router.PUT("/table_value_source/:id", tableValueSourceController.UpdateTableValueSource)

	router.Run(":8080")
}

func initDependencies(ctx context.Context, database *mongo.Database) (
	userController *user_controller.UserController,
	billController *bill_controller.BillController,
	invoiceControler *invoice_controller.InvoiceController,
	tableValueSourceController *table_value_source_controller.TableValueSourceController) {

	userRepository := user.NewUserRepository(ctx, database)

	userController = user_controller.NewUserController(
		user_usecase.NewUserUseCase(userRepository))

	billRepository := bill.NewBillRepository(ctx, database)

	billController = bill_controller.NewBillController(
		bill_usecase.NewBillUseCase(billRepository))

	invoiceRepository := invoice.NewInvoiceRepository(ctx, database)

	invoiceController := invoice_controller.NewInvoiceController(
		invoice_usecase.NewInvoiceUseCase(invoiceRepository))

	tableValueSourceRepository := table_value_source.NewTableValueSourceRepository(ctx, database)

	tableValueSourceController = table_value_source_controller.NewTableValueSourceController(
		table_value_source_usecase.NewTableValueSourceUseCase(tableValueSourceRepository))

	return userController, billController, invoiceController, tableValueSourceController
}
