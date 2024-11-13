package usecase

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/config"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/models/responsemodels"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/repository/interfaces"
	interface_usecase "github.com/ShahabazSulthan/Friendzy_Auth/pkg/usecase/interfaces"
	razorpay "github.com/ShahabazSulthan/Friendzy_Auth/pkg/utils/razopay"
)

type RazopayUseCase struct {
	Repo    interfaces.IPaymentRepository
	Razopay *config.Razopay
}

func NewPaymenyUsecase(repo interfaces.IPaymentRepository, razopay *config.Razopay) interface_usecase.IPaymentUsecase {
	return &RazopayUseCase{
		Repo:    repo,
		Razopay: razopay,
	}
}

// CreateBlueTickVerification saves the verification ID in the database.
func (r *RazopayUseCase) CreateBlueTickPayment(userID uint) (string, error) {
	// Create Razorpay order with specified amount and currency (in this case INR)
	order, err := razorpay.Razopay(r.Razopay.RazopayKey, r.Razopay.RazopaySecret)
	if err != nil {
		return "", err
	}

	// Save the order ID to the database by calling the repository method
	err = r.Repo.CreateBlueTickVerification(userID, order)
	if err != nil {
		return "", err
	}

	// Return the order ID as verification ID
	return order, nil
}

// VerifyBlueTickPayment verifies the payment using Razorpay signature verification and updates the status.
func (r *RazopayUseCase) VerifyBlueTickPayment(verificationID, paymentID, signature string, userID uint) (bool, error) {

	isVerified := razorpay.VerifyPayment(verificationID, paymentID, signature, r.Razopay.RazopaySecret)
	if !isVerified {
		return false, errors.New("payment verification failed")
	}

	fmt.Println("userid = ", userID)

	if isVerified {
		err := r.Repo.UpdateBluetickStatus(userID)
		if err != nil {
			return false, errors.New("cant Update user table")
		}
	}

	_, err := r.Repo.UpdateBlueTickPaymentSuccess(verificationID)
	if err != nil {
		return false, err
	}

	// Step 3: Return success
	// If both verification and status update are successful, the function returns `true`.
	return true, nil
}

// // GetBlueTickVerificationStatus retrieves the payment status for a user.
// func (r *RazopayUseCase) GetBlueTickVerificationStatus(userID uint) (string, error) {
// 	status, err := r.Repo.GetBlueTickVerificationStatus(userID)
// 	if err != nil {
// 		return "", err
// 	}
// 	return status.Status, nil
// }

// GetBlueTickVerification retrieves the blue tick verification details for a given userID and verificationID.
func (u *RazopayUseCase) OnlinePayment(userID, verificationID string) (*responsemodels.OnlinePayment, error) {
	// Call the OnlinePayment method from the repository
	paymentDetails, err := u.Repo.OnlinePayment(userID, verificationID)
	if err != nil {
		return nil, err
	}

	paymentDetails.VerificationFee = 1000

	return paymentDetails, nil
}

func (u *RazopayUseCase) GetAllVerifiedUsers(limit string, offset string) (*[]responsemodels.BlueTickResponse, error) {
	// Convert limit and offset from string to int
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return nil, errors.New("invalid limit value")
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		return nil, errors.New("invalid offset value")
	}

	// Call the repository method to fetch users
	users, err := u.Repo.GetAllVerifiedUsers(limitInt, offsetInt)
	if err != nil {
		return nil, err // Return any errors from the repository
	}

	return users, nil // Return the list of users
}

func (u *RazopayUseCase) IsUserVerified(userID string) (bool, error) {
	verified, err := u.Repo.IsUserVerified(userID)
	if err != nil {
		return false, err
	}

	return verified, nil
}
