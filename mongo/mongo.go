package mongo

import (
	"context"
	"fmt"
	"manager/config"
	"manager/debugger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func connect() *mongo.Collection {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s/?directConnection=true&serverSelectionTimeoutMS=2000&authSource=test&appName=mongosh+1.6.0&ssl=false", config.Config.Mongo_ip))
	client, err := mongo.Connect(context.TODO(), clientOptions)
	debugger.CheckError("Connect", err)

	coll := client.Database(config.Config.Mongo_database).Collection(config.Config.Mongo_collection)
	return coll
}

func CreatePerson(email string) string {
	coll := connect()
	insert, err := coll.InsertOne(context.Background(), bson.D{{Key: "email", Value: email}})
	debugger.CheckError("Insert One", err)
	return insert.InsertedID.(primitive.ObjectID).Hex()
}

func CheckHash(hash string) string {
	coll := connect()
	id, err := primitive.ObjectIDFromHex(hash)
	debugger.CheckError("Object ID from hex", err)
	cursor := coll.FindOne(context.Background(), bson.D{{Key: "_id", Value: id}}, options.FindOne())
	debugger.CheckError("Find One", cursor.Err())
	var email struct {
		Id    primitive.ObjectID `json:"_id"`
		Email string             `json:"email"`
	}
	debugger.CheckError("Decode", cursor.Decode(&email))
	return email.Email
}

func DeleteDocument(email string) error {
	return connect().FindOneAndDelete(context.Background(), bson.D{{Key: "email", Value: email}}).Err()
}
