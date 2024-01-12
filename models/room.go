package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Room struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type      RoomType           `bson:"type" json:"type"`
	BasePrice float64            `bson:"basePrice" json:"basePrice"`
	Beds      int                `bson:"beds" json:"beds"`
	Price     float64            `bson:"price" json:"price"`
	HotelID   primitive.ObjectID `bson:"hotelID,omitempty" json:"hotelID,omitempty"`
}
type UpdateRoomParams struct {
	Type      RoomType           `bson:"type" json:"type" validation:""`
	BasePrice float64            `bson:"basePrice" json:"basePrice" validation:""`
	Beds      int                `bson:"beds" json:"beds" validation:""`
	Price     float64            `bson:"price" json:"price" validation:""`
	HotelID   primitive.ObjectID `bson:"hotelID,omitempty" json:"hotelID,omitempty"`
}
type CreateRoomParams struct {
	Type      RoomType           `bson:"type" json:"type" validation:"required"`
	BasePrice float64            `bson:"basePrice" json:"basePrice" validation:"required"`
	Beds      int                `bson:"beds" json:"beds" validation:"required"`
	Price     float64            `bson:"price" json:"price" validation:"required"`
	HotelID   primitive.ObjectID `bson:"hotelID,omitempty" json:"hotelID,omitempty"`
}

type RoomType int

const (
	_ RoomType = iota
	SingleRoomType
	DoubleRoomType
	SeaSideRoomType
	DeluxeRoomType
)
