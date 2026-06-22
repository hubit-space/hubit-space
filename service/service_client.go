package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func OpenRedisConnection() (*redis.Client, error) {
	redisAddr := GetEnv("REDIS_HOST", "localhost:6379")
	redisPassword := GetEnv("REDIS_PASSWORD", "")
	redisDB := GetEnv("REDIS_DB", "0")
	redisTLS := GetEnv("REDIS_TLS", "true")

	redisDBInt, err := strconv.Atoi(redisDB)
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_DB value %q: %w", redisDB, err)
	}

	opts := &redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDBInt,

		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     20,
		MinIdleConns: 10,

		MaxRetries:      3,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,

		PoolTimeout:     4 * time.Second,
		ConnMaxIdleTime: 10 * time.Minute,
		ConnMaxLifetime: 0,
	}

	if redisTLS == "true" {
		opts.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
