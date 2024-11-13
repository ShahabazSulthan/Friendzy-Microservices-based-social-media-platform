package handler_auth

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"time"

	requestmodel_auth "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Auth_Service/models/requestModel"
	responsemodel_auth "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Auth_Service/models/responseModel"
	"github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Auth_Service/pb"
	byteconverter "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Utils/byteConverter"
	mediafileformatchecker "github.com/ShahabazSulthan/Friendzy_apiGateway/pkg/Utils/mediaFileFormatChecker"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	Client pb.AuthServiceClient
}

func NewAuthUserHandler(client *pb.AuthServiceClient) *UserHandler {
	return &UserHandler{Client: *client}
}

func (svc *UserHandler) UserSignUp(ctx *fiber.Ctx) error {

	var userSignupData requestmodel_auth.UserSignUpReq
	var resSignUp responsemodel_auth.SignupData

	if err := ctx.BodyParser(&userSignupData); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "signup failed(possible-reason:no json input)",
				Error:      err.Error(),
			})
	}

	fmt.Println("User SignbUp data = ", userSignupData)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(userSignupData)
	if err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, e := range ve {
				switch e.Field() {
				case "Name":
					resSignUp.Name = "should be a valid Name. "
				case "UserName":
					resSignUp.UserName = "should be a valid username. "
				case "Email":
					resSignUp.Email = "should be a valid email address. "
				case "Password":
					resSignUp.Password = "Password should have four or more digit"
				case "ConfirmPassword":
					resSignUp.ConfirmPassword = "should match the first password"
				}
			}
		}
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "signup failed",
				Data:       resSignUp,
				Error:      "did't fullfill the signup requirement ",
			})
	}

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.UserSignUp(context, &pb.SignUpRequest{
		Name:            userSignupData.Name,
		UserName:        userSignupData.UserName,
		Email:           userSignupData.Email,
		Password:        userSignupData.Password,
		ConfirmPassword: userSignupData.ConfirmPassword,
	})

	if err != nil {
		fmt.Println("----------auth service down--------")

		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "signup failed",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "signup failed",
				Data:       resp,
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_auth.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "signup success",
			Data:       resp,
			Error:      nil,
		})
}

func (svc *UserHandler) UserOTPVerication(ctx *fiber.Ctx) error {

	var otpData requestmodel_auth.OtpVerification
	var otpveriRes responsemodel_auth.OtpVerifResult

	temptoken := ctx.Get("x-temp-token")

	if err := ctx.BodyParser(&otpData); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "OTP verification failed(possible-reason:no json input)",
				Error:      err.Error(),
			})
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(otpData)
	if err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, e := range ve {
				switch e.Field() {
				case "Otp":
					otpData.Otp = "otp should be a 4 digit number"
				}
			}
		}
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "OTP verification failed",
				Data:       otpveriRes,
				Error:      otpveriRes.Otp,
			})
	}

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.UserOTPVerification(context, &pb.RequestOtpVefification{
		TempToken: temptoken,
		Otp:       otpData.Otp,
	})

	if err != nil {
		fmt.Println("----------auth service down--------")

		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "OTP verification failed",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "OTP verification failed",
				Data:       resp,
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_auth.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "OTP verification success",
			Data:       resp,
			Error:      nil,
		})
}

func (svc *UserHandler) UserLogin(ctx *fiber.Ctx) error {

	var loginData requestmodel_auth.UserLoginReq
	var resLogin responsemodel_auth.UserLoginRes

	if err := ctx.BodyParser(&loginData); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "login failed(possible-reason:no json input)",
				Error:      err.Error(),
			})
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(loginData)
	if err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, e := range ve {
				switch e.Field() {
				case "Email":
					resLogin.Email = "Enter a valid email"
				case "Password":
					resLogin.Password = "Password should have four or more digit"
				}
			}
		}
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "User Login verification failed",
				Data:       resLogin,
				Error:      err.Error(),
			})
	}

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.UserLogin(context, &pb.RequestUserLogin{
		Email:    loginData.Email,
		Password: loginData.Password,
	})

	if err != nil {
		fmt.Println("----------auth service down--------")

		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "login failed",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "login failed",
				Data:       resp,
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_auth.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "login success",
			Data:       resp,
			Error:      nil,
		})
}

func (svc *UserHandler) ForgotPasswordRequest(ctx *fiber.Ctx) error {

	var forgotReqData requestmodel_auth.ForgotPasswordReq
	var resData responsemodel_auth.ForgotPasswordRes

	if err := ctx.BodyParser(&forgotReqData); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "failed request(possible-reason:no json input)",
				Error:      err.Error(),
			})
	}

	fmt.Println("req = ", forgotReqData)
	fmt.Println("res = ", resData)

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(forgotReqData); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			if len(ve) > 0 && ve[0].Field() == "Email" {
				resData.Email = "Enter a valid email"
			}
		}
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "failed request err1",
				Data:       resData,
				Error:      err.Error(),
			})
	}

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.ForgotPasswordRequest(context, &pb.RequestForgotPass{
		Email: forgotReqData.Email,
	})

	if err != nil {
		fmt.Println("----------auth service down--------")

		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "failed request errr",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "failed request",
				Data:       resp,
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_auth.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "success",
			Data:       resp,
			Error:      nil,
		})

}

func (svc *UserHandler) ResetPassword(ctx *fiber.Ctx) error {

	// Get the temporary token from the request header
	temptoken := ctx.Get("x-temp-token")
	if temptoken == "" {
		return ctx.Status(fiber.StatusUnauthorized).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusUnauthorized,
				Message:    "No temporary token provided",
				Error:      "temp token is missing",
			})
	}

	fmt.Println("Received temp token: ", temptoken)

	var requestData requestmodel_auth.ForgotPasswordData
	var resData responsemodel_auth.ForgotPasswordData

	// Parse the request body
	if err := ctx.BodyParser(&requestData); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "Failed request (possible reason: no JSON input)",
				Error:      err.Error(),
			})
	}

	// Validate the request data
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(requestData)
	if err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, e := range ve {
				switch e.Field() {
				case "Otp":
					resData.Otp = "OTP should be a 4-digit number"
				case "Password":
					resData.Password = "Password should have four or more digits"
				case "ConfirmPassword":
					resData.ConfirmPassword = "Passwords should match"
				}
			}
		}
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "Failed to reset password",
				Data:       resData,
				Error:      err.Error(),
			})
	}

	// Create a context with a timeout
	c, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Call the ResetPassword function on the client (gRPC service)
	resp, err := svc.Client.ResetPassword(c, &pb.RequestResetPass{
		Otp:             requestData.Otp,
		Password:        requestData.Password,
		ConfirmPassword: requestData.ConfirmPassword,
		TempToken:       temptoken,
	})

	// Handle service unavailability or errors
	if err != nil {
		fmt.Println("----------auth service down--------")
		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "Failed to reset password",
				Error:      err.Error(),
			})
	}

	// Check if the response contains an error message
	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "Failed to reset password",
				Data:       resp,
				Error:      resp.ErrorMessage,
			})
	}

	// Return a success response
	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_auth.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "Password reset successfully",
			Data:       resp,
			Error:      nil,
		})
}

func (svc *UserHandler) GetUserProfile(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId")

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.GetUserProfile(context, &pb.RequestGetUserProfile{
		UserId: fmt.Sprint(userId),
	})

	if err != nil {
		fmt.Println("----------auth service down--------")

		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "failed to get user profile",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "failed to get user profile",
				Data:       resp,
				Error:      resp.ErrorMessage,
			})
	}

	var respStruct responsemodel_auth.UserProfileA //used to show the zero count of posts,following,followers etc

	intValueuserId, _ := strconv.Atoi(fmt.Sprint(userId))
	uintValueuserId := uint(intValueuserId)

	respStruct.UserId = uintValueuserId
	respStruct.Name = resp.Name
	respStruct.UserName = resp.UserName
	respStruct.Bio = resp.Bio
	respStruct.Links = resp.Links
	respStruct.UserProfileImgURL = resp.ProfileImageURL
	respStruct.PostsCount = uint(resp.PostsCount)
	respStruct.FollowersCount = uint(resp.FollowerCount)
	respStruct.FollowingCount = uint(resp.FollowingCount)

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_auth.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "fetched user profile successfully",
			Data:       respStruct,
			Error:      nil,
		})
}

func (svc *UserHandler) EditUserProfile(ctx *fiber.Ctx) error {

	var editInput requestmodel_auth.EditUserProfile
	var respEditUsr responsemodel_auth.EditUserProfileResp

	userId := ctx.Locals("userId")
	editInput.UserId = fmt.Sprint(userId)

	if err := ctx.BodyParser(&editInput); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "failed request(possible-reason:no json input)",
				Error:      err.Error(),
			})
	}

	fmt.Println("profile res = ", editInput)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(editInput)
	if err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, e := range ve {
				switch e.Field() {
				case "Name":
					respEditUsr.Name = "should be a valid Name. "
				case "UserName":
					respEditUsr.UserName = "should be a valid username. "
				case "Bio":
					respEditUsr.Bio = "Bio can't exceed 60 characters "
				case "Links":
					respEditUsr.Links = "Links can't exceed 20 characters"
				}
			}
		}
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "can't edit user details",
				Data:       respEditUsr,
				Error:      err.Error(),
			})
	}

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.EditUserProfile(context, &pb.RequestEditUserProfile{
		UserId:   editInput.UserId,
		Name:     editInput.Name,
		UserName: editInput.UserName,
		Bio:      editInput.Bio,
		Links:    editInput.Links,
	})

	if err != nil {
		fmt.Println("----------auth service down--------")

		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "can't edit user details",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "can't edit user details",
				Data:       resp,
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_auth.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "edited user profile successfully",
			Error:      nil,
		})

}

func (svc *UserHandler) GetFollowersDetails(ctx *fiber.Ctx) error {

	userId := ctx.Locals("userId")

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.GetFollowersDetails(context, &pb.RequestUserId{UserId: fmt.Sprint(userId)})
	if err != nil {
		fmt.Println("----------auth service down--------")

		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "can't fetch followers details",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "can't fetch followers details",
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_auth.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "fetched followers details successfully",
			Data:       resp,
			Error:      nil,
		})
}

func (svc *UserHandler) GetFollowingsDetails(ctx *fiber.Ctx) error {

	userId := ctx.Locals("userId")

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.GetFollowingsDetails(context, &pb.RequestUserId{UserId: fmt.Sprint(userId)})
	if err != nil {
		fmt.Println("----------auth service down--------")

		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "can't fetch followings details",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "can't fetch followings details",
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_auth.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "fetched followings details successfully",
			Data:       resp,
			Error:      nil,
		})

}

func (svc *UserHandler) GetAnotherUserProfile(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId")
	UserId := fmt.Sprint(userId)

	userBId := ctx.Params("userbid")

	if fmt.Sprint(userId) == "" || userBId == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "can't fetch user profile",
				Data:       nil,
				Error:      "no userbid found in request",
			})
	}

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.GetUserProfile(context, &pb.RequestGetUserProfile{
		UserId:  UserId,
		UserBId: userBId,
	})

	if err != nil {
		fmt.Println("----------auth service down--------")

		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "failed to get user profile",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "failed to get user profile",
				Data:       resp,
				Error:      resp.ErrorMessage,
			})
	}

	var respStruct responsemodel_auth.UserProfileB //used to show the zero count of posts,following,followers etc

	intValueuserBId, _ := strconv.Atoi(userBId)
	uintValueuserBId := uint(intValueuserBId)

	respStruct.UserId = uintValueuserBId
	respStruct.Name = resp.Name
	respStruct.UserName = resp.UserName
	respStruct.Bio = resp.Bio
	respStruct.Links = resp.Links
	respStruct.UserProfileImgURL = resp.ProfileImageURL
	respStruct.PostsCount = uint(resp.PostsCount)
	respStruct.FollowersCount = uint(resp.FollowerCount)
	respStruct.FollowingCount = uint(resp.FollowingCount)
	respStruct.FollowingStatus = resp.FollowingStat

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_auth.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "fetched user profile successfully",
			Data:       respStruct,
			Error:      nil,
		})
}

func (svc *UserHandler) SearchUser(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId")

	searchText := ctx.Params("searchtext")
	limit, offset := ctx.Query("limit", "5"), ctx.Query("offset", "0")

	fmt.Println("Search ", searchText)
	if searchText == "" {
		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "failed to get search result",
				Error:      "enter a valid name or username",
			})
	}

	validSearch := regexp.MustCompile(`^[a-zA-Z0-9_ ]+$`).MatchString
	if len(searchText) > 12 || !validSearch(searchText) {
		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "failed to get search result",
				Error:      "searchtext should contain only less than 12 letters and search input can only contain letters, numbers, spaces, or underscores",
			})
	}

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.SearchUser(context, &pb.RequestUserSearch{
		UserId:     fmt.Sprint(userId),
		SearchText: searchText,
		Limit:      limit,
		Offset:     offset,
	})

	if err != nil {
		fmt.Println("----------auth service down--------")

		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "failed to get search result",
				Error:      err.Error(),
			})
	}

	if resp.Errormessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "failed to get search result",
				Error:      resp.Errormessage,
			})
	}
	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_auth.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "fetched search result successfully",
			Data:       resp,
			Error:      nil,
		})
}

func (svc *UserHandler) SetProfileImage(ctx *fiber.Ctx) error {
	userId := ctx.Locals("userId")

	//fiber's ctx.BodyParser can't parse files(*multipart.FileHeader),
	//so we have to manually access the Multipart form and read the files form it.
	form, err := ctx.MultipartForm()
	if err != nil {
		return err
	}
	img := form.File["ProfileImg"]
	if len(img) == 0 || len(img) > 1 {

		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "can't add profile image",
				Error:      "no ProfileImg found in request,you should exactly upload only one img",
			})

	}
	ProfileImg := img[0]

	if ProfileImg.Size > 2*1024*1024 { // 2 MB limit
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "can't add profile image",
				Error:      "ProfileImg size exceeds the limit (2MB)",
			})
	}

	file, err := ProfileImg.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	contentType, err := mediafileformatchecker.ProfileImageFileFormatChecker(file)
	if err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "can't add profile image",
				Error:      err.Error(),
			})
	}

	content, err := byteconverter.MultipartFileheaderToBytes(&file)
	if err != nil {
		fmt.Println("-------------byteconverter-down---------")
	}

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.SetUserProfileImage(context, &pb.RequestSetProfileImg{
		UserId:      fmt.Sprint(userId),
		ContentType: *contentType,
		Img:         content,
	})

	if err != nil {
		fmt.Println("----------auth service down--------")

		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "can't add profile image",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "can't add profile image",
				Error:      resp.ErrorMessage,
			})
	}
	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_auth.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "profile image set successfully",
			Error:      nil,
		})

}

func (svc *UserHandler) AdminLogin(ctx *fiber.Ctx) error {
	var adminLoginData requestmodel_auth.AdminLoginData

	if err := ctx.BodyParser(&adminLoginData); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "Admin login failed (no JSON input)",
				Error:      err.Error(),
			})
	}

	if adminLoginData.Email == "" || adminLoginData.Password == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "Admin login failed (email and password required)",
				Error:      "Email and Password are required",
			})
	}

	// Send the request to the gRPC service
	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.AdminLogin(context, &pb.AdminLoginRequest{
		Email:    adminLoginData.Email,
		Password: adminLoginData.Password,
	})

	if err != nil {
		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "Admin login failed (service unavailable)",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "Admin login failed",
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_auth.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "Admin login successful",
			Data:       resp,
		})
}

func (svc *UserHandler) GetAllUsers(ctx *fiber.Ctx) error {
	limit := ctx.Query("limit")
	offset := ctx.Query("offset")

	if limit == "" || offset == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "Get all users failed (limit and offset required)",
				Error:      "Limit and Offset are required",
			})
	}

	// Send the request to the gRPC service
	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.GetAllUsers(context, &pb.GetAllUsersRequest{
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "Get all users failed (service unavailable)",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "Get all users failed",
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_auth.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "Get all users successful",
			Data:       resp,
		})
}

func (svc *UserHandler) BlockUser(ctx *fiber.Ctx) error {
	var blockUserReq struct {
		UserId string `json:"userId" validate:"required"`
	}

	if err := ctx.BodyParser(&blockUserReq); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "Block user failed (no JSON input)",
				Error:      err.Error(),
			})
	}

	if blockUserReq.UserId == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "Block user failed (user ID required)",
				Error:      "UserId is required",
			})
	}

	// Send the request to the gRPC service
	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.BlockUser(context, &pb.BlockUserRequest{
		UserId: blockUserReq.UserId,
	})

	if err != nil {
		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "Block user failed (service unavailable)",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "Block user failed",
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_auth.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "User blocked successfully",
		})
}

func (svc *UserHandler) UnblockUser(ctx *fiber.Ctx) error {
	var unblockUserReq struct {
		UserId string `json:"userId" validate:"required"`
	}

	if err := ctx.BodyParser(&unblockUserReq); err != nil {
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "Unblock user failed (no JSON input)",
				Error:      err.Error(),
			})
	}

	if unblockUserReq.UserId == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "Unblock user failed (user ID required)",
				Error:      "UserId is required",
			})
	}

	// Send the request to the gRPC service
	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.UnblockUser(context, &pb.UnblockUserRequest{
		UserId: unblockUserReq.UserId,
	})

	if err != nil {
		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "Unblock user failed (service unavailable)",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "Unblock user failed",
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_auth.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "User unblocked successfully",
		})
}

func (svc *UserHandler) GetLogFile(ctx *fiber.Ctx) error {
	// Read the log file
	logData, err := ioutil.ReadFile("app.log")
	if err != nil {
		// Respond with an internal server error if the file can't be read
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusInternalServerError,
				Message:    "Failed to read log file",
				Error:      err.Error(),
			})
	}

	// Set the Content-Type header to plain text and send the log data
	ctx.Set("Content-Type", "text/plain; charset=utf-8")
	return ctx.Status(fiber.StatusOK).Send(logData)
}

// CreateBlueTickPaymentHandler handles the creation of a blue tick payment order.
func (h *UserHandler) CreateBlueTickPaymentHandler(ctx *fiber.Ctx) error {
	var req struct {
		UserId uint `json:"userId" validate:"required"`
	}

	// Parse and validate request body
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "Invalid request data",
				Error:      err.Error(),
			})
	}

	// Call the gRPC service method
	resp, err := h.Client.CreateBlueTickPayment(context.Background(), &pb.CreateBlueTickPaymentRequest{
		UserId: uint32(req.UserId),
	})
	if err != nil {
		log.Printf("Error creating blue tick payment: %v", err)
		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "Service unavailable",
				Error:      err.Error(),
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_auth.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    resp.Message,
			Data: map[string]string{
				"verificationId": resp.VerificationId,
			},
		})
}

// VerifyBlueTickPaymentHandler verifies the blue tick payment.
func (h *UserHandler) VerifyBlueTickPaymentHandler(ctx *fiber.Ctx) error {
	var req struct {
		VerificationId string `json:"verificationId" validate:"required"`
		PaymentId      string `json:"paymentId" validate:"required"`
		Signature      string `json:"signature" validate:"required"`
		UserId         uint   `json:"userId" validate:"required"`
	}

	// Parse and validate request body
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "Invalid request data",
				Error:      err.Error(),
			})
	}

	// Call the gRPC service method
	resp, err := h.Client.VerifyBlueTickPayment(context.Background(), &pb.VerifyBlueTickPaymentRequest{
		VerificationId: req.VerificationId,
		PaymentId:      req.PaymentId,
		Signature:      req.Signature,
		UserId:         uint32(req.UserId),
	})
	if err != nil {
		log.Printf("Error verifying blue tick payment: %v", err)
		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "Service unavailable",
				Error:      err.Error(),
			})
	}

	// Check for an error message in the gRPC response
	if !resp.Success {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "Payment verification failed",
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_auth.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "Payment verified successfully",
		})
}

func (h *UserHandler) OnlinePayment(ctx *fiber.Ctx) error {
	// Extract query parameters for userID and verificationID
	userID := ctx.Params("userID")
	verificationID := ctx.Params("verificationID")

	// Log userID and verificationID for debugging
	fmt.Printf("Received OnlinePayment request for UserID: %s, VerificationID: %s\n", userID, verificationID)

	// Call the gRPC service to get blue tick verification details
	resp, err := h.Client.OnlinePayment(context.Background(), &pb.OnlinePaymentRequest{
		UserId:         userID,
		VerificationId: verificationID,
	})

	// Error handling and response based on gRPC service results
	if err != nil || resp == nil || resp.Paymentstatus == "" {
		// Render template with an error message
		if err := ctx.Status(fiber.StatusBadRequest).
			Render("razopay", fiber.Map{
				"badRequest": "Verification ID not found or invalid request data",
			}); err != nil {
			return fmt.Errorf("failed to render template: %w", err)
		}
		return nil
	}

	// Render template with verification details
	if err := ctx.Status(fiber.StatusOK).
		Render("razopay", fiber.Map{
			"user_id":          resp.UserId,
			"payment_status":   resp.Paymentstatus,
			"verification_fee": resp.Verificationfee,
			"order_id":         verificationID,
		}); err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}
	return nil
}
func (svc *UserHandler) GetAllVerifiedUsers(ctx *fiber.Ctx) error {
	limit := "10"
	offset := "0"

	if limit == "" || offset == "" {
		return ctx.Status(fiber.ErrBadRequest.Code).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.ErrBadRequest.Code,
				Message:    "Get all users failed (limit and offset required)",
				Error:      "Limit and Offset are required",
			})
	}

	// Send the request to the gRPC service
	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := svc.Client.GetAllVerifiedUsers(context, &pb.GetAllVerifiedUsersRequest{
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		return ctx.Status(fiber.StatusServiceUnavailable).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "Get all users failed (service unavailable)",
				Error:      err.Error(),
			})
	}

	if resp.ErrorMessage != "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(responsemodel_auth.CommonResponse{
				StatusCode: fiber.StatusBadRequest,
				Message:    "Get all verified users failed",
				Error:      resp.ErrorMessage,
			})
	}

	return ctx.Status(fiber.StatusOK).
		JSON(responsemodel_auth.CommonResponse{
			StatusCode: fiber.StatusOK,
			Message:    "Get all verified users successful",
			Data:       resp,
		})
}