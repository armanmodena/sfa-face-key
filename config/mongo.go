package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func OpenMongoConnection() (*mongo.Client, error) {
	url := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin",
		MONGO_USER,
		MONGO_PASSWORD,
		MONGO_HOST,
		MONGO_PORT)

	clientOptions := options.Client().ApplyURI(url).
		SetMaxPoolSize(50).
		SetMinPoolSize(10).
		SetConnectTimeout(10 * time.Second).
		SetServerSelectionTimeout(10 * time.Second)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		return nil, err
	}

	log.Println("Successfully Connected to MongoDB")
	return client, nil
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database(MONGO_DB).Collection(collectionName)
}
