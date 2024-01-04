package mongo

import (
	"context"
	"crypto/tls"
	"github.com/arithfi/arithfi-periphery/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var (
	MONGODB *mongo.Client
)

// init function to initialize the MongoDB client and the Redis client
func init() {
	MONGODB = connectMongoDB()
}

func connectMongoDB() *mongo.Client {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	tslConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	opts := options.Client().ApplyURI(configs.EnvMongoURI()).SetServerAPIOptions(serverAPI).SetTLSConfig(tslConfig)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB", err)
	}

	// Set timeout for the client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB: ", err)
	}
	log.Println("Connected to MongoDB")
	return client
}
