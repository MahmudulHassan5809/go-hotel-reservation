package db

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	DBUri = "mongodb://localhost:27017"
	DBName = "hotel-reservation"
	TestDBName = "hotel-reservation-test"
)

func ToObjectID(id string) primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	return oid
}