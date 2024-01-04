package main

import (
	"context"
	"hotel-project/api"
	"hotel-project/types"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("DB_DRIVER")))
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database(os.Getenv("DB_NAME"))
	user := types.User{
		FirstName: "test",
		LastName:  "test",
	}
	api.HandlePostUser(db, user)
	app := fiber.New()
	apiv1 := app.Group("/api/v1")

	apiv1.Get("/", api.HandleHomePage)
	apiv1.Get("/users", api.HandleGetUsers)
	apiv1.Get("/users/:id", api.HandleGetUser)
	// apiv1.Post("/users",)

	app.Listen(os.Getenv("LISTEN_ADDRESS"))
}
