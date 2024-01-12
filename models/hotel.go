package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Hotel struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
	Rating   int                  `bson:"rating" json:"rating"`
}
type CreateHotelParams struct {
	Name     string `bson:"name" json:"name" validate:"required"`
	Location string `bson:"location" json:"location" validate:"required"`
	Rooms    []Room `bson:"rooms" json:"rooms" validate:"required"`
}
type UpdateHotelParams struct {
	Name     string `bson:"name,omitempty" json:"name" validate:""`
	Location string `bson:"location,omitempty" json:"location" validate:""`
	Rooms    []Room `bson:"rooms,omitempty" json:"rooms" validate:""`
}

func NewHotelFromParams(params CreateHotelParams) (*Hotel, error) {
	return &Hotel{
		Name:     params.Name,
		Location: params.Location,
		Rooms:    []primitive.ObjectID{},
	}, nil
}
