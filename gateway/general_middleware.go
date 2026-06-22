package main

import (
	"context"
	"hubit-space/gateway/config"
	"hubit-space/gateway/middleware"
	"hubit-space/gateway/router"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

var ctx context.Context
var client *resty.Client
var mongoClient *mongo.Client
var redisClient *redis.Client

func init() {
	config.LoadEnv()

	var err error
	mongoClient, err = config.OpenMongoConnection()
	if err != nil {
		log.Fatal(err)
	}

	redisClient, err = config.OpenRedisConnection()
	if err != nil {
		log.Fatal(err)
	}

	client, err = config.OpenServiceConnection()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	gin := gin.Default()

	// CORS middleware
	gin.Use(middleware.CORSMiddleware())

	// Panic middleware
	gin.Use(middleware.PanicMiddleware())

	// Public routes
	router.SetupAuthRouter(gin, mongoClient, redisClient, client)

	// Auth middleware
	gin.Use(middleware.AuthMiddleware())

	// Non public routes
	router.SetupUserRouter(gin, mongoClient, client)

	// Run the server
	log.Println("Starting server on port 8080")
	gin.Run(":8080")

	// Close all connections
	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Fatalf(`Error disconnectiong from MongoDB: %v`, err)
		} else {
			log.Println(`Successfully disconnected from MongoDB`)
		}

		if err := redisClient.Close(); err != nil {
			log.Fatalf(`Error disconnectiong from Redis: %v`, err)
		} else {
			log.Println(`Successfully disconnected from Redis`)
		}

		client.GetClient().CloseIdleConnections()
		log.Println("Server shutdown complete")
	}()
}
