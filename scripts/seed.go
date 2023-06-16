package main

import (
	"context"
	"fmt"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	roomStore db.RoomStore
	hotelStore db.HotelStore
	ctx = context.Background()
)

func seedHotel(name string, location string, rating int) {
	hotel := types.Hotel{
		Name: name, 
		Location: location,
		Rooms: []primitive.ObjectID{},
		Rating: rating,
	}

	rooms := []types.Room{
		{
			Size: "samll",
			Price: 99.9,
		},
		{
			Size: "noraml",
			Price: 199.9,
		},
		{
			Size: "large",
			Price: 299.9,
		},
	}

	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID	
		insertedRoom , err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(insertedRoom)
	}
	fmt.Println(insertedHotel)
}  

func main()  {
	seedHotel("Bellucia", "France", 5)
	seedHotel("Westin", "Bangladesh", 4)
}

func init() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBUri))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBName).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
}