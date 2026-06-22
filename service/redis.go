package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func OpenPostgresConnection() *gorm.DB {
	dsn := fmt.Sprintf(`host=%s user=%s password=%s dbname=%s port=%s sslmode=require connect_timeout=30 TimeZone=Asia/Jakarta`,
		GetEnv("DB_HOST", ""),
		GetEnv("DB_USER", ""),
		GetEnv("DB_PASSWORD", ""),
		GetEnv("DB_NAME", "postgres"),
		GetEnv("DB_PORT", "5432"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Warn),
		PrepareStmt: true,
	})

	if err != nil {
		log.Fatal("Failed to connect PostgreSQL database: " + err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get PostgreSQL DB instance: " + err.Error())
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)

	return db
}
