package configs

import (
	"context"
	"crypto/tls"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var (
	MONGODB *mongo.Client
	CACHE   *redis.Client
	MYSQL   *sql.DB
)

// init function to initialize the MongoDB client and the Redis client
func init() {
	MONGODB = connectMongoDB()
	CACHE = connectCache()
	MYSQL = connectMysql()
}

func connectCache() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     EnvRedisURI(),
		Password: "",
		DB:       0,
	})
	return rdb
}

// connectMongoDB function to connect to MongoDB
func connectMongoDB() *mongo.Client {
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
	collection := MONGODB.Database("periphery").Collection(collectionName)
	return collection
}

func connectMysql() *sql.DB {
	db, err := sql.Open("mysql", EnvMysqlURI())
	if err != nil {
		log.Fatal("Failed to connect to Mysql", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping: %v", err)
	}
	log.Println("Successfully connected to Mysql!")

	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}
