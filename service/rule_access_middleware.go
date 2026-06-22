package main

import (
	"context"
	"database/sql"
	"hubit-space/service/config"
	"hubit-space/service/middleware"
	"hubit-space/service/router"
	"log"

	"firebase.google.com/go/v4/messaging"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var db *gorm.DB
var hdb *sql.DB
var red *redis.Client
var mdb *mongo.Client
var fbm *messaging.Client

func init() {
	config.LoadConfig()

	var err error
	db = config.OpenPostgresConnection()
	log.Println("Connected to PostgreSQL database")

	red, err = config.OpenRedisConnection()
	if err != nil {
		log.Fatal("failed to connect to redis: " + err.Error())
	}
	log.Println("connected to redis")

	mdb, err = config.OpenMongoConnection()
	if err != nil {
		log.Fatal("Failed to connect to MongoDB service: " + err.Error())
	} else {
		log.Println("Connected to MongoDB service")
	}

}

func main() {
	r := gin.Default()

	// Panic recovery middleware
	r.Use(middleware.PanicMiddleware())

	// Middleware
	r.Use(middleware.AuthMiddleware())

	// All routes
	router.SetupOptionRouter(db, r)

	// Run the server
	log.Println("Starting server on port 8081")
	r.Run(":8081")

	// Close all connections
	defer func() {
		if db != nil {
			sqlDB, err := db.DB()
			if err != nil {
				log.Printf("Error getting SQL DB: %v", err)
			} else {
				sqlDB.Close()
				log.Println("PostgreSQL connection closed")
			}
		}

		if hdb != nil {
			hdb.Close()
			log.Println("Hana database connection closed")
		}

		if red != nil {
			err := red.Close()
			if err != nil {
				log.Printf("Error closing Redis connection: %v", err)
			} else {
				log.Println("Redis connection closed")
			}
		}

		if mdb != nil {
			err := mdb.Disconnect(context.Background())
			if err != nil {
				log.Printf("Error closing MongoDB connection: %v", err)
			} else {
				log.Println("MongoDB connection closed")
			}
		}

		log.Println("Server shutdown complete")
	}()
}
