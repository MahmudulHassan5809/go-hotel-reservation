package api

import (
	"context"
	"hotel-reservation/db"
	"log"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


const testMongoUri = "mongodb://localhost:27017"

type testDb struct {
	db.UserStore
}


func (tdb *testDb) tearDown(t *testing.T) {
	// if err := tdp.UserStore.Drop(); err != nil {
	// 	t.Fatal(err)
	// }
}


func setup(t *testing.T) *testDb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testMongoUri))
	if err != nil {
		log.Fatal(err)
	}

	return &testDb{
		UserStore: db.NewMongoUserStore(client),
	}
}

func TestCreateUser(t *testing.T){
	testDB := setup()
}