package server

import (
	"context"
	"fmt"
	"github.com/ilhamtubagus/go-shorten-url/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
)

func ConnectMongoClient(mongoConfig config.MongoConfig) *mongo.Client {
	connStr := fmt.Sprintf("mongodb://%s:%s@%s/%s?%s",
		mongoConfig.User,
		mongoConfig.Password,
		mongoConfig.Host,
		mongoConfig.Database,
		mongoConfig.Options)

	clientOptions := options.Client().ApplyURI(connStr)

	log.Println("connecting to MongoDB server", connStr)

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
