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
	GetUsers(context.Context) ([]*types.User, error)
	CreateUser(context.Context, *types.User) (*types.User, error)
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


func (ms *MongoUserStore) CreateUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := ms.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}


func (ms *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	cur, err := ms.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []*types.User
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
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