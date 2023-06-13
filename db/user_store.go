package db

import (
	"context"
	"hotel-reservation/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


const userColl = "users"

type UserStore interface {
	GetUseById(context.Context, string) (*types.User, error)
}

type MongoUserStore struct {
	client *mongo.Client
	coll  *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore{
	coll := client.Database(DBName).Collection(userColl)
	return &MongoUserStore{
		client: client,
		coll: coll,

	}
}

func (ms *MongoUserStore) GetUseById(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user types.User
	if err := ms.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}