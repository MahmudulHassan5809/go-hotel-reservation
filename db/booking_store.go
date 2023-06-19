package db

import (
	"context"
	"hotel-reservation/types"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingStore interface {
	InsertBooking(context.Context, *types.Booking) (*types.Booking, error)
	GetBookings(context.Context, bson.M) ([]*types.Booking, error)
	GetBookingByID(context.Context, string) (*types.Booking, error)
	UpdateBooking(context.Context, string, bson.M) error
}

type MongoBookingStore struct {
	client *mongo.Client
	coll  *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore{
	dbname := os.Getenv(MongoDBNameEnvName)
	coll := client.Database(dbname).Collection("bookings")
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


func (ms *MongoBookingStore) GetBookingByID(ctx context.Context, id string) (*types.Booking, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var booking types.Booking
	if err := ms.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&booking); err != nil {
		return nil, err
	}
	return &booking, nil
}


func (ms *MongoBookingStore) UpdateBooking(ctx context.Context, id string, update bson.M) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	m := bson.M{"$set": update}
	_, err = ms.coll.UpdateByID(ctx, oid, m)
	return err
}