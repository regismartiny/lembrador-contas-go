package invoice

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/regismartiny/lembrador-contas-go/configuration/logger"
	"github.com/regismartiny/lembrador-contas-go/internal/entity/invoice_entity"
	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InvoiceEntityMongo struct {
	Id        string                       `bson:"_id"`
	Name      string                       `bson:"name"`
	DueDate   int64                        `bson:"due_date"`
	Amount    int64                        `bson:"amount"`
	Status    invoice_entity.InvoiceStatus `bson:"status"`
	CreatedAt int64                        `bson:"created_at"`
	UpdatedAt int64                        `bson:"updated_at"`
}

type InvoiceRepository struct {
	Collection *mongo.Collection
}

func NewInvoiceRepository(ctx context.Context, database *mongo.Database) *InvoiceRepository {
	coll := database.Collection("invoices")

	return &InvoiceRepository{
		Collection: coll,
	}
}

func (ur *InvoiceRepository) CreateInvoice(
	ctx context.Context,
	invoiceEntity *invoice_entity.Invoice) *internal_error.InternalError {

	InvoiceEntityMongo := &InvoiceEntityMongo{
		Id:        invoiceEntity.Id,
		Name:      invoiceEntity.Name,
		DueDate:   invoiceEntity.DueDate.Unix(),
		Amount:    int64(invoiceEntity.Amount * 100),
		Status:    invoiceEntity.Status,
		CreatedAt: invoiceEntity.CreatedAt.Unix(),
		UpdatedAt: invoiceEntity.UpdatedAt.Unix(),
	}

	if _, err := ur.Collection.InsertOne(ctx, InvoiceEntityMongo); err != nil {
		logger.Error("Error trying to insert invoice", err)
		return internal_error.NewInternalServerError("Error trying to insert invoice")
	}

	return nil
}

func (ur *InvoiceRepository) FindInvoiceById(
	ctx context.Context, invoiceId string) (*invoice_entity.Invoice, *internal_error.InternalError) {
	filter := bson.M{"_id": invoiceId}

	var invoiceEntityMongo InvoiceEntityMongo
	err := ur.Collection.FindOne(ctx, filter).Decode(&invoiceEntityMongo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Error(fmt.Sprintf("Invoice not found with this id = %s", invoiceId), err)
			return nil, internal_error.NewNotFoundError(
				fmt.Sprintf("Invoice not found with this id = %s", invoiceId))
		}

		logger.Error("Error trying to find invoice by invoiceId", err)
		return nil, internal_error.NewInternalServerError("Error trying to find invoice by invoiceId")
	}

	invoiceEntity := &invoice_entity.Invoice{
		Id:        invoiceEntityMongo.Id,
		Name:      invoiceEntityMongo.Name,
		DueDate:   time.Unix(invoiceEntityMongo.DueDate, 0),
		Amount:    float64(invoiceEntityMongo.Amount / 100),
		Status:    invoiceEntityMongo.Status,
		CreatedAt: time.Unix(invoiceEntityMongo.CreatedAt, 0),
		UpdatedAt: time.Unix(invoiceEntityMongo.UpdatedAt, 0),
	}

	return invoiceEntity, nil
}

func (repo *InvoiceRepository) FindInvoices(
	ctx context.Context,
	status invoice_entity.InvoiceStatus,
	name string) ([]invoice_entity.Invoice, *internal_error.InternalError) {
	filter := bson.M{}

	if status != 0 {
		filter["status"] = status
	}

	if name != "" {
		filter["name"] = primitive.Regex{Pattern: name, Options: "i"}
	}

	cursor, err := repo.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error("Error finding invoices", err)
		return nil, internal_error.NewInternalServerError("Error finding invoices")
	}
	defer cursor.Close(ctx)

	var invoicesMongo []InvoiceEntityMongo
	if err := cursor.All(ctx, &invoicesMongo); err != nil {
		logger.Error("Error decoding invoices", err)
		return nil, internal_error.NewInternalServerError("Error decoding invoices")
	}

	var invoicesEntity []invoice_entity.Invoice
	for _, invoice := range invoicesMongo {
		invoicesEntity = append(invoicesEntity, invoice_entity.Invoice{
			Id:        invoice.Id,
			Name:      invoice.Name,
			DueDate:   time.Unix(invoice.DueDate, 0),
			Amount:    float64(invoice.Amount) / 100,
			Status:    invoice.Status,
			CreatedAt: time.Unix(invoice.CreatedAt, 0),
			UpdatedAt: time.Unix(invoice.UpdatedAt, 0),
		})
	}

	return invoicesEntity, nil
}
