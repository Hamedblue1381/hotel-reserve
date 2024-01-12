package api

import (
	"errors"
	"fmt"
	"hotel-project/models"
	"hotel-project/store"
	"hotel-project/validation"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	store         *store.Store
	userValidator *validation.XValidator
}

func NewUserHandler(store *store.Store, userValidator *validation.XValidator) *UserHandler {
	return &UserHandler{
		store:         store,
		userValidator: userValidator,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	user, err := h.store.Users.GetUserByID(c.Context(), oid)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(map[string]string{"error": "not found"})
		}
	}
	return c.JSON(user)
}

func HandleHomePage(c *fiber.Ctx) error {
	return c.JSON("hello world!")
}
func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.store.Users.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params models.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	errs := h.userValidator.Validate(params)
	if len(errs) > 0 {
		errMap := make(map[string]string)
		for _, err := range errs {
			errMap[err.FailedField] = fmt.Sprintf(
				"'%v' Needs to implement '%s'",
				err.Value,
				err.Tag,
			)
		}
		return c.JSON(errMap)
	}

	user, err := models.NewUserFromParams(params)
	if err != nil {
		return err
	}
	insertedUser, err := h.store.Users.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	userId := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}
	var params models.UpdateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	errs := h.userValidator.Validate(params)
	if len(errs) > 0 {
		errMap := make(map[string]string)
		for _, err := range errs {
			errMap[err.FailedField] = fmt.Sprintf(
				"'%v' Needs to implement '%s'",
				err.Value,
				err.Tag,
			)
		}
		return c.JSON(errMap)
	}
	updatedUser, err := h.store.Users.UpdateUser(c.Context(), oid, &params)
	if err != nil {
		return err
	}

	return c.JSON(updatedUser)

}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userId := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}
	if err := h.store.Users.DeleteUser(c.Context(), oid); err != nil {
		return err
	}

	return c.Status(fiber.StatusNoContent).JSON(nil)
}
