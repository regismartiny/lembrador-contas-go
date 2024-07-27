package email_value_source

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/regismartiny/lembrador-contas-go/configuration/logger"
	"github.com/regismartiny/lembrador-contas-go/internal/entity/email_value_source_entity"
	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EmailValueSourceEntityMongo struct {
	Id            string                                                  `bson:"_id"`
	Address       string                                                  `bson:"address"`
	Subject       string                                                  `bson:"subject"`
	DataExtractor email_value_source_entity.EmailValueSourceDataExtractor `bson:"data_extractor"`
	CreatedAt     int64                                                   `bson:"created_at"`
	UpdatedAt     int64                                                   `bson:"updated_at"`
}

type EmailValueSourceRepository struct {
	Collection *mongo.Collection
}

func NewEmailValueSourceRepository(ctx context.Context, database *mongo.Database) *EmailValueSourceRepository {
	coll := database.Collection("emailValueSources")

	createEmailValueSourceNameUniqueIndex(ctx, coll)

	return &EmailValueSourceRepository{
		Collection: coll,
	}
}

func createEmailValueSourceNameUniqueIndex(ctx context.Context, coll *mongo.Collection) {
	_, err := coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"name": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		logger.Error("Error creating emailValueSource name unique index", err)
	}
}

func (ur *EmailValueSourceRepository) CreateEmailValueSource(
	ctx context.Context,
	emailValueSourceEntity *email_value_source_entity.EmailValueSource) *internal_error.InternalError {

	EmailValueSourceEntityMongo := &EmailValueSourceEntityMongo{
		Id:            emailValueSourceEntity.Id,
		Address:       emailValueSourceEntity.Address,
		Subject:       emailValueSourceEntity.Subject,
		DataExtractor: emailValueSourceEntity.DataExtractor,
		CreatedAt:     emailValueSourceEntity.CreatedAt.Unix(),
		UpdatedAt:     emailValueSourceEntity.UpdatedAt.Unix(),
	}

	if _, err := ur.Collection.InsertOne(ctx, EmailValueSourceEntityMongo); err != nil {
		logger.Error("Error trying to insert emailValueSource", err)
		return internal_error.NewInternalServerError("Error trying to insert emailValueSource")
	}

	return nil
}

func (ur *EmailValueSourceRepository) UpdateEmailValueSource(
	ctx context.Context,
	emailValueSourceEntity *email_value_source_entity.EmailValueSource) *internal_error.InternalError {

	filter := bson.M{"_id": emailValueSourceEntity.Id}

	EmailValueSourceEntityMongo := &EmailValueSourceEntityMongo{
		Id:            emailValueSourceEntity.Id,
		Address:       emailValueSourceEntity.Address,
		Subject:       emailValueSourceEntity.Subject,
		DataExtractor: emailValueSourceEntity.DataExtractor,
		CreatedAt:     emailValueSourceEntity.CreatedAt.Unix(),
		UpdatedAt:     emailValueSourceEntity.UpdatedAt.Unix(),
	}

	_, err := ur.Collection.UpdateOne(ctx, filter, bson.M{"$set": EmailValueSourceEntityMongo})
	if err != nil {
		logger.Error("Error trying to update emailValueSource", err)
		return internal_error.NewInternalServerError("Error trying to update emailValueSource")
	}

	return nil
}

func (ur *EmailValueSourceRepository) FindEmailValueSourceById(
	ctx context.Context, emailValueSourceId string) (*email_value_source_entity.EmailValueSource, *internal_error.InternalError) {
	filter := bson.M{"_id": emailValueSourceId}

	var emailValueSourceEntityMongo EmailValueSourceEntityMongo
	err := ur.Collection.FindOne(ctx, filter).Decode(&emailValueSourceEntityMongo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Error(fmt.Sprintf("EmailValueSource not found with this id = %s", emailValueSourceId), err)
			return nil, internal_error.NewNotFoundError(
				fmt.Sprintf("EmailValueSource not found with this id = %s", emailValueSourceId))
		}

		logger.Error("Error trying to find emailValueSource by emailValueSourceId", err)
		return nil, internal_error.NewInternalServerError("Error trying to find emailValueSource by emailValueSourceId")
	}

	emailValueSourceEntity := &email_value_source_entity.EmailValueSource{
		Id:            emailValueSourceEntityMongo.Id,
		Address:       emailValueSourceEntityMongo.Address,
		Subject:       emailValueSourceEntityMongo.Subject,
		DataExtractor: emailValueSourceEntityMongo.DataExtractor,
		CreatedAt:     time.Unix(emailValueSourceEntityMongo.CreatedAt, 0),
		UpdatedAt:     time.Unix(emailValueSourceEntityMongo.UpdatedAt, 0),
	}

	return emailValueSourceEntity, nil
}

func (repo *EmailValueSourceRepository) FindEmailValueSources(
	ctx context.Context,
	address, subject string) ([]email_value_source_entity.EmailValueSource, *internal_error.InternalError) {
	filter := bson.M{}

	if address != "" {
		filter["address"] = primitive.Regex{Pattern: address, Options: "i"}
	}

	if subject != "" {
		filter["subject"] = primitive.Regex{Pattern: address, Options: "i"}
	}

	cursor, err := repo.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error("Error finding emailValueSources", err)
		return nil, internal_error.NewInternalServerError("Error finding emailValueSources")
	}
	defer cursor.Close(ctx)

	var emailValueSourcesMongo []EmailValueSourceEntityMongo
	if err := cursor.All(ctx, &emailValueSourcesMongo); err != nil {
		logger.Error("Error decoding emailValueSources", err)
		return nil, internal_error.NewInternalServerError("Error decoding emailValueSources")
	}

	var emailValueSourcesEntity []email_value_source_entity.EmailValueSource
	for _, emailValueSource := range emailValueSourcesMongo {
		emailValueSourcesEntity = append(emailValueSourcesEntity, email_value_source_entity.EmailValueSource{
			Id:            emailValueSource.Id,
			Address:       emailValueSource.Address,
			Subject:       emailValueSource.Subject,
			DataExtractor: emailValueSource.DataExtractor,
			CreatedAt:     time.Unix(emailValueSource.CreatedAt, 0),
			UpdatedAt:     time.Unix(emailValueSource.UpdatedAt, 0),
		})
	}

	return emailValueSourcesEntity, nil
}
