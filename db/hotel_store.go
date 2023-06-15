package db

import (
	"context"
	"hotel-reservation/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelStore interface {
	InsertHotel(context.Context, string) (*types.Hotel, error)
}

type MongoHotelStore struct {
	client *mongo.Client
	coll  *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client, dbname string) *MongoHotelStore{
	coll := client.Database(dbname).Collection("hotels")
	return &MongoHotelStore{
		client: client,
		coll: coll,

	}
}

func (ms *MongoUserStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := ms.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = res.InsertedID.(primitive.ObjectID)
	return hotel, nil
}