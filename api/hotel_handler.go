package api

import (
	"context"
	"errors"
	"fmt"
	"hotel-project/models"
	"hotel-project/store"
	"hotel-project/validation"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelHandler struct {
	store          *store.Store
	hotelValidator *validation.XValidator
}

func NewHotelHandler(store *store.Store, hotelValidator *validation.XValidator) *HotelHandler {
	return &HotelHandler{
		store:          store,
		hotelValidator: hotelValidator,
	}
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	var (
		id = c.Params("id")
	)
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	hotel, err := h.store.Hotels.GetHotelByID(c.Context(), oid)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(map[string]string{"error": "not found"})
		}
	}
	return c.JSON(hotel)
}

type HotelQueryParams struct {
	Rooms  bool
	Rating int
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var qparams HotelQueryParams
	if err := c.QueryParser(&qparams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"error": err.Error()})
	}
	var filter bson.M
	if qparams.Rating != 0 {
		filter = bson.M{"rating": qparams.Rating}
	}
	hotels, err := h.store.Hotels.GetHotels(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}
func (h *HotelHandler) HandlePostHotel(c *fiber.Ctx) error {
	var params models.CreateHotelParams
	ctx := context.Background()

	if err := c.BodyParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(validation.GlobalErrorHandlerResp{
			Success: false,
			Message: err.Error(),
		})
	}
	errs := h.hotelValidator.Validate(params)
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
	hotel, err := models.NewHotelFromParams(params)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(validation.GlobalErrorHandlerResp{
			Success: false,
			Message: err.Error(),
		})
	}
	insertedHotel, err := h.store.Hotels.InsertHotel(c.Context(), hotel)
	if err != nil {
		return err
	}
	for i := range params.Rooms {
		params.Rooms[i].HotelID = insertedHotel.ID
		_, err := h.store.Rooms.InsertRoom(ctx, &params.Rooms[i])
		if err != nil {
			return err
		}
	}

	result, err := h.store.Hotels.GetHotelByID(ctx, insertedHotel.ID)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

func (h *HotelHandler) HandlePutHotel(c *fiber.Ctx) error {
	hotelID := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		return err
	}
	var params models.UpdateHotelParams
	if err := c.BodyParser(&params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(validation.GlobalErrorHandlerResp{
			Success: false,
			Message: err.Error(),
		})
	}
	errs := h.hotelValidator.Validate(params)
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
	updatedHotel, err := h.store.Hotels.UpdateHotel(c.Context(), oid, &params)
	if err != nil {
		return err
	}
	return c.JSON(updatedHotel)
}
func (h *HotelHandler) HandleDeleteHotel(c *fiber.Ctx) error {
	hotelID := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		return err
	}
	err = h.store.Hotels.DestroyHotel(c.Context(), oid)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusNoContent).JSON(nil)
}
