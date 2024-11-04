package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ShahabazSulthan/Friendzy_Notification/pkg/config"
	"github.com/ShahabazSulthan/Friendzy_Notification/pkg/domain"
	interface_hash "github.com/ShahabazSulthan/Friendzy_Notification/pkg/utils/hash/interface"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDatabase establishes a connection, checks for database existence, and performs migrations.
func ConnectDatabase(config *config.DataBase, hash interface_hash.Ihash) (*gorm.DB, error) {
	// Initial connection to PostgreSQL to verify or create database
	connectionString := fmt.Sprintf("host=%s user=%s password=%s port=%s sslmode=disable", config.DBHost, config.DBUser, config.DBPassword, config.DBPort)
	sqlDB, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Postgres: %w", err)
	}
	defer sqlDB.Close()

	// Check for database existence
	rows, err := sqlDB.Query("SELECT 1 FROM pg_database WHERE datname = $1", config.DBName)
	if err != nil {
		return nil, fmt.Errorf("error checking database existence: %w", err)
	}
	defer rows.Close()

	// Create database if it doesn't exist
	if !rows.Next() {
		_, err = sqlDB.Exec("CREATE DATABASE " + config.DBName)
		if err != nil {
			return nil, fmt.Errorf("error creating database: %w", err)
		}
		fmt.Println("Database created:", config.DBName)
	} else {
		fmt.Println("Database", config.DBName, "already exists")
	}

	psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", config.DBHost, config.DBUser, config.DBName, config.DBPort, config.DBPassword)
	DB, dberr := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC() // Set the timezone to UTC
		},
	})
	if dberr != nil {
		return DB, nil
	}

	// Table Creation
	if err := DB.AutoMigrate(&domain.Notification{}); err != nil {
		return DB, err
	}

	return DB, nil
}
