package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Client membungkus connection pool ke Postgres (Supabase).
type Client struct {
	Pool *pgxpool.Pool
}

// Connect membuka connection pool dan memverifikasi dengan ping.
func Connect(dsn string) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &Client{Pool: pool}, nil
}
