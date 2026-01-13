package storage

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// Initialize MongoDB with authentication
func InitMongo() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// MongoDB connection config with credentials
	clientOptions := options.Client().
		ApplyURI("mongodb://localhost:27017/?authSource=admin").
		SetAuth(options.Credential{
			Username:      "nandini-y",
			Password:      "Meizo@_2026",
			AuthSource:    "admin",
			AuthMechanism: "SCRAM-SHA-256",
		})

	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println("❌ Mongo Connect Error:", err)
		return err
	}

	// Ping DB
	if err := client.Ping(ctx, nil); err != nil {
		log.Println("❌ Mongo Ping Error:", err)
		return err
	}

	log.Println("✅ MongoDB Connected Successfully")
	return nil
}

// Return Mongo client
func GetMongo() *mongo.Client {
	return client
}

// Helper: get a database
func GetDatabase(dbName string) *mongo.Database {
	return client.Database(dbName)
}

// Helper: get a collection inside "fitpro"
func GetCollection(collectionName string) *mongo.Collection {
	return client.Database("fitpro").Collection(collectionName)
}

// Disconnect cleanly
func CloseMongo() {
	if client != nil {
		_ = client.Disconnect(context.Background())
	}
}
