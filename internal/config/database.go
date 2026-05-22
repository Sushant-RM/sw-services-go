package config

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func NewDatabaseConnection(cfg *Config) (*sql.DB, error) {

	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	var db *sql.DB
	var err error

	maxRetries := 10

	for attempt := 1; attempt <= maxRetries; attempt++ {

		db, err = sql.Open(
			"postgres",
			connectionString,
		)

		if err == nil {

			err = db.Ping()
			if err == nil {

				log.Println("Database connection established")

				return db, nil
			}
		}

		log.Printf(
			"Database connection failed (attempt %d/%d): %v",
			attempt,
			maxRetries,
			err,
		)

		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf(
		"failed to connect to database after %d attempts: %w",
		maxRetries,
		err,
	)
}
