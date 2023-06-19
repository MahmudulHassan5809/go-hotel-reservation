package db

import (
	"context"
	"hotel-reservation/types"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HotelStore interface {
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	UpdateHotel(context.Context, Map, Map) error
	GetHotels(context.Context, Map, *Pagination) ([]*types.Hotel, error)
	GetHotelByID(context.Context, string) (*types.Hotel, error)
}

type MongoHotelStore struct {
	client *mongo.Client
	coll  *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore{
	dbname := os.Getenv(MongoDBNameEnvName)
	coll := client.Database(dbname).Collection("hotels")
	return &MongoHotelStore{
		client: client,
		coll: coll,

	}
}

func (ms *MongoHotelStore) UpdateHotel(ctx context.Context, filter Map, update Map) error {
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


func (ms *MongoHotelStore) GetHotels(ctx context.Context, filter Map, page *Pagination) ([]*types.Hotel, error) {
	opts := options.FindOptions{}
	opts.SetSkip(int64((page.Page -1) * page.Limit))
	opts.SetLimit(int64(page.Limit))
	resp, err := ms.coll.Find(ctx, filter, &opts)
	if err != nil {
		return nil, err
	}
	var hotels []*types.Hotel 
	if err := resp.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}

func (ms *MongoHotelStore) GetHotelByID(ctx context.Context, id string) (*types.Hotel, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	
	var hotel *types.Hotel
	if err := ms.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&hotel); err != nil {
		return nil, err
	}

	return hotel, nil
}