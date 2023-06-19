package main

import (
	"context"
	"flag"
	"hotel-reservation/api"
	"hotel-reservation/db"
	"hotel-reservation/middleware"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



var config = fiber.Config{
    ErrorHandler: func(ctx *fiber.Ctx, err error) error {
        if apiError, ok := err.(api.Error); ok {
            return ctx.Status(apiError.Code).JSON(apiError)
        }
        if apiError, ok := err.(middleware.Error); ok {
            return ctx.Status(apiError.Code).JSON(apiError)
        }
        apiError := api.NewError(http.StatusInternalServerError, err.Error())
        return ctx.Status(apiError.Code).JSON(apiError)
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
        bookingStore = db.NewMongoBookingStore(client)
        store = db.Store{
            Hotel: hotelStore,
            Room: roomStore,
            User: userStore,
            Booking: bookingStore,
        }
        userHandler = api.NewUserHandler(userStore)
        hotelHandler = api.NewHotelHandler(store)
        authHandler = api.NewAuthHandler(userStore)
        roomHandler = api.NewRoomHandler(store)
        bookingHandler = api.NewBookingHandler(store)

        app = fiber.New(config)
        apiV1NoAuth = app.Group("/api")
        apiV1 = app.Group("/api/v1", middleware.JWTAuthentication(userStore))
        admin = apiV1.Group("/admin", middleware.AdminAuth)
    )

    // Auth Handler
    apiV1NoAuth.Post("/auth", authHandler.HandleAuthenticate)

    
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

    // Rooms Handlers 
    apiV1.Get("/rooms", roomHandler.HandleGetRooms)
    apiV1.Post("/room/:id/book", roomHandler.HandleBookRoom)

    // Bookings Handlers
    apiV1.Get("/bookings/:id", bookingHandler.HandleGetBooking)

    // Admin routes
    admin.Get("/bookings", bookingHandler.HandleGetBookings)
    admin.Patch("/bookings/:id/cancelled", bookingHandler.HandleCancelBooking)
    
    app.Listen(*listenAddr)
    
}
