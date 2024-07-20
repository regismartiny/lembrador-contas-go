package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/regismartiny/lembrador-contas-go/configuration/database/mongodb"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/api/web/controller/user_controller"
	"github.com/regismartiny/lembrador-contas-go/internal/infra/database/user"
	"github.com/regismartiny/lembrador-contas-go/internal/usecase/user_usecase"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {

	ctx := context.Background()

	if err := godotenv.Load("cmd/lembrador-contas/.env"); err != nil {
		log.Fatal("Error trying to load env variables")
		return
	}

	databaseConnection, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	log.Println(databaseConnection.Name())

	router := gin.Default()

	userController := initDependencies(ctx, databaseConnection)

	router.GET("/user/:userId", userController.FindUserById)
	router.POST("/user", userController.CreateUser)

	router.Run(":8080")
}

func initDependencies(ctx context.Context, database *mongo.Database) (
	userController *user_controller.UserController) {

	userRepository := user.NewUserRepository(ctx, database)

	userController = user_controller.NewUserController(
		user_usecase.NewUserUseCase(userRepository))

	return userController
}
