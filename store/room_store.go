package store

import (
	"context"
	"hotel-project/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	InsertRoom(context context.Context, room *models.Room) (*models.Room, error)
	UpdateRoom(context context.Context, id primitive.ObjectID, params *models.UpdateRoomParams) (*models.Room, error)
	DestroyRoom(context context.Context, id primitive.ObjectID) error
	GetRoomById(context context.Context, id primitive.ObjectID) (*models.Room, error)
	GetRooms(context context.Context, filter bson.M) ([]*models.Room, error)
}
type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, coll *mongo.Collection, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		coll:       coll,
		HotelStore: hotelStore,
	}
}
func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *models.Room) (*models.Room, error) {
	resp, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = resp.InsertedID.(primitive.ObjectID)
	filter := bson.M{"_id": room.HotelID}
	update := bson.M{"$push": bson.M{"rooms": room.ID}}
	if err = s.HotelStore.PushRoom(ctx, filter, update); err != nil {
		return nil, err
	}
	return room, nil
}
func (s *MongoRoomStore) GetRoomById(ctx context.Context, id primitive.ObjectID) (*models.Room, error) {
	var room models.Room
	err := s.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&room)
	if err != nil {
		return nil, err
	}
	return &room, nil
}
func (s *MongoRoomStore) GetRooms(ctx context.Context, filter bson.M) ([]*models.Room, error) {
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var rooms []*models.Room
	if err = cur.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}
func (s *MongoRoomStore) UpdateRoom(ctx context.Context, id primitive.ObjectID, params *models.UpdateRoomParams) (*models.Room, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	update, err := models.ToDoc(params)
	if err != nil {
		return nil, err
	}

	result, err := s.coll.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		return nil, err
	}
	//TODO: cleaner method for this response
	_ = result
	return s.GetRoomById(ctx, id)
}
func (s *MongoRoomStore) DestroyRoom(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}
