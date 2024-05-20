package database

import (
	"database/sql"
	"fmt"
	"log"
	"miracummapper/internal/config"
	"time"

	_ "github.com/lib/pq"
)

func NewDBConnection(config *config.Config) *sql.DB {
	db, err := getDBConnection(config)
	if err != nil {
		panic(err)
	}
	return db
}

func getDBConnection(config *config.Config) (*sql.DB, error) {

	db, err := createDatabaseConnection(config)
	if err != nil {
		return nil, err
	}

	for i := 0; i < config.File.DatabaseConfig.Retry; i++ {
		err = pingDb(db)
		if err == nil {
			log.Printf("Successfully connected to database: %s", config.Env.DBName)
			break
		}
		log.Printf("Failed to connect to database: %s. Retrying in %d seconds", config.Env.DBName, config.File.DatabaseConfig.Sleep)
		if i != config.File.DatabaseConfig.Retry-1 {
			time.Sleep(time.Duration(config.File.DatabaseConfig.Sleep) * time.Second)
		}
	}

	return db, nil
}

func createDatabaseConnection(config *config.Config) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Env.DBHost,
		config.Env.DBPort,
		config.Env.DBUser,
		config.Env.DBPassword,
		config.Env.DBName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}

func pingDb(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	return nil
}
