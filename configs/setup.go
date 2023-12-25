package configs

import (
	"context"
	"crypto/tls"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var (
	DB    *mongo.Client
	CACHE *redis.Client
)

// init function to initialize the MongoDB client and the Redis client
func init() {
	DB = connectDB()
	CACHE = connectCache()
}

func connectCache() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     EnvRedisURI(),
		Password: "",
		DB:       0,
	})
	return rdb
}

// connectDB function to connect to MongoDB
func connectDB() *mongo.Client {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	tslConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	opts := options.Client().ApplyURI(EnvMongoURI()).SetServerAPIOptions(serverAPI).SetTLSConfig(tslConfig)
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
		log.Fatal("Failed to ping MongoDB", err)
	}
	log.Println("Connected to MongoDB")
	return client
}

// GetCollection function to get database collections
func GetCollection(collectionName string) *mongo.Collection {
	collection := DB.Database("periphery").Collection(collectionName)
	return collection
}
