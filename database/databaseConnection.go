package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// returns a mongo client
func DBinstance() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error Loading the .env file..")
	}

	MongoDb := os.Getenv("MONGODB_URL")

	// instantiating a mongo client
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
	if err != nil{
		log.Fatal(err)
	}

	// setting a timeout 
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
		if err != nil {
		log.Fatal(err)
	}

	fmt.Println("connecting to mongoDB..")

	// DBinstance function always returns the &client
	return client
}

// calling the function and capturing it in Client..
var Client *mongo.Client = DBinstance()

// Accessing a database collection; the function returns the collection
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("cluster0").Collection(collectionName)
	return collection
}