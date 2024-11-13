package server

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/models/requestmodels"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/pb"
	interface_usecase "github.com/ShahabazSulthan/Friendzy_Auth/pkg/usecase/interfaces"
)

type AuthService struct {
	UserUsecase    interface_usecase.IUserUseCase
	JwtUseCase     interface_usecase.IJwt
	AdminUsecase   interface_usecase.IAdminUsecase
	PaymentUsecase interface_usecase.IPaymentUsecase
	pb.AuthServiceServer
}

func NewAuthServer(userUseCase interface_usecase.IUserUseCase,
	jwtusecase interface_usecase.IJwt,
	adminusecase interface_usecase.IAdminUsecase,
	payemtUseacse interface_usecase.IPaymentUsecase) *AuthService {
	if jwtusecase == nil {
		log.Fatal("JwtUseCase cannot be nil in NewAuthServer")
	}
	return &AuthService{
		UserUsecase:    userUseCase,
		JwtUseCase:     jwtusecase,
		AdminUsecase:   adminusecase,
		PaymentUsecase: payemtUseacse,
	}
}

func (s *AuthService) UserSignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {
	var inputData requestmodels.UserSignUpRequest

	inputData.Name = req.Name
	inputData.UserName = req.UserName
	inputData.Email = req.Email
	inputData.Password = req.Password
	inputData.ConformPassword = req.ConfirmPassword

	respData, err := s.UserUsecase.UserSignUp(&inputData)
	if err != nil {
		return &pb.SignUpResponse{
			ErrorMessage: err.Error(),
		}, nil
	}

	fmt.Println("input", inputData)
	return &pb.SignUpResponse{
		Token: respData.Token,
	}, nil
}

func (s *AuthService) UserOTPVerification(ctx context.Context, req *pb.RequestOtpVefification) (*pb.ResponseOtpVerification, error) {
	respData, err := s.UserUsecase.VerifyOtp(req.Otp, &req.TempToken)
	if err != nil {
		return &pb.ResponseOtpVerification{
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.ResponseOtpVerification{
		Otp:          respData.OTP,
		AccessToken:  respData.AccessToken,
		RefreshToken: respData.RefereshToken,
	}, nil
}

func (s *AuthService) UserLogin(ctx context.Context, req *pb.RequestUserLogin) (*pb.ResponseUserLogin, error) {
	var loginData requestmodels.UserLoginRequest

	loginData.Email = req.Email
	loginData.Password = req.Password

	respData, err := s.UserUsecase.UserLogin(&loginData)
	if err != nil {
		return &pb.ResponseUserLogin{
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.ResponseUserLogin{
		AccessToken:  respData.AccessToken,
		RefreshToken: respData.RefereshToken,
	}, nil
}

func (s *AuthService) ForgotPasswordRequest(ctx context.Context, req *pb.RequestForgotPass) (*pb.ResponseForgotPass, error) {
	if req == nil {
		return nil, errors.New("ForgotPasswordRequest: request cannot be nil")
	}

	// Check if email is provided
	if req.Email == "" {
		return &pb.ResponseForgotPass{ErrorMessage: "Email is required"}, nil
	}

	// Process the ForgotPasswordRequest
	respData, err := s.UserUsecase.ForgetPasswordRequest(&req.Email)
	if err != nil {
		return &pb.ResponseForgotPass{ErrorMessage: err.Error()}, nil
	}

	// Check if the response data is not nil
	if respData == nil {
		return &pb.ResponseForgotPass{ErrorMessage: "Generated token is nil"}, nil
	}

	return &pb.ResponseForgotPass{Token: *respData}, nil
}

func (s *AuthService) ResetPassword(ctx context.Context, req *pb.RequestResetPass) (*pb.ResponseErrorMessage, error) {
	var requestData requestmodels.ForgotPasswordRequest

	requestData.OTP = req.Otp
	requestData.Password = req.Password
	requestData.ConformPassword = req.ConfirmPassword

	err := s.UserUsecase.ResetPassword(&requestData, &req.TempToken)

	if err != nil {
		return &pb.ResponseErrorMessage{
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.ResponseErrorMessage{}, nil

}

func (s *AuthService) VerifyAccessToken(ctx context.Context, req *pb.RequestVerifyAccess) (*pb.ResponseVerifyAccess, error) {
	uesrId, err := s.JwtUseCase.VerifyAccessToken(&req.AccessToken)

	if err != nil {
		return &pb.ResponseVerifyAccess{
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.ResponseVerifyAccess{
		UserId: *uesrId,
	}, nil
}

func (s *AuthService) GetEmailAndUsername(ctx context.Context, req *pb.GetEmailAndUsernameRequest) (*pb.GetEmailAndUsernameResponse, error) {
	email, username, err := s.UserUsecase.GetEmailAndUsername(int(req.UserId))
	if err != nil {
		fmt.Println("Error retrieving user email and username:", err)
		return &pb.GetEmailAndUsernameResponse{
			Email:    "",
			Username: "",
			Error:    err.Error(),
		}, nil
	}

	return &pb.GetEmailAndUsernameResponse{
		Email:    email,
		Username: username,
	}, nil
}

func (s *AuthService) AccessReGenerator(ctx context.Context, req *pb.RequestAccessGenerator) (*pb.ResponseAccessGenerator, error) {
	// Initial log to confirm function entry
	log.Println("Entering AccessReGenerator function")

	// Validate access and refresh tokens
	if req.GetAccessToken() == "" || req.GetRefreshToken() == "" {
		return nil, errors.New("access or refresh token is missing")
	}

	// Check if the request is nil
	if req == nil {
		log.Println("Request is nil")
		return &pb.ResponseAccessGenerator{ErrorMessage: "Request cannot be nil"}, nil
	}

	// Validate access token
	if req.AccessToken == "" {
		log.Println("Access token is missing")
		return &pb.ResponseAccessGenerator{ErrorMessage: "Access token is required"}, nil
	}

	// Validate refresh token
	if req.RefreshToken == "" {
		log.Println("Refresh token is missing")
		return &pb.ResponseAccessGenerator{ErrorMessage: "Refresh token is required"}, nil
	}

	// Ensure JwtUseCase is initialized
	if s.JwtUseCase == nil {
		log.Println("JwtUseCase is not initialized")
		return &pb.ResponseAccessGenerator{ErrorMessage: "JwtUseCase is not initialized"}, nil
	}

	// Log the request details
	log.Printf("Access Token: %s", req.AccessToken)
	log.Printf("Refresh Token: %s", req.RefreshToken)

	// Call JwtUseCase's AccessRegenerator method
	newAccessToken, err := s.JwtUseCase.AccessRegenerator(req.AccessToken, req.RefreshToken)
	if err != nil {
		log.Printf("AccessReGenerator error: %v", err)
		return &pb.ResponseAccessGenerator{ErrorMessage: "Failed to regenerate access token"}, nil
	}

	// Handle case where the new access token is nil or empty
	if newAccessToken == "" {
		log.Println("Generated access token is empty or nil")
		return &pb.ResponseAccessGenerator{ErrorMessage: "Generated access token is empty"}, nil
	}

	// Log the new access token
	log.Printf("New Access Token: %s", newAccessToken)

	// Return success with the new access token
	return &pb.ResponseAccessGenerator{
		AccesToken: newAccessToken,
	}, nil
}

func (s *AuthService) GetUserProfile(ctx context.Context, req *pb.RequestGetUserProfile) (*pb.ResponseUserProfile, error) {

	respData, err := s.UserUsecase.UserProfile(&req.UserId, &req.UserBId)
	if err != nil {
		return &pb.ResponseUserProfile{
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.ResponseUserProfile{
		Name:            respData.Name,
		UserName:        respData.UserName,
		Bio:             respData.Bio,
		Links:           respData.Links,
		ProfileImageURL: respData.UserProfileImgUrl,
		PostsCount:      uint64(respData.PostCount),
		FollowerCount:   uint64(respData.FollowCount),
		FollowingCount:  uint64(respData.FollowingCount),
		FollowingStat:   respData.FollowingStatus,
	}, nil

}

func (s *AuthService) EditUserProfile(ctx context.Context, req *pb.RequestEditUserProfile) (*pb.ResponseErrorMessage, error) {
	var editInput requestmodels.EditUserProfileRequest

	// Validate request
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}
	if req.UserId == "" {
		return &pb.ResponseErrorMessage{
			ErrorMessage: "User ID cannot be empty",
		}, nil
	}

	// Map fields from request
	editInput.Name = req.Name
	editInput.UserName = req.UserName
	editInput.Bio = req.Bio
	editInput.Links = req.Links
	editInput.UserId = req.UserId

	// Check if UserUsecase is properly initialized
	if s.UserUsecase == nil {
		return &pb.ResponseErrorMessage{
			ErrorMessage: "User use case is not initialized",
		}, nil
	}

	// Call the use case to edit user details
	err := s.UserUsecase.EditUserDetails(&editInput)
	if err != nil {
		return &pb.ResponseErrorMessage{
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.ResponseErrorMessage{}, nil
}

func (s *AuthService) CheckUserExist(ctx context.Context, req *pb.RequestUserId) (*pb.ResponseBool, error) {
	stat, err := s.UserUsecase.CheckUserExist(&req.UserId)

	if err != nil {
		return &pb.ResponseBool{
			ErrorMessage: (*err).Error(),
		}, nil
	}

	return &pb.ResponseBool{
		ExistStatus: stat,
	}, nil
}

func (s *AuthService) GetUserDetailsLiteForPostView(ctc context.Context, req *pb.RequestUserId) (*pb.ResponseUserDetailsLite, error) {

	respData, err := s.UserUsecase.GetUserDetailsLiteForPostView(&req.UserId)
	if err != nil {
		return &pb.ResponseUserDetailsLite{
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.ResponseUserDetailsLite{
		UserName:          respData.UserName,
		UserProfileImgURL: respData.UserProfileImgUrl,
	}, nil
}

func (s *AuthService) SearchUser(ctx context.Context, req *pb.RequestUserSearch) (*pb.ResponseUserSearch, error) {
	respData, err := s.UserUsecase.SearchUser(&req.UserId, &req.SearchText, &req.Limit, &req.Offset)
	if err != nil {
		return &pb.ResponseUserSearch{
			Errormessage: err.Error(),
		}, nil
	}

	var respSlice []*pb.SingleResponseGetFollowers
	for i := range *respData {
		respSlice = append(respSlice, &pb.SingleResponseGetFollowers{
			UserId:        (*respData)[i].Id,
			Name:          (*respData)[i].Name,
			UserName:      (*respData)[i].UserName,
			ProfileImgUrl: (*respData)[i].UserProfileImgUrl})
	}

	return &pb.ResponseUserSearch{
		SearchResult: respSlice,
	}, nil
}

func (s *AuthService) GetFollowersDetails(ctx context.Context, req *pb.RequestUserId) (*pb.ResponseGetUsersDetails, error) {
	respData, err := s.UserUsecase.GetFollowersDetails(&req.UserId)
	if err != nil {
		return &pb.ResponseGetUsersDetails{
			ErrorMessage: (*err).Error(),
		}, nil
	}

	var respLoader []*pb.SingleResponseGetFollowers
	for i := range *respData {
		respLoader = append(respLoader, &pb.SingleResponseGetFollowers{
			UserId:        (*respData)[i].Id,
			Name:          (*respData)[i].Name,
			UserName:      (*respData)[i].UserName,
			ProfileImgUrl: (*respData)[i].UserProfileImgUrl})

	}
	return &pb.ResponseGetUsersDetails{
		UserData: respLoader,
	}, nil
}

func (s *AuthService) GetFollowingsDetails(ctx context.Context, req *pb.RequestUserId) (*pb.ResponseGetUsersDetails, error) {
	respData, err := s.UserUsecase.GetFollowingDetails(&req.UserId)
	if err != nil {
		return &pb.ResponseGetUsersDetails{
			ErrorMessage: (*err).Error(),
		}, nil
	}

	var respLoader []*pb.SingleResponseGetFollowers
	for i := range *respData {
		respLoader = append(respLoader, &pb.SingleResponseGetFollowers{
			UserId:        (*respData)[i].Id,
			Name:          (*respData)[i].Name,
			UserName:      (*respData)[i].UserName,
			ProfileImgUrl: (*respData)[i].UserProfileImgUrl})

	}
	return &pb.ResponseGetUsersDetails{
		UserData: respLoader,
	}, nil
}

func (s *AuthService) SetUserProfileImage(ctx context.Context, req *pb.RequestSetProfileImg) (*pb.ResponseErrorMessage, error) {

	err := s.UserUsecase.SetUserProfileImage(&req.UserId, &req.ContentType, &req.Img)
	if err != nil {
		return &pb.ResponseErrorMessage{
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.ResponseErrorMessage{}, nil
}

func (s *AuthService) AdminLogin(ctx context.Context, req *pb.AdminLoginRequest) (*pb.AdminLoginResponse, error) {
	// Input validation
	if req.Email == "" || req.Password == "" {
		return &pb.AdminLoginResponse{
			ErrorMessage: "Email and Password are required",
		}, nil
	}

	// Prepare the login request data
	adminLoginData := &requestmodels.AdminLoginData{
		Email:    req.Email,
		Password: req.Password,
	}

	// Call AdminUsecase to process the login
	respData, err := s.AdminUsecase.AdminLogin(adminLoginData)
	if err != nil {
		return &pb.AdminLoginResponse{
			ErrorMessage: err.Error(),
		}, nil
	}

	// Return success response with token
	return &pb.AdminLoginResponse{
		Token: respData.Token,
	}, nil
}

// GetAllUsers fetches a list of all users with pagination support
func (s *AuthService) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	// Validate request parameters
	if req.Limit == "" || req.Offset == "" {
		return &pb.GetAllUsersResponse{
			ErrorMessage: "Limit and Offset are required",
		}, nil
	}

	// Call the use case to get all users
	usersData, err := s.AdminUsecase.GetAllUsers(req.Limit, req.Offset)
	if err != nil {
		return &pb.GetAllUsersResponse{
			ErrorMessage: err.Error(),
		}, nil
	}

	// Prepare the response
	var userResponses []*pb.UserAdminResponse
	for _, user := range *usersData {
		userResponses = append(userResponses, &pb.UserAdminResponse{
			ID:              uint64(user.ID),
			Name:            user.Name,
			UserName:        user.UserName,
			Email:           user.Email,
			Bio:             user.Bio,
			ProfileImageURL: user.ProfileImageURL,
			Links:           user.Links,
			Status:          user.Status,
		})
	}

	return &pb.GetAllUsersResponse{
		Users: userResponses,
	}, nil
}

// BlockUser blocks a user by their ID
func (s *AuthService) BlockUser(ctx context.Context, req *pb.BlockUserRequest) (*pb.ResponseErrorMessage, error) {
	// Validate UserID
	if req.UserId == "" {
		return &pb.ResponseErrorMessage{
			ErrorMessage: "UserId is required",
		}, nil
	}

	// Call the use case to block the user
	err := s.AdminUsecase.BlcokUser(req.UserId)
	if err != nil {
		return &pb.ResponseErrorMessage{
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.ResponseErrorMessage{}, nil
}

// UnblockUser unblocks a user by their ID
func (s *AuthService) UnblockUser(ctx context.Context, req *pb.UnblockUserRequest) (*pb.ResponseErrorMessage, error) {
	// Validate UserID
	if req.UserId == "" {
		return &pb.ResponseErrorMessage{
			ErrorMessage: "UserId is required",
		}, nil
	}

	// Call the use case to unblock the user
	err := s.AdminUsecase.UnblockUser(req.UserId)
	if err != nil {
		return &pb.ResponseErrorMessage{
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.ResponseErrorMessage{}, nil
}

func (s *AuthService) VerifyAdminToken(ctx context.Context, req *pb.RequestVerifyAdmin) (*pb.ResponseVerifyAdmin, error) {
	adminEmail, err := s.JwtUseCase.VerifyAdminToken(&req.AdminToken)

	if err != nil {
		return &pb.ResponseVerifyAdmin{
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.ResponseVerifyAdmin{
		AdminEmail:   *adminEmail, // Assuming adminEmail returns a pointer to a string
		ErrorMessage: "",          // No error, so we set this to an empty string
	}, nil
}

func (s *AuthService) CreateBlueTickPayment(ctx context.Context, req *pb.CreateBlueTickPaymentRequest) (*pb.CreateBlueTickPaymentResponse, error) {
	// Call the use case to create a Razorpay order and get the verification ID
	verificationID, err := s.PaymentUsecase.CreateBlueTickPayment(uint(req.UserId))
	if err != nil {
		log.Printf("Error creating blue tick payment: %v", err)
		return &pb.CreateBlueTickPaymentResponse{
			Message: err.Error(),
		}, nil
	}

	// Return the successful response
	return &pb.CreateBlueTickPaymentResponse{
		VerificationId: verificationID,
		Message:        "Payment order created successfully",
	}, nil
}

func (s *AuthService) VerifyBlueTickPayment(ctx context.Context, req *pb.VerifyBlueTickPaymentRequest) (*pb.VerifyBlueTickPaymentResponse, error) {
	// Call the Razorpay use case to verify the payment
	isVerified, err := s.PaymentUsecase.VerifyBlueTickPayment(req.GetVerificationId(), req.GetPaymentId(), req.GetSignature(), uint(req.UserId))
	if err != nil {
		return &pb.VerifyBlueTickPaymentResponse{
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}

	// Return the success response if verification is successful
	return &pb.VerifyBlueTickPaymentResponse{
		Success:      isVerified,
		ErrorMessage: "",
	}, nil
}

// GetBlueTickVerification handles the gRPC request to fetch blue tick verification details.
func (s *AuthService) OnlinePayment(ctx context.Context, req *pb.OnlinePaymentRequest) (*pb.OnlinePaymentResponse, error) {
	// Call the use case to get blue tick verification details
	_, err := s.PaymentUsecase.OnlinePayment(req.GetUserId(), req.GetVerificationId())
	if err != nil {
		// Log error and return a failed response
		return &pb.OnlinePaymentResponse{
			UserId:          req.UserId,
			Paymentstatus:   "Failed",
			Verificationfee: "1000",
		}, nil
	}

	// Return the success response with the verification details
	return &pb.OnlinePaymentResponse{
		UserId:          req.UserId,
		Paymentstatus:   "Success",
		Verificationfee: "1000",
	}, nil
}

func (s *AuthService) GetAllVerifiedUsers(ctx context.Context, req *pb.GetAllVerifiedUsersRequest) (*pb.GetAllverifiedUsers, error) {
	// Validate request parameters
	if req.Limit == "" || req.Offset == "" {
		return &pb.GetAllverifiedUsers{
			ErrorMessage: "Limit and Offset are required",
		}, nil
	}

	// Call the use case to get all users
	usersData, err := s.PaymentUsecase.GetAllVerifiedUsers(req.Limit, req.Offset)
	if err != nil {
		return &pb.GetAllverifiedUsers{
			ErrorMessage: err.Error(),
		}, nil
	}

	// Prepare the response
	var userResponses []*pb.BlueTickResponse
	for _, user := range *usersData {
		userResponses = append(userResponses, &pb.BlueTickResponse{
			ID:              uint64(user.ID),
			BlueTick:        "☑️",
			Name:            user.Name,
			UserName:        user.UserName,
			Email:           user.Email,
			Bio:             user.Bio,
			ProfileImageURL: user.ProfileImageURL,
			Links:           user.Links,
			Status:          user.Status,
		})
	}

	return &pb.GetAllverifiedUsers{
		Users: userResponses,
	}, nil
}

func (s *AuthService) CheckUserVerified(ctx context.Context, req *pb.RequestUserId) (*pb.ResponseBool, error) {
	stat, err := s.PaymentUsecase.IsUserVerified(req.UserId)
	if err != nil {
		return &pb.ResponseBool{
			ExistStatus:  false,
			ErrorMessage: "Error in verifying user status",
		}, err // Returning the actual error so it can be handled by the caller
	}

	return &pb.ResponseBool{
		ExistStatus:  stat,
		ErrorMessage: "",
	}, nil
}
