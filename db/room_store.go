package db

import (
	"context"
	"hotel-reservation/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll  *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore{
	coll := client.Database(DBName).Collection("rooms")
	return &MongoRoomStore{
		client: client,
		coll: coll,
		HotelStore: hotelStore,

	}
}

func (ms *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := ms.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = res.InsertedID.(primitive.ObjectID)

	// update the hotel with this room id
	filter := bson.M{"_id": room.HotelID}
	update := bson.M{"$push": bson.M{"rooms": room.ID}}

	if err := ms.HotelStore.UpdateHotel(ctx, filter, update); err != nil {
		return nil, err
	}

	return room, nil
}