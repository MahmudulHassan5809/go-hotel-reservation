package db

import "go.mongodb.org/mongo-driver/bson/primitive"

const DBName = "hotel-reservation"

func ToObjectID(id string) primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	return oid
}