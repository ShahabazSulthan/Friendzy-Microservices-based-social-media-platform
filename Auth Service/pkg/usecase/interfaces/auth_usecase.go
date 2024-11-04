package interface_usecase

import (
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/models/requestmodels"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/models/responsemodels"
)

type IUserUseCase interface {
	UserSignUp(userData *requestmodels.UserSignUpRequest) (responsemodels.SignUpResponse, error)
	VerifyOtp(otp string, TempVerificationToken *string) (responsemodels.OTPVerificationResponse, error)
	UserLogin(loginData *requestmodels.UserLoginRequest) (responsemodels.LoginResponse, error)
	ForgetPasswordRequest(email *string) (*string, error)

	ResetPassword(userData *requestmodels.ForgotPasswordRequest, TempVerification *string) error
	UserProfile(userId, UserBId *string) (*responsemodels.UserProfileResponse, error)
	EditUserDetails(editInput *requestmodels.EditUserProfileRequest) error
	GetUserDetailsLiteForPostView(userId *string) (*responsemodels.UserLiteResponse, error)
	GetEmailAndUsername(userID int) (string, string, error)

	GetFollowersDetails(userId *string) (*[]responsemodels.UserListResponse, *error)
	GetFollowingDetails(userId *string) (*[]responsemodels.UserListResponse, *error)
	SearchUser(myId, SearchText, limit, offset *string) (*[]responsemodels.UserListResponse, error)
	SetUserProfileImage(userId, ContenType *string, Img *[]byte) error
	
	CheckUserExist(userId *string) (bool, *error)
}
