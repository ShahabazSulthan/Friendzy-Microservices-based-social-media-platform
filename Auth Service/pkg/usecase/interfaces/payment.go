package interface_usecase

import "github.com/ShahabazSulthan/Friendzy_Auth/pkg/models/responsemodels"

type IPaymentUsecase interface {
	CreateBlueTickPayment(userID uint) (string, error)
	VerifyBlueTickPayment(orderID, paymentID, signature string, userID uint) (bool, error)
	OnlinePayment(userID, verificationID string) (*responsemodels.OnlinePayment, error)
	GetAllVerifiedUsers(limit string, offset string) (*[]responsemodels.BlueTickResponse, error)
	IsUserVerified(userID string) (bool, error)
}
