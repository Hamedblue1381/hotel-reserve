package main

import (
	"context"
	"hotel-project/api"
	"hotel-project/store"
	"hotel-project/validation"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.Status(fiber.StatusBadRequest).JSON(validation.GlobalErrorHandlerResp{
			Success: false,
			Message: err.Error(),
		})
	},
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("DB_DRIVER")))
	if err != nil {
		log.Fatal(err)
	}
	coll := client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("USER_COLLECTION"))
	userHandler := api.NewUserHandler(store.NewMongoUserStore(client, coll), validation.NewXValidator(validator.New()))

	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")

	apiv1.Get("/", api.HandleHomePage)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Post("/user", userHandler.HandlePostUser)

	app.Listen(os.Getenv("LISTEN_ADDRESS"))
}
