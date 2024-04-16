package api

import (
	"bytes"
	"context"
	"encoding/json"
	"hotel-project/models"
	"hotel-project/store"
	"hotel-project/validation"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func insertMockUser(t *testing.T, userStore store.UserStore) *models.User {
	user, err := models.NewUserFromParams(models.CreateUserParams{
		FirstName: "james",
		LastName:  "Bond",
		Email:     "hunter@gmail.com",
		Password:  "password12345",
	})
	if err != nil {
		t.Fatal(err)
	}
	createdUser, err := userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}
	return createdUser
}
func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	createdUser := insertMockUser(t, tdb.store.Users)
	defer tdb.teardown(t)

	app := fiber.New()
	AuthHandler := NewAuthHandler(tdb.store, validation.NewXValidator(validator.New()))
	app.Post("/auth", AuthHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "hunter@gmail.com",
		Password: "password12345",
	}
	b, err := json.Marshal(params)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	var authResponse AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	if err != nil {
		t.Fatal(err)
	}
	// Set the encrypted password to null (response doesnt contain password)
	createdUser.EncryptedPassword = ""
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, createdUser, authResponse.User)
	assert.NotEmpty(t, authResponse.Token)
}
func TestAuthenticateFailure(t *testing.T) {
	tdb := setup(t)
	createdUser := insertMockUser(t, tdb.store.Users)
	_ = createdUser
	defer tdb.teardown(t)

	app := fiber.New()
	AuthHandler := NewAuthHandler(tdb.store, validation.NewXValidator(validator.New()))
	app.Post("/auth", AuthHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "hunter@gmail.com",
		Password: "invalid",
	}
	b, err := json.Marshal(params)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
