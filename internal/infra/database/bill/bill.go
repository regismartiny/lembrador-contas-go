package bill

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/regismartiny/lembrador-contas-go/configuration/logger"
	"github.com/regismartiny/lembrador-contas-go/internal/entity/bill_entity"
	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BillEntityMongo struct {
	Id              string                      `bson:"_id"`
	UserId          string                      `bson:"user_id"`
	Name            string                      `bson:"name"`
	Company         string                      `bson:"company"`
	ValueSourceType bill_entity.ValueSourceType `bson:"value_source_type"`
	ValueSourceId   string                      `bson:"value_source_id"`
	DueDay          uint8                       `bson:"due_day"`
	Status          bill_entity.BillStatus      `bson:"status"`
	CreatedAt       int64                       `bson:"created_at"`
	UpdatedAt       int64                       `bson:"updated_at"`
}

type BillRepository struct {
	Collection *mongo.Collection
}

func NewBillRepository(ctx context.Context, database *mongo.Database) *BillRepository {
	coll := database.Collection("bills")

	createBillNameUniqueIndex(ctx, coll)

	return &BillRepository{
		Collection: coll,
	}
}

func createBillNameUniqueIndex(ctx context.Context, coll *mongo.Collection) {
	_, err := coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"name": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		logger.Error("Error creating bill name unique index", err)
	}
}

func (ur *BillRepository) CreateBill(
	ctx context.Context,
	billEntity *bill_entity.Bill) *internal_error.InternalError {

	BillEntityMongo := &BillEntityMongo{
		Id:              billEntity.Id,
		UserId:          billEntity.UserId,
		Name:            billEntity.Name,
		Company:         billEntity.Company,
		ValueSourceType: billEntity.ValueSourceType,
		ValueSourceId:   billEntity.ValueSourceId,
		DueDay:          billEntity.DueDay,
		Status:          billEntity.Status,
		CreatedAt:       billEntity.CreatedAt.Unix(),
		UpdatedAt:       billEntity.UpdatedAt.Unix(),
	}

	if _, err := ur.Collection.InsertOne(ctx, BillEntityMongo); err != nil {
		logger.Error("Error trying to insert bill", err)
		return internal_error.NewInternalServerError("Error trying to insert bill")
	}

	return nil
}

func (ur *BillRepository) FindBillById(
	ctx context.Context, billId string) (*bill_entity.Bill, *internal_error.InternalError) {
	filter := bson.M{"_id": billId}

	var billEntityMongo BillEntityMongo
	err := ur.Collection.FindOne(ctx, filter).Decode(&billEntityMongo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Error(fmt.Sprintf("Bill not found with this id = %s", billId), err)
			return nil, internal_error.NewNotFoundError(
				fmt.Sprintf("Bill not found with this id = %s", billId))
		}

		logger.Error("Error trying to find bill by billId", err)
		return nil, internal_error.NewInternalServerError("Error trying to find bill by billId")
	}

	billEntity := &bill_entity.Bill{
		Id:              billEntityMongo.Id,
		UserId:          billEntityMongo.UserId,
		Name:            billEntityMongo.Name,
		Company:         billEntityMongo.Company,
		ValueSourceType: billEntityMongo.ValueSourceType,
		ValueSourceId:   billEntityMongo.ValueSourceId,
		DueDay:          billEntityMongo.DueDay,
		Status:          billEntityMongo.Status,
		CreatedAt:       time.Unix(billEntityMongo.CreatedAt, 0),
		UpdatedAt:       time.Unix(billEntityMongo.UpdatedAt, 0),
	}

	return billEntity, nil
}

func (repo *BillRepository) FindBills(
	ctx context.Context,
	status bill_entity.BillStatus,
	userId string,
	name, company string) ([]*bill_entity.Bill, *internal_error.InternalError) {
	filter := bson.M{}

	if status != 0 {
		filter["status"] = status
	}

	if userId != "" {
		filter["user_id"] = userId
	}

	if name != "" {
		filter["name"] = primitive.Regex{Pattern: name, Options: "i"}
	}

	if company != "" {
		filter["company"] = primitive.Regex{Pattern: company, Options: "i"}
	}

	cursor, err := repo.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error("Error finding bills", err)
		return nil, internal_error.NewInternalServerError("Error finding bills")
	}
	defer cursor.Close(ctx)

	var billsMongo []BillEntityMongo
	if err := cursor.All(ctx, &billsMongo); err != nil {
		logger.Error("Error decoding bills", err)
		return nil, internal_error.NewInternalServerError("Error decoding bills")
	}

	billsEntity := make([]*bill_entity.Bill, len(billsMongo))
	for i, bill := range billsMongo {
		billsEntity[i] = &bill_entity.Bill{
			Id:              bill.Id,
			UserId:          bill.UserId,
			Name:            bill.Name,
			Company:         bill.Company,
			ValueSourceType: bill.ValueSourceType,
			ValueSourceId:   bill.ValueSourceId,
			DueDay:          bill.DueDay,
			Status:          bill.Status,
			CreatedAt:       time.Unix(bill.CreatedAt, 0),
			UpdatedAt:       time.Unix(bill.UpdatedAt, 0),
		}
	}

	return billsEntity, nil
}
