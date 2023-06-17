package db

import (
	"context"
	"hotel-reservation/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingStore interface {
	InsertBooking(context.Context, *types.Booking) (*types.Booking, error)
	GetBookings(context.Context, bson.M) ([]*types.Booking, error)
}

type MongoBookingStore struct {
	client *mongo.Client
	coll  *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore{
	coll := client.Database(DBName).Collection("bookings")
	return &MongoBookingStore{
		client: client,
		coll: coll,

	}
}


func (ms *MongoBookingStore) GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error){
	cur, err := ms.coll.Find(ctx, filter)

	if err != nil {
		return nil, err
	}
	var bookings []*types.Booking

	if err := cur.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}

func(ms *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error){
	res, err := ms.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.ID = res.InsertedID.(primitive.ObjectID)
	return booking, nil
}