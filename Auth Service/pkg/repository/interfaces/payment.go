package interfaces

import (
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/domain"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/models/responsemodels"
)

type IPaymentRepository interface {
	GetBlueTickVerificationPrice() uint
	CreateBlueTickVerification(userID uint, verificationID string) error
	UpdateBlueTickPaymentSuccess(verificationID string) (*domain.BlueTickVerification, error) //GetBlueTickVerificationStatus(userID uint) (*responsemodels.OnlinePayment, error)
	OnlinePayment(userID, verificationID string) (*responsemodels.OnlinePayment, error)
	UpdateBluetickStatus(userID uint) error
	GetAllVerifiedUsers(limit, offset int) (*[]responsemodels.BlueTickResponse, error)
	IsUserVerified(userID string) (bool, error) 
}
