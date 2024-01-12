package api

import (
	"hotel-project/store"
	"hotel-project/validation"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomHandler struct {
	store         *store.Store
	RoomValidator *validation.XValidator
}

func NewRoomHandler(store *store.Store, rsValidator *validation.XValidator) *RoomHandler {
	return &RoomHandler{
		store:         store,
		RoomValidator: rsValidator,
	}
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"error": err.Error()})
	}
	filter := bson.M{"hotelID": oid}
	rooms, err := h.store.Rooms.GetRooms(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}
