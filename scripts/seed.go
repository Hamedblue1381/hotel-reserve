package main

import (
	"context"
	"hotel-project/models"
	"hotel-project/store"
	"hotel-project/util"
	"log"
	"math/rand"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	roomStore  store.RoomStore
	hotelStore store.HotelStore
	ctx        = context.Background()
	wg         = sync.WaitGroup{}
)

func seedHotel(name, location string, rooms []models.Room, rating int) {
	defer wg.Done()
	hotel := models.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
		// CreatedAt: time.Now(),
		// UpdatedAt: time.Now(),
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	//TODO:remove this random
	numRooms := rand.Intn(4)

	for _, room := range rooms[:numRooms] {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func main() {
	rooms := []models.Room{}
	rooms = append(rooms, *NewRoom(models.SingleRoomType, 99.9, 2, 120))
	rooms = append(rooms, *NewRoom(models.DoubleRoomType, 199.9, 4, 250))
	rooms = append(rooms, *NewRoom(models.DeluxeRoomType, 299.9, 6, 350))
	rooms = append(rooms, *NewRoom(models.SeaSideRoomType, 399.9, 4, 450))
	wg.Add(4)
	go seedHotel("Palladium", "Iran", rooms, 2)
	go seedHotel("Morvarid", "Iran", rooms, 3)
	go seedHotel("Cozy Corner", "Iran", rooms, 5)
	go seedHotel("Panjeh", "Iran", rooms, 5)

	wg.Wait()
	log.Println("Seeding the database...")
}
func NewRoom(typeOfRoom models.RoomType, basePrice float64, beds int, price float64) *models.Room {
	return &models.Room{
		Type:      typeOfRoom,
		BasePrice: basePrice,
		Beds:      beds,
		Price:     price,
	}
}
func init() {
	envConfig, err := util.LoadConfig("./")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(envConfig.DBDriver))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(envConfig.DBName).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelColl := client.Database(envConfig.DBName).Collection(envConfig.HotelsCollection)
	roomColl := client.Database(envConfig.DBName).Collection(envConfig.RoomsCollection)
	hotelStore = store.NewMongoHotelStore(client, hotelColl)
	roomStore = store.NewMongoRoomStore(client, roomColl, hotelStore)
}
