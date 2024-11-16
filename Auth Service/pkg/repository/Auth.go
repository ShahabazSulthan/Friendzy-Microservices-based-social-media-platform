package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/models/requestmodels"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/models/responsemodels"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/repository/interfaces"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) interfaces.IUserRepository {
	return &UserRepo{
		DB: db,
	}
}

func (u *UserRepo) UserExistsByEmail(email string) bool {
	var UserCount int

	delUncomplete := "DELETE FROM users WHERE email = $1 AND status= $2"
	result := u.DB.Exec(delUncomplete, email, "pending")
	if result.Error != nil {
		fmt.Println("Error in deleting already existing user with this email and status pending")
	}

	query := "SELECT COUNT(*) FROM users WHERE email=$1 AND status !=$2"
	err := u.DB.Raw(query, email, "deleted").Row().Scan(&UserCount)
	if err != nil {
		fmt.Println("Error in users Count query")
	}

	if UserCount >= 1 {
		return true
	}

	return false
}

func (u *UserRepo) GetEmailAndUsernameByUserID(userID int) (string, string, error) {
	var email, username string

	// Query to fetch email and username by user_id
	query := "SELECT email, user_name FROM users WHERE user_id = $1 "
	err := u.DB.Raw(query, userID).Row().Scan(&email, &username)
	if err != nil {
		fmt.Println("Error retrieving email and username:", err)
		return "", "", err
	}

	// Return email and username if the user exists
	return email, username, nil
}

func (u *UserRepo) UserExistsByUsername(username string) bool {
	var UserCount int

	// Correct the query method to use Raw for a SELECT statement
	query := "SELECT COUNT(*) FROM users WHERE user_name = $1 AND status != $2"
	err := u.DB.Raw(query, username, "deleted").Scan(&UserCount).Error
	if err != nil {
		fmt.Println("Error in UserName Count query:", err)
		return false
	}

	// Return true if a user exists, otherwise return false
	return UserCount >= 1
}

func (u *UserRepo) DeleteRecentOtpRequests() error {
	query := "DELETE FROM otps WHERE expiration < CURRENT_TIMESTAMP - INTERVAL '5 minutes'"
	err := u.DB.Exec(query).Error
	if err != nil {
		fmt.Println("Error In Delete Recent OtP Request")
		return err
	}
	return nil
}

func (u *UserRepo) SaveUserOtp(otp int, userEmail string, expiry time.Time) error {
	query := `INSERT INTO otps (email,otp,expiration) VALUES ($1,$2,$3)`
	err := u.DB.Exec(query, userEmail, otp, expiry).Error
	if err != nil {
		fmt.Println("Error In Save User Otp")
	}

	fmt.Println("OTP = ", otp)
	return nil
}

func (u *UserRepo) CreateUser(userData *requestmodels.UserSignUpRequest) error {
	query := "INSERT INTO users (name,user_name,email,password) VALUES ($1,$2,$3,$4)"
	err := u.DB.Exec(query, userData.Name, userData.UserName, userData.Email, userData.Password).Error
	if err != nil {
		fmt.Println("Error In Creating Users")
		return err
	}
	return nil
}

func (u *UserRepo) GetOtpInfo(email string) (string, time.Time, error) {
	var exp time.Time

	type OTPINFO struct {
		OTP        string    `gorm:"column:otp"`
		Expiration time.Time `gorm:"column:expiration"`
	}

	var otpInfo OTPINFO

	if err := u.DB.Raw("SELECT otp,expiration FROM otps WHERE email = ? ORDER BY expiration DESC LIMIT 1;", email).Scan(&otpInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", exp, errors.New("otp verification failed, no data found for this user-email")
		}
		return "", exp, errors.New("otp verification failed, error finding user data: " + err.Error())
	}

	return otpInfo.OTP, otpInfo.Expiration, nil
}

func (u *UserRepo) ActivateUser(email string) error {
	query := "UPDATE users SET status='active' WHERE email = $1"
	result := u.DB.Exec(query, email)
	if result.Error != nil {
		fmt.Println("Error In Activate Users : ", result.Error)
		return result.Error
	}
	return nil
}

func (u *UserRepo) GetUserIDByEmail(email string) (string, error) {
	var UserID string
	query := "SELECT id FROM users WHERE email=$1 And status=$2"
	err := u.DB.Raw(query, email, "active").Row().Scan(&UserID)
	if err != nil {
		fmt.Println("Error in Get user Id my email: ", err)
		return "", err
	}
	return UserID, nil
}

func (u *UserRepo) GetUserPasswordHashAndStatus(email string) (string, string, string, error) {
	var hashedPassword, status, userid string

	query := "SELECT password, id, status FROM users WHERE email=? AND status!='delete'"
	err := u.DB.Raw(query, email).Row().Scan(&hashedPassword, &userid, &status)
	if err != nil {
		return "", "", "", errors.New("no user exist with the specified email,signup first")
	}

	return hashedPassword, userid, status, nil
}

func (u *UserRepo) UpdateUserPassword(email, hashedPassword *string) error {
	query := `UPDATE users SET password=$1 WHERE email=$2`
	err := u.DB.Exec(query, *hashedPassword, *email).Error
	if err != nil {
		fmt.Println("Error In update User password")
		return err
	}
	return nil
}

func (u *UserRepo) GetUserStateForAccessTokenGeneration(userID *string) (*string, error) {
	var UserCurrentstatus string

	query := "SELECT status FROM users WHERE id=?"
	result := u.DB.Raw(query, userID).Scan(&UserCurrentstatus)

	if result.RowsAffected == 0 {
		return &UserCurrentstatus, errors.New("no result found this user id")
	}

	if result.Error != nil {
		fmt.Println("Error in Access Token Generation")
		return &UserCurrentstatus, result.Error
	}

	return &UserCurrentstatus, nil
}

func (u *UserRepo) GetUserLiteProfile(userID *string) (*responsemodels.UserProfileResponse, error) {
	var response responsemodels.UserProfileResponse

	query := "SELECT id,name,user_name,bio,links,profile_img_url FROM users WHERE id=$1"
	result := u.DB.Raw(query, userID).Scan(&response)

	if result.Error != nil {
		fmt.Println("Error in Get user Profile:", result.Error)
		return nil, result.Error
	}

	// Check if no rows were affected (i.e., user not found)
	if result.RowsAffected == 0 {
		return nil, errors.New("no user with specified user id")
	}

	return &response, nil
}

func (u *UserRepo) UpdateUserLiteProfile(editInput *requestmodels.EditUserProfileRequest) error {
	query := "UPDATE users SET name=$1,user_name=$2,bio=$3,links=$4 WHERE id=$5"
	err := u.DB.Exec(query, editInput.Name, editInput.UserName, editInput.Bio, editInput.Links, editInput.UserId).Error
	if err != nil {
		fmt.Println("Error Update User Profile")
		return err
	}
	return nil
}

func (u *UserRepo) GetUserProfileUrlAndUsername(userID *string) (*responsemodels.UserLiteResponse, error) {
	var response responsemodels.UserLiteResponse

	query := "SELECT user_name,profile_img_url FROM users WHERE id=$1"
	result := u.DB.Raw(query, userID).Scan(&response)

	// Check if there's an error during query execution
	if result.Error != nil {
		fmt.Println("Error in Get user Profile And Username:", result.Error)
		return nil, result.Error
	}

	// Check if no rows were affected (i.e., user not found)
	if result.RowsAffected == 0 {
		return nil, errors.New("no user with specified user id")
	}
	return &response, nil
}

func (u *UserRepo) UserExistsByID(userID string) (bool, *error) {
	var UserCount int

	query := "SELECT COUNT(*) FROM users WHERE id=$1 AND status != $2"
	err := u.DB.Raw(query, userID, "deleted").Row().Scan(&UserCount)

	if err != nil {
		fmt.Println("Error in user exit by id")
		return false, &err
	}

	if UserCount >= 1 {
		return true, nil
	}
	return false, nil
}

func (u *UserRepo) SearchUsersByName(myID, searchText, limit, offset *string) (*[]responsemodels.UserListResponse, error) {
	var response []responsemodels.UserListResponse

	query := "SELECT id,name,user_name,profile_img_url FROM users WHERE (name ILIKE $2 OR user_name ILIKE $1) AND status = 'active' AND id != $2 LIMIT $3 OFFSET $4"
	err := u.DB.Raw(query, "%"+*searchText+"%", myID, limit, offset).Scan(&response).Error
	if err != nil {
		fmt.Println("Error in search By UserName")
		return nil, err
	}
	return &response, nil
}

func (u *UserRepo) GetFollowersDetails(userIDs *[]uint64) (*[]responsemodels.UserListResponse, error) {
	var userData []responsemodels.UserListResponse

	// Check if userIDs is empty
	if len(*userIDs) == 0 {
		return &userData, nil // Return empty slice if no user IDs
	}

	// Prepare the query and arguments
	interfaceIds := make([]interface{}, len(*userIDs))
	for i, id := range *userIDs {
		interfaceIds[i] = id
	}

	// Construct the query
	query := "SELECT id, name, user_name, profile_img_url FROM users WHERE id IN ("
	for i := range *userIDs {
		query += "?"

		if i < len(*userIDs)-1 {
			query += ","
		}
	}
	query += ")"

	// Execute the query
	err := u.DB.Raw(query, interfaceIds...).Scan(&userData).Error
	if err != nil {
		fmt.Println("Error in Get Followers : ", err)
		return nil, err
	}
	return &userData, nil
}

func (u *UserRepo) GetFollowingDetails(userIDs *[]uint64) (*[]responsemodels.UserListResponse, error) {
	var userData []responsemodels.UserListResponse

	// Check if userIDs is empty
	if len(*userIDs) == 0 {
		return &userData, nil // Return empty slice if no user IDs
	}

	// Prepare the query and arguments
	interfaceIds := make([]interface{}, len(*userIDs))
	for i, id := range *userIDs {
		interfaceIds[i] = id
	}

	// Construct the query
	query := "SELECT id, name, user_name, profile_img_url FROM users WHERE id IN ("
	for i := range *userIDs {
		query += "?"

		if i < len(*userIDs)-1 {
			query += ","
		}
	}
	query += ")"

	// Execute the query
	err := u.DB.Raw(query, interfaceIds...).Scan(&userData).Error
	if err != nil {
		fmt.Println("Error in Get Following : ", err)
		return nil, err
	}
	return &userData, nil
}

func (u *UserRepo) SetUserProfileImage(userID, imageUrl *string) error {
	query := "UPDATE users SET profile_img_url=$1 WHERE id=$2"
	err := u.DB.Exec(query, imageUrl, userID).Error
	if err != nil {
		fmt.Println("Error in Set user Profile")
		return err
	}
	return nil
}

func (u *UserRepo) IsUserBlocked(userID string) (bool, error) {
	var status string

	// Query to check the user's status
	query := "SELECT status FROM users WHERE id=$1"
	err := u.DB.Raw(query, userID).Row().Scan(&status)

	if err != nil {
		fmt.Println("Error in checking user block status:", err)
		return false, err
	}

	// Assuming 'blocked' is the status representing a blocked user
	if status == "blocked" {
		return true, nil
	}

	return false, nil
}

