package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/config"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/domain"
	interface_hash "github.com/ShahabazSulthan/Friendzy_Auth/pkg/utils/hashed_password/interfaces"
	_ "github.com/lib/pq" // Import the postgres driver
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDatabase connects to the database, checks if the specified database exists, and performs migrations.
func ConnectDatabase(config *config.DataBase, hash interface_hash.IHashPassword) (*gorm.DB, error) {
	// Connection string to check if the database exists
	connectionString := fmt.Sprintf("host=%s user=%s password=%s port=%s sslmode=disable", config.DBHost, config.DBUser, config.DBPassword, config.DBPort)
	sqlDB, err := sql.Open("postgres", connectionString) // Ensure you import the driver
	if err != nil {
		fmt.Println("Error connecting to Postgres:", err)
		return nil, err
	}
	defer sqlDB.Close()

	// Check if the database exists
	rows, err := sqlDB.Query("SELECT 1 FROM pg_database WHERE datname = $1", config.DBName)
	if err != nil {
		fmt.Println("Error querying for database existence:", err)
		return nil, err
	}
	defer rows.Close()

	// If database does not exist, create it
	if !rows.Next() {
		_, err = sqlDB.Exec("CREATE DATABASE " + config.DBName)
		if err != nil {
			fmt.Println("Error in creating database:", err)
			return nil, err
		}
		fmt.Println("Database created:", config.DBName)
	} else {
		fmt.Println("Database", config.DBName, "already exists")
	}

	// Connection string for GORM to connect to the newly created database
	psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=disable", config.DBHost, config.DBUser, config.DBName, config.DBPassword, config.DBPort)
	DB, dberr := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if dberr != nil {
		fmt.Println("Error connecting to database with GORM:", dberr)
		return nil, dberr
	}

	// Perform migrations
	if err := DB.AutoMigrate(&domain.Admin{}, &domain.User{}, &domain.OTP{},&domain.BlueTickVerification{}); err != nil {
		fmt.Println("Error in migrating database:", err)
		return nil, err
	}

	// Create default admin
	CreateAdmin(DB, hash)

	return DB, nil
}

// CreateAdmin creates a default admin if no admin exists in the database
func CreateAdmin(DB *gorm.DB, hash interface_hash.IHashPassword) {
	var count int64

	Name := "Friendzy"
	Email := "friendzy@gmail.com"
	Password := "friendzy123"

	hashedPassword := hash.HashedPassword(Password)

	// Count number of admins
	DB.Table("admins").Count(&count)

	if count == 0 {
		admin := domain.Admin{
			Name:     Name,
			Email:    Email,
			Password: hashedPassword,
		}
		if err := DB.Create(&admin).Error; err != nil {
			fmt.Println("Error in creating admin:", err)
		} else {
			fmt.Println("Default admin created.")
		}
	} else {
		fmt.Println("Admin already exists.")
	}
}
