package bill_processing

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/regismartiny/lembrador-contas-go/configuration/logger"
	"github.com/regismartiny/lembrador-contas-go/internal/entity/bill_processing_entity"
	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BillProcessingEntityMongo struct {
	Id        string                                      `bson:"_id"`
	Status    bill_processing_entity.BillProcessingStatus `bson:"status"`
	CreatedAt int64                                       `bson:"created_at"`
	UpdatedAt int64                                       `bson:"updated_at"`
}

type BillProcessingRepository struct {
	Collection *mongo.Collection
}

func NewBillProcessingRepository(ctx context.Context, database *mongo.Database) *BillProcessingRepository {
	coll := database.Collection("billProcessings")

	return &BillProcessingRepository{
		Collection: coll,
	}
}

func (ur *BillProcessingRepository) CreateBillProcessing(
	ctx context.Context,
	billProcessingEntity *bill_processing_entity.BillProcessing) *internal_error.InternalError {

	BillProcessingEntityMongo := &BillProcessingEntityMongo{
		Id:        billProcessingEntity.Id,
		Status:    billProcessingEntity.Status,
		CreatedAt: billProcessingEntity.CreatedAt.Unix(),
		UpdatedAt: billProcessingEntity.UpdatedAt.Unix(),
	}

	if _, err := ur.Collection.InsertOne(ctx, BillProcessingEntityMongo); err != nil {
		logger.Error("Error trying to insert billProcessing", err)
		return internal_error.NewInternalServerError("Error trying to insert billProcessing")
	}

	return nil
}

func (ur *BillProcessingRepository) FindBillProcessingById(
	ctx context.Context, billProcessingId string) (*bill_processing_entity.BillProcessing, *internal_error.InternalError) {
	filter := bson.M{"_id": billProcessingId}

	var billProcessingEntityMongo BillProcessingEntityMongo
	err := ur.Collection.FindOne(ctx, filter).Decode(&billProcessingEntityMongo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Error(fmt.Sprintf("BillProcessing not found with this id = %s", billProcessingId), err)
			return nil, internal_error.NewNotFoundError(
				fmt.Sprintf("BillProcessing not found with this id = %s", billProcessingId))
		}

		logger.Error("Error trying to find billProcessing by billProcessingId", err)
		return nil, internal_error.NewInternalServerError("Error trying to find billProcessing by billProcessingId")
	}

	billProcessingEntity := &bill_processing_entity.BillProcessing{
		Id:        billProcessingEntityMongo.Id,
		Status:    billProcessingEntityMongo.Status,
		CreatedAt: time.Unix(billProcessingEntityMongo.CreatedAt, 0),
		UpdatedAt: time.Unix(billProcessingEntityMongo.UpdatedAt, 0),
	}

	return billProcessingEntity, nil
}

func (repo *BillProcessingRepository) FindBillProcessings(
	ctx context.Context,
	status bill_processing_entity.BillProcessingStatus) ([]bill_processing_entity.BillProcessing, *internal_error.InternalError) {
	filter := bson.M{}

	if status != 0 {
		filter["status"] = status
	}

	cursor, err := repo.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error("Error finding billProcessings", err)
		return nil, internal_error.NewInternalServerError("Error finding billProcessings")
	}
	defer cursor.Close(ctx)

	var billProcessingsMongo []BillProcessingEntityMongo
	if err := cursor.All(ctx, &billProcessingsMongo); err != nil {
		logger.Error("Error decoding billProcessings", err)
		return nil, internal_error.NewInternalServerError("Error decoding billProcessings")
	}

	var billProcessingsEntity []bill_processing_entity.BillProcessing
	for _, billProcessing := range billProcessingsMongo {
		billProcessingsEntity = append(billProcessingsEntity, bill_processing_entity.BillProcessing{
			Id:        billProcessing.Id,
			Status:    billProcessing.Status,
			CreatedAt: time.Unix(billProcessing.CreatedAt, 0),
			UpdatedAt: time.Unix(billProcessing.UpdatedAt, 0),
		})
	}

	return billProcessingsEntity, nil
}

func (repo *BillProcessingRepository) GetProcessingsInProgressCount(
	ctx context.Context) (int64, *internal_error.InternalError) {
	filter := bson.M{}

	filter["status"] = bill_processing_entity.Started

	count, err := repo.Collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

	return count, nil
}

func (repo *BillProcessingRepository) UpdateBillProcessing(
	ctx context.Context,
	billProcessingEntity *bill_processing_entity.BillProcessing) *internal_error.InternalError {

	filter := bson.M{"_id": billProcessingEntity.Id}

	BillProcessingEntityMongo := &BillProcessingEntityMongo{
		Id:        billProcessingEntity.Id,
		Status:    billProcessingEntity.Status,
		CreatedAt: billProcessingEntity.CreatedAt.Unix(),
		UpdatedAt: billProcessingEntity.UpdatedAt.Unix(),
	}

	_, err := repo.Collection.UpdateOne(ctx, filter, bson.M{"$set": BillProcessingEntityMongo})
	if err != nil {
		logger.Error("Error trying to update billProcessing", err)
		return internal_error.NewInternalServerError("Error trying to update billProcessing")
	}

	return nil
}
