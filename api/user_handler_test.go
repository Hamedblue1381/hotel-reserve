package api

import (
	"bytes"
	"context"
	"encoding/json"
	"hotel-project/models"
	"hotel-project/store"
	"hotel-project/util"
	"hotel-project/validation"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	store *store.Store
}

func (testdb *testdb) teardown(t *testing.T) error {
	if err := testdb.store.Users.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
	return nil
}

func setup(t *testing.T) *testdb {
	envConfig, err := util.LoadConfig("./..")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(envConfig.TestDBDriver))
	if err != nil {
		log.Fatal(err)
	}
	coll := client.Database(envConfig.TestDBName).Collection(envConfig.TestUsersCollection)
	return &testdb{
		store: &store.Store{
			Users: store.NewMongoUserStore(client, coll),
		},
	}
}
func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	UserHandler := NewUserHandler(tdb.store, validation.NewXValidator(validator.New()))
	app.Post("/api/v1/user", UserHandler.HandlePostUser)
	params := models.CreateUserParams{
		Email:     "test@example.com",
		FirstName: "James",
		LastName:  "Bond",
		Password:  "123456789",
	}
	b, err := json.Marshal(params)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/api/v1/user", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var user models.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, params.FirstName, user.FirstName)
	assert.Equal(t, params.LastName, user.LastName)
	assert.Equal(t, params.Email, user.Email)
	assert.Empty(t, user.EncryptedPassword)
	assert.NotEmpty(t, user.ID)
}
