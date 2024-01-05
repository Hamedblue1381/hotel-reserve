package api

import (
	"errors"
	"fmt"
	"hotel-project/models"
	"hotel-project/store"
	"hotel-project/validation"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore     store.UserStore
	userValidator *validation.XValidator
}

func NewUserHandler(userStore store.UserStore, userValidator *validation.XValidator) *UserHandler {
	return &UserHandler{
		userStore:     userStore,
		userValidator: userValidator,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)
	user, err := h.userStore.GetUserByID(c.Context(), id)
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
	users, err := h.userStore.GetUsers(c.Context())
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
	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	userId := c.Params("id")
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
	updatedUser, err := h.userStore.UpdateUser(c.Context(), userId, &params)
	if err != nil {
		return err
	}

	return c.JSON(updatedUser)

}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userId := c.Params("id")
	if err := h.userStore.DeleteUser(c.Context(), userId); err != nil {
		return err
	}

	return c.Status(fiber.StatusNoContent).JSON(nil)
}
