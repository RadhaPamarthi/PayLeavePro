// db/mongo.go
package db

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client *mongo.Client
	once   sync.Once
	mu     sync.Mutex
)

// Connect establishes connection to MongoDB
func Connect() error {
	var err error
	once.Do(func() {
		// Create a context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Set client options
		clientOptions := options.Client().ApplyURI("mongodb+srv://iamradha0246:IcrOVusDHSrGh0X2@cluster0.rhzjo.mongodb.net/?retryWrites=true&w=majority")

		// Connect to MongoDB
		Client, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Printf("Failed to connect to MongoDB: %v", err)
			return
		}

		// Ping the database
		err = Client.Ping(ctx, nil)
		if err != nil {
			log.Printf("Failed to ping MongoDB: %v", err)
			return
		}

		log.Println("Successfully connected to MongoDB!")
	})
	return err
}

// GetCollection returns a MongoDB collection with connection check
func GetCollection(collectionName string) (*mongo.Collection, error) {
	mu.Lock()
	defer mu.Unlock()

	if Client == nil {
		if err := Connect(); err != nil {
			return nil, err
		}
	}

	// Ping to ensure connection is still alive
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Client.Ping(ctx, nil); err != nil {
		log.Printf("Lost connection to MongoDB, attempting to reconnect: %v", err)
		if err := Connect(); err != nil {
			return nil, err
		}
	}

	return Client.Database("PayLeavePro").Collection(collectionName), nil
}
