package usecase

import (
	"errors"
	"fmt"
	"log"

	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/config"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/repository/interfaces"
	interface_usecase "github.com/ShahabazSulthan/Friendzy_Auth/pkg/usecase/interfaces"
	interface_jwt "github.com/ShahabazSulthan/Friendzy_Auth/pkg/utils/JWT/Interface"
)

type JwtUsecase struct {
	JwtKeys  *config.Token
	JwtUtil  interface_jwt.Ijwt
	UserRepo interfaces.IUserRepository
}

func NewJwtUseCase(jwtKeys *config.Token, jwtUtil interface_jwt.Ijwt, userRepo interfaces.IUserRepository) interface_usecase.IJwt {
	return &JwtUsecase{
		JwtKeys:  jwtKeys,
		JwtUtil:  jwtUtil,
		UserRepo: userRepo,
	}
}

func (j *JwtUsecase) VerifyAccessToken(token *string) (*string, error) {

	userId, err := j.JwtUtil.VerifyAccessToken(*token, j.JwtKeys.UserSecurityKey)
	if err != nil {
		if userId == " " {
			fmt.Println("UserId Is Empty")
			return nil, err
		}
		fmt.Println("Error in verify accesToken")
		return nil, err
	}

	if ok, _ := j.UserRepo.IsUserBlocked(userId); ok {
		return nil, errors.New("user blocked by admin")
	}

	return &userId, nil
}

func (j *JwtUsecase) AccessRegenerator(accessToken string, refreshToken string) (string, error) {
	// Ensure JwtUtil and JwtKeys are initialized
	if j.JwtUtil == nil || j.JwtKeys == nil {
		log.Println("JwtUtil or JwtKeys not initialized")
		return "", errors.New("JwtUtil or JwtKeys not initialized")
	}

	// Log the tokens for debugging
	log.Printf("Access Token: %s, Refresh Token: %s", accessToken, refreshToken)

	// Verify access token
	userId, err := j.JwtUtil.VerifyAccessToken(accessToken, j.JwtKeys.UserSecurityKey)
	if err != nil {
		log.Printf("Error verifying access token: %v", err)
		return "", err
	}

	// Check if userId is valid
	if userId == "" {
		log.Println("Invalid userId")
		return "", errors.New("invalid userId")
	}

	// Verify refresh token
	err = j.JwtUtil.VerifyRefreshToken(refreshToken, j.JwtKeys.UserSecurityKey)
	if err != nil {
		log.Printf("Error verifying refresh token: %v", err)
		return "", err
	}

	// Check user status
	status, err := j.UserRepo.GetUserStateForAccessTokenGeneration(&userId)
	if err != nil {
		log.Printf("Error fetching user status: %v", err)
		return "", err
	}

	// Check if status is nil before dereferencing
	if status == nil {
		log.Println("User status is nil")
		return "", errors.New("user status is nil")
	}

	if *status == "blocked" {
		log.Println("User is blocked by admin")
		return "", errors.New("user id blocked by admin")
	}

	// Generate a new access token
	newAccessToken, err := j.JwtUtil.GenerateAccessToken(j.JwtKeys.UserSecurityKey, userId)
	if err != nil {
		log.Printf("Error generating new access token: %v", err)
		return "", err
	}

	log.Printf("New Access Token generated: %s", newAccessToken)

	return newAccessToken, nil
}

func (j *JwtUsecase) VerifyAdminToken(adminToken *string) (*string, error) {
	// Verify the admin token using the security key
	adminEmail, err := j.JwtUtil.VerifyAdminToken(*adminToken, j.JwtKeys.AdminSecurityKey)
	if err != nil {
		log.Printf("Error verifying admin token: %v", err)
		return nil, err
	}

	// Ensure adminEmail is valid
	if adminEmail == "" {
		log.Println("Invalid admin email")
		return nil, errors.New("invalid admin email")
	}

	log.Printf("Admin token verified successfully: %s", adminEmail)

	return &adminEmail, nil
}
