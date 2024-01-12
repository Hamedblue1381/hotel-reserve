package store

import (
	"context"
	"fmt"
	"hotel-project/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelStore interface {
	InsertHotel(context context.Context, hotel *models.Hotel) (*models.Hotel, error)
	PushRoom(context context.Context, filter, update bson.M) error
	UpdateHotel(context context.Context, id primitive.ObjectID, params *models.UpdateHotelParams) (*models.Hotel, error)
	DestroyHotel(context context.Context, id primitive.ObjectID) error
	GetHotelByID(context context.Context, id primitive.ObjectID) (*models.Hotel, error)
	GetHotels(context context.Context, filter bson.M) ([]*models.Hotel, error)

	Dropper
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client, coll *mongo.Collection) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   coll,
	}
}
func (s *MongoHotelStore) InsertHotel(ctx context.Context, hotel *models.Hotel) (*models.Hotel, error) {
	resp, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = resp.InsertedID.(primitive.ObjectID)
	return hotel, nil
}
func (s *MongoHotelStore) UpdateHotel(ctx context.Context, id primitive.ObjectID, params *models.UpdateHotelParams) (*models.Hotel, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	update, err := models.ToDoc(params)
	if err != nil {
		return nil, err
	}
	resp, err := s.coll.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		return nil, err
	}
	//TODO: cleaner method for this response
	_ = resp
	return s.GetHotelByID(ctx, id)
}
func (s *MongoHotelStore) PushRoom(ctx context.Context, filter, update bson.M) error {
	_, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
func (s *MongoHotelStore) DestroyHotel(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
func (s *MongoHotelStore) GetHotelByID(ctx context.Context, id primitive.ObjectID) (*models.Hotel, error) {
	var hotel models.Hotel
	err := s.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&hotel)
	if err != nil {
		return nil, err
	}
	return &hotel, nil
}
func (s *MongoHotelStore) GetHotels(ctx context.Context, filter bson.M) ([]*models.Hotel, error) {
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var hotels []*models.Hotel
	if err := cur.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}

func (s *MongoHotelStore) Drop(context context.Context) error {
	fmt.Println("--- dropping hotel collection ---")
	return s.coll.Drop(context)
}
