package table_value_source

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/regismartiny/lembrador-contas-go/configuration/logger"
	"github.com/regismartiny/lembrador-contas-go/internal/entity/table_value_source_entity"
	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TableValueSourceEntityMongo struct {
	Id        string                                           `bson:"_id"`
	Name      string                                           `bson:"name"`
	Data      []table_value_source_entity.TableValueSourceData `bson:"company"`
	Status    table_value_source_entity.TableValueSourceStatus `bson:"status"`
	CreatedAt int64                                            `bson:"created_at"`
	UpdatedAt int64                                            `bson:"updated_at"`
}

type TableValueSourceRepository struct {
	Collection *mongo.Collection
}

func NewTableValueSourceRepository(ctx context.Context, database *mongo.Database) *TableValueSourceRepository {
	coll := database.Collection("tableValueSources")

	createTableValueSourceNameUniqueIndex(ctx, coll)

	return &TableValueSourceRepository{
		Collection: coll,
	}
}

func createTableValueSourceNameUniqueIndex(ctx context.Context, coll *mongo.Collection) {
	_, err := coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"name": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		logger.Error("Error creating tableValueSource name unique index", err)
	}
}

func (ur *TableValueSourceRepository) CreateTableValueSource(
	ctx context.Context,
	tableValueSourceEntity *table_value_source_entity.TableValueSource) *internal_error.InternalError {

	TableValueSourceEntityMongo := &TableValueSourceEntityMongo{
		Id:        tableValueSourceEntity.Id,
		Name:      tableValueSourceEntity.Name,
		Data:      tableValueSourceEntity.Data,
		Status:    tableValueSourceEntity.Status,
		CreatedAt: tableValueSourceEntity.CreatedAt.Unix(),
		UpdatedAt: tableValueSourceEntity.UpdatedAt.Unix(),
	}

	if _, err := ur.Collection.InsertOne(ctx, TableValueSourceEntityMongo); err != nil {
		logger.Error("Error trying to insert tableValueSource", err)
		return internal_error.NewInternalServerError("Error trying to insert tableValueSource")
	}

	return nil
}

func (ur *TableValueSourceRepository) UpdateTableValueSource(
	ctx context.Context,
	tableValueSourceEntity *table_value_source_entity.TableValueSource) *internal_error.InternalError {

	filter := bson.M{"_id": tableValueSourceEntity.Id}

	TableValueSourceEntityMongo := &TableValueSourceEntityMongo{
		Id:        tableValueSourceEntity.Id,
		Name:      tableValueSourceEntity.Name,
		Data:      tableValueSourceEntity.Data,
		Status:    tableValueSourceEntity.Status,
		CreatedAt: tableValueSourceEntity.CreatedAt.Unix(),
		UpdatedAt: tableValueSourceEntity.UpdatedAt.Unix(),
	}

	_, err := ur.Collection.UpdateOne(ctx, filter, bson.M{"$set": TableValueSourceEntityMongo})
	if err != nil {
		logger.Error("Error trying to update tableValueSource", err)
		return internal_error.NewInternalServerError("Error trying to update tableValueSource")
	}

	return nil
}

func (ur *TableValueSourceRepository) FindTableValueSourceById(
	ctx context.Context, tableValueSourceId string) (*table_value_source_entity.TableValueSource, *internal_error.InternalError) {
	filter := bson.M{"_id": tableValueSourceId}

	var tableValueSourceEntityMongo TableValueSourceEntityMongo
	err := ur.Collection.FindOne(ctx, filter).Decode(&tableValueSourceEntityMongo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Error(fmt.Sprintf("TableValueSource not found with this id = %s", tableValueSourceId), err)
			return nil, internal_error.NewNotFoundError(
				fmt.Sprintf("TableValueSource not found with this id = %s", tableValueSourceId))
		}

		logger.Error("Error trying to find tableValueSource by tableValueSourceId", err)
		return nil, internal_error.NewInternalServerError("Error trying to find tableValueSource by tableValueSourceId")
	}

	tableValueSourceEntity := &table_value_source_entity.TableValueSource{
		Id:        tableValueSourceEntityMongo.Id,
		Name:      tableValueSourceEntityMongo.Name,
		Data:      tableValueSourceEntityMongo.Data,
		Status:    tableValueSourceEntityMongo.Status,
		CreatedAt: time.Unix(tableValueSourceEntityMongo.CreatedAt, 0),
		UpdatedAt: time.Unix(tableValueSourceEntityMongo.UpdatedAt, 0),
	}

	return tableValueSourceEntity, nil
}

func (repo *TableValueSourceRepository) FindTableValueSources(
	ctx context.Context,
	status table_value_source_entity.TableValueSourceStatus,
	name string) ([]table_value_source_entity.TableValueSource, *internal_error.InternalError) {
	filter := bson.M{}

	if status != 0 {
		filter["status"] = status
	}

	if name != "" {
		filter["name"] = primitive.Regex{Pattern: name, Options: "i"}
	}

	cursor, err := repo.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error("Error finding tableValueSources", err)
		return nil, internal_error.NewInternalServerError("Error finding tableValueSources")
	}
	defer cursor.Close(ctx)

	var tableValueSourcesMongo []TableValueSourceEntityMongo
	if err := cursor.All(ctx, &tableValueSourcesMongo); err != nil {
		logger.Error("Error decoding tableValueSources", err)
		return nil, internal_error.NewInternalServerError("Error decoding tableValueSources")
	}

	var tableValueSourcesEntity []table_value_source_entity.TableValueSource
	for _, tableValueSource := range tableValueSourcesMongo {
		tableValueSourcesEntity = append(tableValueSourcesEntity, table_value_source_entity.TableValueSource{
			Id:        tableValueSource.Id,
			Name:      tableValueSource.Name,
			Data:      tableValueSource.Data,
			Status:    tableValueSource.Status,
			CreatedAt: time.Unix(tableValueSource.CreatedAt, 0),
			UpdatedAt: time.Unix(tableValueSource.UpdatedAt, 0),
		})
	}

	return tableValueSourcesEntity, nil
}
