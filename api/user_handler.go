package api

import (
	"fmt"
	"hotel-project/models"
	"hotel-project/store"
	"hotel-project/validation"
	"strings"

	"github.com/gofiber/fiber/v2"
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
		return fiber.NewError(fiber.StatusNotFound, err.Error())
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
		errMsgs := make([]string, 0)
		for _, err := range errs {
			errMsgs = append(errMsgs, fmt.Sprintf(
				"[%s]: '%v' Needs to implement '%s'",
				err.FailedField,
				err.Value,
				err.Tag,
			))
		}
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: strings.Join(errMsgs, " and "),
		}
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
