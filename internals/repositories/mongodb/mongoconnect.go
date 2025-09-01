package mongodb

import (
	"context"
	"crypto/tls"
	"fmt"
	"grpc-school-mgnt/pkg/utils"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var client *mongo.Client

// Connect initializes a global MongoDB client (reused with connection pool)
func Connect(uri string) error {
	if client != nil {
		return nil
	}

	opts := options.Client().
		ApplyURI(uri).
		SetTLSConfig(&tls.Config{InsecureSkipVerify: true})

	c, err := mongo.Connect(opts)
	if err != nil {
		return utils.ErrorHandler(err, "❌ Failed to connect to MongoDB")
	}

	// Ping the database
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := c.Ping(ctx, readpref.Primary()); err != nil {
		return utils.ErrorHandler(err, "⚠️ Ping failed")
	}

	fmt.Println("✅ Connected to MongoDB Atlas!")

	client = c

	return nil
}

func Client() *mongo.Client {
	if client == nil {
		log.Fatal("MongoDB client is not initialized. Call Connect() first.")
	}
	return client
}

func Disconnect() {
	if client == nil {
		return
	}
	if err := client.Disconnect(context.Background()); err != nil {
		log.Fatalf("Error disconnecting MongoDB: %v", err)
	}
	fmt.Println("👋 MongoDB disconnected.")
}
