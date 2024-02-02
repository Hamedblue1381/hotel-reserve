package api

import (
	"errors"
	"fmt"
	"hotel-project/middleware"
	"hotel-project/models"
	"hotel-project/store"
	"hotel-project/validation"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	store         *store.Store
	authValidator *validation.XValidator
}
type AuthResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

func NewAuthHandler(store *store.Store, authValidator *validation.XValidator) *AuthHandler {
	return &AuthHandler{
		store:         store,
		authValidator: authValidator,
	}
}

type AuthParams struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	user, err := h.store.Users.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("invalid credentials")
		}
		return err
	}
	if !models.IsValidPassword(user.EncryptedPassword, params.Password) {
		return fmt.Errorf("invalid credentials")
	}
	resp := AuthResponse{
		User:  user,
		Token: middleware.CreateTokenFromUser(user),
	}
	return c.JSON(resp)
}
