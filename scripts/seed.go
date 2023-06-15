package main

import (
	"context"
	"fmt"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main()  {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBUri))
	if err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client, db.DBName)
	hotel := types.Hotel{
		Name: "Bellucia",
		Location: "France",
	}

	room := types.Room{
		Type: types.SingleRoomType,
		BasePrice: 99.9,
	}
	fmt.Println("Seeding the database")
}