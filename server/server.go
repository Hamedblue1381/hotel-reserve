package server

import (
	"context"
	"hotel-project/api"
	"hotel-project/store"
	"hotel-project/util"
	"hotel-project/validation"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	envConfig    *util.Config
	userHandler  *api.UserHandler
	hotelHandler *api.HotelHandler
	roomHandler  *api.RoomHandler
}

func NewServer(envConfig *util.Config) *Server {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(envConfig.DBDriver))
	if err != nil {
		log.Fatal(err)
	}

	var (
		//TODO: these must happen in different layers
		collUser  = client.Database(envConfig.DBName).Collection(envConfig.UsersCollection)
		collHotel = client.Database(envConfig.DBName).Collection(envConfig.HotelsCollection)
		collRoom  = client.Database(envConfig.DBName).Collection(envConfig.RoomsCollection)

		userStore  = store.NewMongoUserStore(client, collUser)
		hotelStore = store.NewMongoHotelStore(client, collHotel)
		roomStore  = store.NewMongoRoomStore(client, collRoom, hotelStore)

		store = &store.Store{
			Users:  userStore,
			Hotels: hotelStore,
			Rooms:  roomStore,
		}

		userHandler  = api.NewUserHandler(store, validation.NewXValidator(validator.New()))
		hotelHandler = api.NewHotelHandler(store, validation.NewXValidator(validator.New()))
		roomHandler  = api.NewRoomHandler(store, validation.NewXValidator(validator.New()))
	)
	return &Server{
		envConfig:    envConfig,
		userHandler:  userHandler,
		hotelHandler: hotelHandler,
		roomHandler:  roomHandler,
	}
}

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.Status(fiber.StatusBadRequest).JSON(validation.GlobalErrorHandlerResp{
			Success: false,
			Message: err.Error(),
		})
	},
}

func (s *Server) RunServer() *fiber.App {

	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")

	apiv1.Get("/user", s.userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", s.userHandler.HandleGetUser)
	apiv1.Post("/user", s.userHandler.HandlePostUser)
	apiv1.Put("/user/:id", s.userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", s.userHandler.HandleDeleteUser)

	apiv1.Get("/hotels", s.hotelHandler.HandleGetHotels)
	apiv1.Get("/hotels/:id", s.hotelHandler.HandleGetHotel)
	apiv1.Post("/hotels", s.hotelHandler.HandlePostHotel)
	// apiv1.Put("/hotels/:id", s.hotelHandler.HandlePutHotel)
	// apiv1.Delete("/hotels/:id", s.hotelHandler.HandleDeleteHotel)

	apiv1.Get("hotel/:id/rooms", s.roomHandler.HandleGetRooms)

	return app
}
