package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

var Pool *pgxpool.Pool

// ConnectDB establishes a connection pool to the PostgreSQL database
func ConnectDB(uri string) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(uri)
	if err != nil {
		log.Fatalf("Unable to parse DATABASE_URL: %v\n", err)
	}

	// Set pool configurations
	config.MaxConns = 100
	config.MinConns = 10
	config.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Ping to ensure the connection is established
	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	log.Println("Connected to PostgreSQL!")
	Pool = pool
	return Pool
}

// ClosePool gracefully closes the database connection pool
func ClosePool() {
	if Pool != nil {
		Pool.Close()
	}
}

// GetPool returns the current database pool
func GetPool() *pgxpool.Pool {
	if Pool == nil {
		log.Fatal("PostgreSQL pool is not initialized")
	}
	return Pool
}
