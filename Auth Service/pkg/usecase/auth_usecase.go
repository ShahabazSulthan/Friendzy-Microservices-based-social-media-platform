package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/config"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/models/requestmodels"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/models/responsemodels"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/pb"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/repository/interfaces"
	interface_usecase "github.com/ShahabazSulthan/Friendzy_Auth/pkg/usecase/interfaces"
	interface_jwt "github.com/ShahabazSulthan/Friendzy_Auth/pkg/utils/JWT/Interface"
	interface_hash "github.com/ShahabazSulthan/Friendzy_Auth/pkg/utils/hashed_password/interfaces"
	interface_random "github.com/ShahabazSulthan/Friendzy_Auth/pkg/utils/randomNumber/interface"
	interface_regex "github.com/ShahabazSulthan/Friendzy_Auth/pkg/utils/regex/interface"
	interface_smtp "github.com/ShahabazSulthan/Friendzy_Auth/pkg/utils/smtp/interface"
)

type UserUseCase struct {
	UserRepo         interfaces.IUserRepository
	PaymentRepo      interfaces.IPaymentRepository
	SmtpUtil         interface_smtp.Ismtp
	JwtUtil          interface_jwt.Ijwt
	RandomNumUtil    interface_random.IRandGene
	RegexUtil        interface_regex.IRegex
	TokenSecurityKey *config.Token
	HashUtil         interface_hash.IHashPassword
	PostANdClient    pb.PostNrelServiceClient
	Razopay          *config.Razopay
}

func NewUserCase(userRepo interfaces.IUserRepository,
	paymentRepo interfaces.IPaymentRepository,
	smtpUtil interface_smtp.Ismtp,
	jwtUtil interface_jwt.Ijwt,
	randNumUtil interface_random.IRandGene,
	regexUtil interface_regex.IRegex,
	config *config.Token,
	hashUtil interface_hash.IHashPassword,
	postNrel *pb.PostNrelServiceClient,
	razopay *config.Razopay) interface_usecase.IUserUseCase {
	return &UserUseCase{
		UserRepo:         userRepo,
		PaymentRepo:      paymentRepo,
		SmtpUtil:         smtpUtil,
		JwtUtil:          jwtUtil,
		RandomNumUtil:    randNumUtil,
		RegexUtil:        regexUtil,
		TokenSecurityKey: config,
		HashUtil:         hashUtil,
		PostANdClient:    *postNrel,
		Razopay:          razopay,
	}
}

func (u *UserUseCase) UserSignUp(userData *requestmodels.UserSignUpRequest) (responsemodels.SignUpResponse, error) {
	var resSignUp responsemodels.SignUpResponse

	fmt.Println("UserData 1 -----", userData)
	if isUserExist := u.UserRepo.UserExistsByEmail(userData.Email); isUserExist {
		fmt.Println("Error in email exist")
		return resSignUp, errors.New("user exists, try agaim another email")
	}

	stat, message := u.RegexUtil.IsValidPassword(userData.Password)
	if !stat {
		fmt.Println("Error in Password")
		return resSignUp, errors.New(message)
	}

	stat, message = u.RegexUtil.IsValidUsername(userData.UserName)
	if !stat {
		fmt.Println("Error Validate username")
		return resSignUp, errors.New(message)
	}

	fmt.Println("")
	if isUserExistUserName := u.UserRepo.UserExistsByUsername(userData.UserName); isUserExistUserName {
		return resSignUp, errors.New("user name exist, try another username")
	}

	errRemv := u.UserRepo.DeleteRecentOtpRequests()
	if errRemv != nil {
		return resSignUp, errRemv
	}

	otp := u.RandomNumUtil.RandomNumber()
	errOtp := u.SmtpUtil.SendNotificationWithEmailOtp(otp, userData.Email, userData.Name)
	if errOtp != nil {
		return resSignUp, errOtp
	}

	exp := time.Now().Add(5 * time.Minute)

	errTempServe := u.UserRepo.SaveUserOtp(otp, userData.Email, exp)
	if errTempServe != nil {
		fmt.Println("Cant Save temporary data for otp varification in db")
		return resSignUp, errors.New("otp verification down,please try after some time")
	}

	hashedPassword := u.HashUtil.HashedPassword(userData.ConformPassword)
	userData.Password = hashedPassword

	errCreateUser := u.UserRepo.CreateUser(userData)
	if errCreateUser != nil {
		return resSignUp, errCreateUser
	}

	tempToken, err := u.JwtUtil.TempTokenForOtpVerification(u.TokenSecurityKey.TempVerificationKey, userData.Email)
	if err != nil {
		fmt.Println("error creating temp token for otp verification")
		return resSignUp, errors.New("error creating temp token for otp verification")
	}

	resSignUp.Token = tempToken

	return resSignUp, nil
}

func (u *UserUseCase) VerifyOtp(otp string, TempVerificationToken *string) (responsemodels.OTPVerificationResponse, error) {
	var otpverificationres responsemodels.OTPVerificationResponse

	email, unbindErr := u.JwtUtil.UnbindEmailFromClaim(*TempVerificationToken, u.TokenSecurityKey.TempVerificationKey)
	if unbindErr != nil {
		return otpverificationres, unbindErr
	}

	userOtp, exp, errGetinfo := u.UserRepo.GetOtpInfo(email)
	if errGetinfo != nil {
		return otpverificationres, errGetinfo
	}

	if otp != userOtp {
		return otpverificationres, errors.New("invalid otp")
	}

	if time.Now().After(exp) {
		return otpverificationres, errors.New("otp expired")
	}

	changeStatErr := u.UserRepo.ActivateUser(email)
	if changeStatErr != nil {
		return otpverificationres, changeStatErr
	}

	userId, fetchErr := u.UserRepo.GetUserIDByEmail(email)
	if fetchErr != nil {
		return otpverificationres, fetchErr
	}

	accessToken, atokenErr := u.JwtUtil.GenerateAccessToken(u.TokenSecurityKey.UserSecurityKey, userId)
	if atokenErr != nil {
		otpverificationres.AccessToken = atokenErr.Error()
		return otpverificationres, atokenErr
	}

	refreshToken, rtokenErr := u.JwtUtil.GenerateRefreshToken(u.TokenSecurityKey.UserSecurityKey)
	if rtokenErr != nil {
		otpverificationres.RefereshToken = rtokenErr.Error()
		return otpverificationres, rtokenErr
	}

	otpverificationres.OTP = "verified"
	otpverificationres.AccessToken = accessToken
	otpverificationres.RefereshToken = refreshToken

	return otpverificationres, nil
}

func (u *UserUseCase) UserLogin(loginData *requestmodels.UserLoginRequest) (responsemodels.LoginResponse, error) {
	var resLogin responsemodels.LoginResponse

	stat, message := u.RegexUtil.IsValidPassword(loginData.Password)
	if !stat {
		return resLogin, errors.New(message)
	}

	hashedPassword, userId, status, errs := u.UserRepo.GetUserPasswordHashAndStatus(loginData.Email)
	if errs != nil {
		return resLogin, errs
	}

	passwordErr := u.HashUtil.ComparePassword(hashedPassword, loginData.Password)
	if passwordErr != nil {
		return resLogin, passwordErr
	}

	if status == "blocked" {
		return resLogin, errors.New("user is blocked by admin")
	}

	if status == "pending" {
		return resLogin, errors.New("user is pending,otp is not verified")
	}

	accessToken, err := u.JwtUtil.GenerateAccessToken(u.TokenSecurityKey.UserSecurityKey, userId)
	if err != nil {
		return resLogin, err
	}

	refreshToken, err := u.JwtUtil.GenerateRefreshToken(u.TokenSecurityKey.UserSecurityKey)
	if err != nil {
		return resLogin, err
	}

	resLogin.AccessToken = accessToken
	resLogin.RefereshToken = refreshToken
	return resLogin, nil
}

func (r *UserUseCase) ForgetPasswordRequest(email *string) (*string, error) {

	_, _, status, err := r.UserRepo.GetUserPasswordHashAndStatus(*email)
	if err != nil {
		return nil, err
	}

	if status == "blocked" {
		return nil, errors.New("user is blocked by the admin")
	}

	if status == "pending" {
		return nil, errors.New("user is on status pending,OTP not verified")
	}

	err = r.UserRepo.DeleteRecentOtpRequests()
	if err != nil {
		return nil, err
	}

	otp := r.RandomNumUtil.RandomNumber()
	err = r.SmtpUtil.SendRestPasswordEmailOtp(otp, *email)
	if err != nil {
		return nil, err
	}

	expiration := time.Now().Add(5 * time.Minute)

	errTempSave := r.UserRepo.SaveUserOtp(otp, *email, expiration)
	if errTempSave != nil {
		fmt.Println("Cant save temporary data for otp verification in db")
		return nil, errors.New("OTP verification down,please try after some time")
	}

	tempToken, err := r.JwtUtil.TempTokenForOtpVerification(r.TokenSecurityKey.TempVerificationKey, *email)
	if err != nil {
		fmt.Println("----------", err)
		return nil, errors.New("error creating temp token for otp verification")
	}

	return &tempToken, nil
}

func (u *UserUseCase) ResetPassword(userData *requestmodels.ForgotPasswordRequest, TempVerification *string) error {
	stat, message := u.RegexUtil.IsValidPassword(userData.Password)
	if !stat {
		return errors.New(message)
	}

	email, err := u.JwtUtil.UnbindEmailFromClaim(*TempVerification, u.TokenSecurityKey.TempVerificationKey)
	if err != nil {
		return err
	}

	userOTP, expiration, err := u.UserRepo.GetOtpInfo(email)
	if err != nil {
		return err
	}

	if userData.OTP != userOTP {
		return errors.New("invalid OTP")
	}
	if time.Now().After(expiration) {
		return errors.New("OTP expired")
	}

	hashedPassword := u.HashUtil.HashedPassword(userData.ConformPassword)

	err = u.UserRepo.UpdateUserPassword(&email, &hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserUseCase) UserProfile(userId, UserBId *string) (*responsemodels.UserProfileResponse, error) {
	var actualId *string

	if *UserBId == "" {
		actualId = userId
	} else {
		actualId = UserBId
	}
	userData, err := u.UserRepo.GetUserLiteProfile(actualId)
	if err != nil {
		return nil, err
	}

	isVerify, _ := u.PaymentRepo.IsUserVerified(*actualId)
	if isVerify {
		userData.BlueTickVerified = "☑️"
	}

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	respData, err := u.PostANdClient.GetCountsForUserProfile(context, &pb.RequestUserIdPnR{
		UserId: *actualId,
	})
	if err != nil {
		log.Fatal(err)
	}
	if respData.ErrorMessage != "" {
		return nil, errors.New(respData.ErrorMessage)
	}

	userData.PostCount = uint(respData.PostCount)
	userData.FollowCount = uint(respData.FollowerCount)
	userData.FollowingCount = uint(respData.FollowingCount)

	if *UserBId != "" {

		respStat, err := u.PostANdClient.UserAFollowingUserBorNot(context, &pb.RequestFollowUnFollow{
			UserId:  *userId,
			UserBId: *UserBId,
		})
		if err != nil {
			log.Fatal(err)
		}
		if respData.ErrorMessage != "" {
			return nil, errors.New(respData.ErrorMessage)
		}
		userData.FollowingStatus = respStat.BoolStat
	}

	return userData, nil
}

func (u *UserUseCase) EditUserDetails(editInput *requestmodels.EditUserProfileRequest) error {
	stat, message := u.RegexUtil.IsValidUsername(editInput.UserName)
	if !stat {
		return errors.New(message)
	}

	userData, err := u.UserRepo.GetUserLiteProfile(&editInput.UserId)
	if err != nil {
		fmt.Println("Error in edit password : ", err)
		return err
	}

	if userData.UserName != editInput.UserName {
		if isUserExistUserName := u.UserRepo.UserExistsByUsername(editInput.UserName); isUserExistUserName {
			return errors.New("user exist, try again with another username")
		}
	}

	err = u.UserRepo.UpdateUserLiteProfile(editInput)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserUseCase) GetEmailAndUsername(userID int) (string, string, error) {
	// Call the repository function to get the email and username
	email, username, err := u.UserRepo.GetEmailAndUsernameByUserID(userID)
	if err != nil {
		fmt.Println("Error retrieving user email and username:", err)
		return "", "", err
	}

	return email, username, nil
}

func (u *UserUseCase) GetUserDetailsLiteForPostView(userId *string) (*responsemodels.UserLiteResponse, error) {
	respData, err := u.UserRepo.GetUserProfileUrlAndUsername(userId)

	if err != nil {
		return respData, err
	}

	return respData, nil
}

func (u *UserUseCase) CheckUserExist(userId *string) (bool, *error) {
	boolstat, err := u.UserRepo.UserExistsByID(*userId)
	if err != nil {
		return boolstat, err
	}
	return boolstat, nil
}

func (u *UserUseCase) SearchUser(myId, SearchText, limit, offset *string) (*[]responsemodels.UserListResponse, error) {

	respData, err := u.UserRepo.SearchUsersByName(myId, SearchText, limit, offset)
	if err != nil {
		return nil, err
	}
	return respData, nil
}

func (u *UserUseCase) GetFollowersDetails(userId *string) (*[]responsemodels.UserListResponse, *error) {
	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	userIdsSlice, err := u.PostANdClient.GetFollowersIds(context, &pb.RequestUserIdPnR{UserId: *userId})
	if err != nil {
		log.Fatal(err)
	}
	if userIdsSlice.ErrorMessage != "" {
		return nil, &err
	}

	userDetailsSlice, err := u.UserRepo.GetFollowersDetails(&userIdsSlice.UserIds)
	if err != nil {
		return nil, &err
	}

	return userDetailsSlice, nil
}

func (u *UserUseCase) GetFollowingDetails(userId *string) (*[]responsemodels.UserListResponse, *error) {
	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	userIdsSlice, err := u.PostANdClient.GetFollowingsIds(context, &pb.RequestUserIdPnR{UserId: *userId})
	if err != nil {
		log.Fatal(err)
	}
	if userIdsSlice.ErrorMessage != "" {
		return nil, &err
	}

	userDetailsSlice, err := u.UserRepo.GetFollowingDetails(&userIdsSlice.UserIds)
	if err != nil {
		return nil, &err
	}

	return userDetailsSlice, nil
}

func (u *UserUseCase) SetUserProfileImage(userId, contentType *string, Img *[]byte) error {
	// Define the folder where images will be saved locally
	folderPath := "AuthService/userprofileimg/"

	// Ensure the directory exists, create it if not
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return err
	}

	// Generate a unique filename for the image
	fileName := fmt.Sprintf("%s_%d.%s", *userId, time.Now().Unix(), getExtension(*contentType))

	// Define the full path where the file will be saved
	filePath := filepath.Join(folderPath, fileName)

	// Write the image to the local folder
	err = os.WriteFile(filePath, *Img, 0644) // Save file with read/write permissions
	if err != nil {
		fmt.Printf("Error saving file locally: %v\n", err)
		return err
	}

	// You can return the local file path or some relative URL
	mediaURL := filePath // or serve via a web server as URL

	// Update user profile image URL in the repository
	err = u.UserRepo.SetUserProfileImage(userId, &mediaURL)
	if err != nil {
		return err
	}

	return nil
}

// Helper function to get the file extension based on content type
func getExtension(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return "jpg"
	case "image/png":
		return "png"
	default:
		return "bin" // Default to a generic binary extension
	}
}

//==============================================================================//
