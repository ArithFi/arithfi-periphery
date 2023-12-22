package configs

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var (
	DB *mongo.Client
)

// init function to initialize the MongoDB client
func init() {
	DB = connectDB()
}

// connectDB function to connect to MongoDB
func connectDB() *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(EnvMongoURI()))
	if err != nil {
		log.Fatal(err)
	}

	// Set timeout for the client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	return client
}

// GetCollection function to get database collections
func GetCollection(collectionName string) *mongo.Collection {
	collection := DB.Database("periphery").Collection(collectionName)
	return collection
}
