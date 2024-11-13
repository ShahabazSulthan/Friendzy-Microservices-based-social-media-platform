package repository

import (
	"errors"
	"fmt"

	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/domain"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/models/responsemodels"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/repository/interfaces"

	"gorm.io/gorm"
)

type PaymentRepo struct {
	DB *gorm.DB
}

// NewPaymentRepo creates a new instance of PaymentRepo.
func NewPaymentRepo(db *gorm.DB) interfaces.IPaymentRepository {
	return &PaymentRepo{
		DB: db,
	}
}

// GetBlueTickVerificationPrice fetches the current price for Blue Tick verification.
func (p *PaymentRepo) GetBlueTickVerificationPrice() uint {
	// Set a fixed or fetch from config/database as per your requirement
	const blueTickPrice uint = 600 // example price
	return blueTickPrice
}

func (p *PaymentRepo) CreateBlueTickVerification(userID uint, verificationID string) error {
	// Define the query for inserting the verification entry into the database
	query := `INSERT INTO blue_tick_verifications (user_id, verification_id, created_at, status, verification_fee) VALUES (?, ?, NOW(), 'pending', ?)`

	verificationFee := 1000
	// Execute the query using the database connection
	db := p.DB.Exec(query, userID, verificationID, verificationFee)
	if db.Error != nil {
		return fmt.Errorf("error saving blue tick verification: %w", db.Error) // Use db.Error for better debugging
	}

	return nil
}

// UpdateBlueTickPaymentSuccess updates the blue tick verification status to 'Success'.
func (p *PaymentRepo) UpdateBlueTickPaymentSuccess(verificationID string) (*domain.BlueTickVerification, error) {
	var blueTickVerification domain.BlueTickVerification

	// Find the verification record by verification_id
	if err := p.DB.Where("verification_id = ?", verificationID).First(&blueTickVerification).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("blue tick verification record not found")
		}
		return nil, err
	}

	// Update the status to "Success"
	blueTickVerification.Status = "Success"
	if err := p.DB.Save(&blueTickVerification).Error; err != nil {
		return nil, err
	}
	return &blueTickVerification, nil
}

// // GetBlueTickVerificationStatus retrieves the blue tick verification payment status for a user.
// func (p *PaymentRepo) GetBlueTickVerificationStatus(userID uint) (*responsemodels.OnlinePayment, error) {
// 	// var blueTickVerification domain.BlueTickVerification

// 	// if err := p.DB.Where("user_id = ?", userID).First(&blueTickVerification).Error; err != nil {
// 	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
// 	// 		return nil, errors.New("blue tick verification record not found")
// 	// 	}
// 	// 	return nil, err
// 	// }

// 	// // Prepare the response model
// 	// status := &responsemodels.OnlinePayment{
// 	// 	UserID: userID,
// 	// 	Status: blueTickVerification.Status,
// 	// 	Amount: p.GetBlueTickVerificationPrice(),
// 	// }
// 	// return status, nil
// }

func (p *PaymentRepo) OnlinePayment(userID, verificationID string) (*responsemodels.OnlinePayment, error) {
	var orderDetails responsemodels.OnlinePayment

	// Query to fetch necessary fields for OnlinePayment response
	query := `
        SELECT user_id, status AS payment_status, verification_fee
        FROM blue_tick_verifications
        WHERE user_id = ? AND verification_id = ?;
    `
	// Execute query and scan the result into orderDetails struct
	result := p.DB.Raw(query, userID, verificationID).Scan(&orderDetails)

	// Handle errors and no rows affected cases
	if result.Error != nil {
		return nil, fmt.Errorf("error executing database query: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("no records found for user ID %s and verification ID %s", userID, verificationID)
	}

	return &orderDetails, nil
}

func (p *PaymentRepo) UpdateBluetickStatus(userID uint) error {
	// SQL query to update the is_bluetick_verified column
	query := `
        UPDATE users
        SET is_bluetick_verified = TRUE
        WHERE id = ?;
    `

	// Execute the update query
	result := p.DB.Exec(query, userID)

	// Check for errors in the execution
	if result.Error != nil {
		return fmt.Errorf("error executing update query: %w", result.Error)
	}

	// Check if no rows were affected (user not found)
	if result.RowsAffected == 0 {
		return fmt.Errorf("no user found with user Id %v", userID)
	}

	return nil
}

func (p *PaymentRepo) GetAllVerifiedUsers(limit, offset int) (*[]responsemodels.BlueTickResponse, error) {
	var users []responsemodels.BlueTickResponse

	// SQL query to get verified users
	query := `
        SELECT id, name, user_name, email, bio, profile_img_url, links, status 
        FROM users
        WHERE deleted_at IS NULL 
        AND is_bluetick_verified = true
        LIMIT $1 OFFSET $2;
    `

	// Execute the query and scan the results
	err := p.DB.Raw(query, limit, offset).Scan(&users).Error
	if err != nil {
		fmt.Println("Error in fetching users:", err)
		return nil, err
	}

	return &users, nil
}

func (p *PaymentRepo) IsUserVerified(userID string) (bool, error) {
    var isVerified bool

    // SQL query to check if the user is verified
    query := `
        SELECT is_bluetick_verified 
        FROM users 
        WHERE id = $1 AND deleted_at IS NULL;
    `

    // Execute the query and store the result in isVerified
    err := p.DB.Raw(query, userID).Scan(&isVerified).Error
    if err != nil {
        fmt.Println("Error in checking verification status:", err)
        return false, err
    }

    return isVerified, nil
}