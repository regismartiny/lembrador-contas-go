package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/regismartiny/lembrador-contas-go/configuration/logger"
	"github.com/regismartiny/lembrador-contas-go/internal/entity/user_entity"
	"github.com/regismartiny/lembrador-contas-go/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserEntityMongo struct {
	Id        string                 `bson:"_id"`
	Name      string                 `bson:"name"`
	Email     string                 `bson:"email"`
	Status    user_entity.UserStatus `bson:"status"`
	CreatedAt int64                  `bson:"created_at"`
	UpdatedAt int64                  `bson:"updated_at"`
}

type UserRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(ctx context.Context, database *mongo.Database) *UserRepository {
	coll := database.Collection("users")

	createUserEmailUniqueIndex(ctx, coll)

	return &UserRepository{
		Collection: coll,
	}
}

func createUserEmailUniqueIndex(ctx context.Context, coll *mongo.Collection) {
	_, err := coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		logger.Error("Error creating user email unique index", err)
	}
}

func (ur *UserRepository) FindUserById(
	ctx context.Context, userId string) (*user_entity.User, *internal_error.InternalError) {
	filter := bson.M{"_id": userId}

	var userEntityMongo UserEntityMongo
	err := ur.Collection.FindOne(ctx, filter).Decode(&userEntityMongo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Error(fmt.Sprintf("User not found with this id = %s", userId), err)
			return nil, internal_error.NewNotFoundError(
				fmt.Sprintf("User not found with this id = %s", userId))
		}

		logger.Error("Error trying to find user by userId", err)
		return nil, internal_error.NewInternalServerError("Error trying to find user by userId")
	}

	userEntity := &user_entity.User{
		Id:        userEntityMongo.Id,
		Name:      userEntityMongo.Name,
		Email:     userEntityMongo.Email,
		Status:    userEntityMongo.Status,
		CreatedAt: time.Unix(userEntityMongo.CreatedAt, 0),
		UpdatedAt: time.Unix(userEntityMongo.UpdatedAt, 0),
	}

	return userEntity, nil
}

func (ur *UserRepository) CreateUser(
	ctx context.Context,
	userEntity *user_entity.User) *internal_error.InternalError {

	userEntityMongo := &UserEntityMongo{
		Id:        userEntity.Id,
		Name:      userEntity.Name,
		Email:     userEntity.Email,
		Status:    userEntity.Status,
		CreatedAt: userEntity.CreatedAt.Unix(),
		UpdatedAt: userEntity.UpdatedAt.Unix(),
	}

	if _, err := ur.Collection.InsertOne(ctx, userEntityMongo); err != nil {
		logger.Error("Error trying to insert user", err)
		return internal_error.NewInternalServerError("Error trying to insert user")
	}

	return nil
}
