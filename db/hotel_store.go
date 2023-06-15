package db

import (
	"context"
	"hotel-reservation/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelStore interface {
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	UpdateHotel(context.Context, bson.M, bson.M) error
	GetHotels(context.Context, bson.M) ([]*types.Hotel, error)
}

type MongoHotelStore struct {
	client *mongo.Client
	coll  *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore{
	coll := client.Database(DBName).Collection("hotels")
	return &MongoHotelStore{
		client: client,
		coll: coll,

	}
}

func (ms *MongoHotelStore) UpdateHotel(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := ms.coll.UpdateOne(ctx, filter, update)
	return err
}

func (ms *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := ms.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = res.InsertedID.(primitive.ObjectID)
	return hotel, nil
}


func (ms *MongoHotelStore) GetHotels(ctx context.Context, filter bson.M) ([]*types.Hotel, error) {
	resp, err := ms.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var hotels []*types.Hotel 
	if err := resp.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}