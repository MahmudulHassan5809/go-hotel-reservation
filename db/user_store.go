package db

import (
	"context"
	"fmt"
	"hotel-reservation/types"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


const userColl = "users"


type Map map[string]any

type Dropper interface {
	Drop(context.Context) error
}



type UserStore interface {
	Dropper
	GetUseById(context.Context, string) (*types.User, error)
	GetUserByEmail(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	CreateUser(context.Context, *types.User) (*types.User, error)
	DeleteUser(context.Context, string) error
	UpdateUser(context.Context, Map,  types.UpdateUserParams) error
}

type MongoUserStore struct {
	client *mongo.Client
	coll  *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore{
	dbname := os.Getenv(MongoDBNameEnvName)
	coll := client.Database(dbname).Collection(userColl)
	return &MongoUserStore{
		client: client,
		coll: coll,

	}
}

func (ms *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("--- dropping user collection")
	return ms.coll.Drop(ctx)
}

func (ms *MongoUserStore) UpdateUser(ctx context.Context, filter Map, params types.UpdateUserParams) error {
	oid, err := primitive.ObjectIDFromHex(filter["_id"].(string))
	if err != nil {
		return err
	}
	filter["_id"] = oid
	update := bson.M{"$set": params.ToBSON()}
	_, err = ms.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err 
	}
	return nil
}

func (ms *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	// TODO: Maybe it is a good idea to handle if we did not delete any use
	_, err = ms.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	return nil
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


func (ms *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	var user types.User
	if err := ms.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}