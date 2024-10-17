package database

import (
    "database/sql"
    "log"
    _ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
    // Hardcoded connection string example
    connStr := "postgres://postgres:ravjotravjot@localhost:5432/fmc?sslmode=disable"

    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Error opening database: %v", err)
    }

    err = DB.Ping()
    if err != nil {
        log.Fatalf("Error connecting to the database: %v", err)
    }

    log.Println("Successfully connected to the database.")
}
