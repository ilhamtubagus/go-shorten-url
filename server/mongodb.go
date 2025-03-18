package server

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"os"
)

func ConnectMongoClient() *mongo.Client {
	connStr := fmt.Sprintf("mongodb://%s:%s@%s/%s?%s",
		os.Getenv("MONGODB_USER"),
		os.Getenv("MONGODB_PASSWORD"),
		os.Getenv("MONGODB_HOST"),
		os.Getenv("MONGODB_DATABASE_NAME"),
		os.Getenv("MONGODB_OPTIONS"))

	clientOptions := options.Client().ApplyURI(connStr)

	log.Println("connecting to MongoDB server")

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %e\n", err)
	}

	// Ping the database to verify the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("failed to ping MongoDB: %e\n", err)
	}

	log.Printf("connected to MongoDB: %s\n", connStr)

	return client
}
