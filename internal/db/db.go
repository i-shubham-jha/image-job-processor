package db

import (
	"context"
	"log"
	"os"
	"retail_pulse/internal/logger"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	clientInstance *mongo.Client
	once           sync.Once
)

// GetMongoClient returns a singleton instance of the MongoDB client
func GetMongoClient() *mongo.Client {
	once.Do(func() {
		var err error
		mongoURI := os.Getenv("MONGODB_URI") // Read the MongoDB URI from the environment variable
		if mongoURI == "" {
			log.Fatal("MONGODB_URI environment variable is not set")
		}

		clientOptions := options.Client().ApplyURI(mongoURI)
		clientInstance, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			log.Fatal(err)
		}

		// Check the connection
		err = clientInstance.Ping(context.TODO(), nil)
		if err != nil {
			log.Fatal(err)
		}

		logger.GetLogger().Log("Successfully established connection to mongodb")
	})
	return clientInstance
}
