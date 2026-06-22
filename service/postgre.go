package config

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDBName = GetEnv("MONGO_DB_NAME", "hubit_space")

func OpenMongoConnection() (*mongo.Client, error) {
	mongoUser := url.QueryEscape(GetEnv("MONGO_USER", ""))
	mongoPassword := url.QueryEscape(GetEnv("MONGO_PASSWORD", ""))
	mongoHost := GetEnv("MONGO_HOST", "")

	if mongoUser == "" || mongoPassword == "" || mongoHost == "" {
		return nil, fmt.Errorf("MONGO_USER, MONGO_PASSWORD, and MONGO_HOST must be set")
	}

	uri := fmt.Sprintf("mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority",
		mongoUser,
		mongoPassword,
		mongoHost,
	)

	clientOptions := options.Client().ApplyURI(uri).
		SetConnectTimeout(30 * time.Second).
		SetServerSelectionTimeout(5 * time.Second).
		SetMaxConnIdleTime(30 * time.Minute).
		SetMaxPoolSize(20).
		SetMinPoolSize(0).
		SetRetryWrites(true).
		SetRetryReads(true)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		log.Println("Failed to ping MongoDB:", err)
		return nil, err
	}

	log.Println("Successfully connected to MongoDB")
	return client, nil
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database(MongoDBName).Collection(collectionName)
}
