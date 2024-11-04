package interfaces

import (
	"time"

	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/models/requestmodels"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/models/responsemodels"
)

type IUserRepository interface {
    UserExistsByEmail(email string) bool
    UserExistsByUsername(username string) bool
    DeleteRecentOtpRequests() error
    SaveUserOtp(otp int, userEmail string, expiry time.Time) error

    CreateUser(userData *requestmodels.UserSignUpRequest) error
    GetOtpInfo(email string) (string, time.Time, error)
    ActivateUser(email string) error
    GetUserIDByEmail(email string) (string, error)
    GetEmailAndUsernameByUserID(userID int) (string, string, error)

    GetUserPasswordHashAndStatus(email string) (string, string, string, error)
    UpdateUserPassword(email, hashedPassword *string) error
    GetUserStateForAccessTokenGeneration(userID *string) (*string, error)
    GetUserLiteProfile(userID *string) (*responsemodels.UserProfileResponse, error)

    UpdateUserLiteProfile(editInput *requestmodels.EditUserProfileRequest) error
    GetUserProfileUrlAndUsername(userID *string) (*responsemodels.UserLiteResponse, error)
    UserExistsByID(userID string) (bool,*error)
    SearchUsersByName(myID, searchText, limit, offset *string) (*[]responsemodels.UserListResponse, error)

    SetUserProfileImage(userID, imageUrl *string) error
    GetFollowersDetails(userIDs *[]uint64) (*[]responsemodels.UserListResponse, error)
    GetFollowingDetails(userIDs *[]uint64) (*[]responsemodels.UserListResponse, error)
    IsUserBlocked(userID string) (bool, error) 
}
