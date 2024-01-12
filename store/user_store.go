package store

import (
	"context"
	"fmt"
	"hotel-project/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore interface {
	GetUserByID(context.Context, primitive.ObjectID) (*models.User, error)
	GetUsers(context.Context) ([]*models.User, error)
	InsertUser(context.Context, *models.User) (*models.User, error)
	UpdateUser(context.Context, primitive.ObjectID, *models.UpdateUserParams) (*models.User, error)
	DeleteUser(context.Context, primitive.ObjectID) error
	Dropper
}
type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client, coll *mongo.Collection) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll:   coll,
	}
}

func (s *MongoUserStore) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	if err := s.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*models.User, error) {
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []*models.User
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (s *MongoUserStore) InsertUser(ctx context.Context, user *models.User) (*models.User, error) {
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, id primitive.ObjectID, params *models.UpdateUserParams) (*models.User, error) {
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
	return s.GetUserByID(ctx, id)
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
func (s *MongoUserStore) Drop(context context.Context) error {
	fmt.Println("--- dropping user collection ---")
	return s.coll.Drop(context)
}

type PostgresUserStore struct{}
