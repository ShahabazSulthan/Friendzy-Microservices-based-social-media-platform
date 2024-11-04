package repository

import (
	"errors"
	"fmt"
	"log"

	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/models/responsemodels"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/repository/interfaces"
	"gorm.io/gorm"
)

type AdminRepo struct {
	DB *gorm.DB
}

func NewAdminRepo(db *gorm.DB) interfaces.IAdminRepo {
	return &AdminRepo{DB: db}
}

func (a AdminRepo) GetPassword(email string) (string, error) {
	var hashedPassword string

	query := "SELECT password FROM admins WHERE email =?"
	err := a.DB.Raw(query, email).Row().Scan(&hashedPassword)
	if err != nil {
		log.Println("Error fetching admin password: ", err)
		return "", errors.New("error fetching admin password")
	}

	return hashedPassword, nil
}

func (a AdminRepo) AllUsers(limit, offset int) (*[]responsemodels.UserAdminResponse, error) {
	var users []responsemodels.UserAdminResponse

	// SQL query to get users
	query := `
        SELECT id, name, user_name, email, bio, profile_img_url, links, status 
        FROM users
        WHERE deleted_at IS NULL
        LIMIT $1 OFFSET $2;
    `

	// Execute the query and scan the results
	err := a.DB.Raw(query, limit, offset).Scan(&users).Error
	if err != nil {
		fmt.Println("Error in fetching users:", err)
		return nil, err
	}

	return &users, nil
}

func (a *AdminRepo) BlockUser(id string) error {
	fmt.Println("id = ", id)

	// Corrected query for updating user status
	query := "UPDATE users SET status = 'blocked' WHERE id = ?"

	// Use Exec instead of Raw for an update query
	result := a.DB.Exec(query, id)

	// Check if there's any error during the execution
	if result.Error != nil {
		fmt.Println("Error blocking user:", result.Error)
		return errors.New("block user process is not satisfied")
	}

	// Optionally, you can check if any rows were affected
	if result.RowsAffected == 0 {
		return errors.New("no user was found with the provided id")
	}

	fmt.Println("User blocked successfully")
	return nil
}

func (a *AdminRepo) UnblockUser(id string) error {
	query := "UPDATE users SET status = 'active' WHERE id = ?"

	// Use Exec for update queries
	result := a.DB.Exec(query, id)

	// Check if there's any error during execution
	if result.Error != nil {
		return errors.New("active user process is not satisfied")
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return errors.New("no user exists by the provided id")
	}

	fmt.Println("User unblocked successfully")
	return nil
}
