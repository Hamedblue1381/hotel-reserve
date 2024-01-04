package api

import (
	"context"
	"fmt"
	"hotel-project/types"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleHomePage(c *fiber.Ctx) error {
	return c.JSON("hello world!")
}
func HandleGetUsers(c *fiber.Ctx) error {
	user := types.User{
		FirstName: "Hamed",
		LastName:  "Malek",
	}
	return c.JSON(user)
}
func HandleGetUser(c *fiber.Ctx) error {
	return c.JSON("hello world!")
}
func HandlePostUser(db *mongo.Database, user types.User) error {

	coll := db.Collection(os.Getenv("USER_COLLECTION"))
	res, err := coll.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
	return nil
}
