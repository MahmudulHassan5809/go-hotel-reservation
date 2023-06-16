package main

import (
	"context"
	"flag"
	"hotel-reservation/api"
	"hotel-reservation/db"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



var config = fiber.Config{
    ErrorHandler: func(ctx *fiber.Ctx, err error) error {
        return ctx.JSON(map[string]string{"error": err.Error()})
    },
}

func main() {
    listenAddr := flag.String("listenAddr", ":8000", "The listen address of the api server")
    flag.Parse()

    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBUri))
	if err != nil {
		log.Fatal(err)
	}
    // handlers initialization
    var (
        
        hotelStore = db.NewMongoHotelStore(client)
        roomStore = db.NewMongoRoomStore(client, hotelStore)
        userStore = db.NewMongoUserStore(client)
        store = &db.Store{
            Hotel: hotelStore,
            Room: roomStore,
            User: userStore,
        }
        userHandler = api.NewUserHandler(userStore)
        hotelHandler = api.NewHotelHandler(*store)

        app = fiber.New(config)
        apiV1 = app.Group("/api/v1")
    )

    
    // User Handlers
    apiV1.Get("/users", userHandler.HandleGetUsers)
    apiV1.Get("/users/:id", userHandler.HandleGetUser)
    apiV1.Post("/users", userHandler.HandlePostUser)
    apiV1.Delete("/users/:id", userHandler.HandleDeleteUser)
    apiV1.Put("/users/:id", userHandler.HandlePutUser)

    // Hotel Handlers
    apiV1.Get("/hotels", hotelHandler.HandleGetHotels)
    apiV1.Get("/hotels/:id", hotelHandler.HandleGetHotel)
    apiV1.Get("/hotels/:id/rooms", hotelHandler.HandleGetRooms)
    
    
    app.Listen(*listenAddr)
    
}
