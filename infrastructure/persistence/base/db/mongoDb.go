package db

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
func NewMongoDB() (*mongo.Client, error) {
	DbHost := os.Getenv("DB_MONGO_HOST")
	// DBName := os.Getenv("DB_MONGO_NAME")
	// DbCollection := os.Getenv("DB_MONGO_COLLECTION_PRODUCTS")

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(DbHost))
	if err != nil {
		return nil, err
	}

	// Ping the deployment to confirm a successful connection
	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")


	return client, nil
}

