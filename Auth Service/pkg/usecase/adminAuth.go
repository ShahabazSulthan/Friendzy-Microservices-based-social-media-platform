package usecase

import (
	"errors"
	"strconv"

	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/config"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/models/requestmodels"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/models/responsemodels"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/repository/interfaces"
	interface_usecase "github.com/ShahabazSulthan/Friendzy_Auth/pkg/usecase/interfaces"
	interface_jwt "github.com/ShahabazSulthan/Friendzy_Auth/pkg/utils/JWT/Interface"
	interface_hash "github.com/ShahabazSulthan/Friendzy_Auth/pkg/utils/hashed_password/interfaces"
	interface_regex "github.com/ShahabazSulthan/Friendzy_Auth/pkg/utils/regex/interface"
)

type AdminUseCase struct {
	repo             interfaces.IAdminRepo
	JwtUtil          interface_jwt.Ijwt
	RegexUtil        interface_regex.IRegex
	TokenSecurityKey *config.Token
	HashUtil         interface_hash.IHashPassword
}

func NewAdminUseCase(adminRepo interfaces.IAdminRepo,
	jwtUtil interface_jwt.Ijwt,
	regexUtil interface_regex.IRegex,
	config *config.Token,
	hashUtil interface_hash.IHashPassword) interface_usecase.IAdminUsecase {
	return &AdminUseCase{
		repo:             adminRepo,
		JwtUtil:          jwtUtil,
		RegexUtil:        regexUtil,
		TokenSecurityKey: config,
		HashUtil:         hashUtil,
	}
}

func (a *AdminUseCase) AdminLogin(adminData *requestmodels.AdminLoginData) (*responsemodels.AdminLoginres, error) {
	var adminLoginRes responsemodels.AdminLoginres


	// Fetch the hashed password from the repository
	hashedPassword, err := a.repo.GetPassword(adminData.Email) // Assuming you fetch the password using email
	if err != nil {
		return nil, err // Return nil instead of empty response on error
	}

	// Compare the provided password with the hashed password
	passwordErr := a.HashUtil.ComparePassword(hashedPassword, adminData.Password)
	if passwordErr != nil {
		return nil, errors.New("invalid credentials") // More generic error message for security
	}

	// Generate an admin token
	accessToken, err := a.JwtUtil.GenerateAdminToken(a.TokenSecurityKey.AdminSecurityKey, adminData.Email) // Ensure you pass the correct parameters
	if err != nil {
		return nil, err // Return nil instead of empty response on error
	}

	// Set the token in the response
	adminLoginRes.Token = accessToken
	return &adminLoginRes, nil
}

func (a AdminUseCase) GetAllUsers(limit string, offset string) (*[]responsemodels.UserAdminResponse, error) {
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
	users, err := a.repo.AllUsers(limitInt, offsetInt)
	if err != nil {
		return nil, err // Return any errors from the repository
	}

	return users, nil // Return the list of users
}

func (a AdminUseCase) BlcokUser(userId string) error {
	err := a.repo.BlockUser(userId)
	if err != nil {
		return err
	}
	return nil
}

func (a AdminUseCase) UnblockUser(userId string) error {
	err := a.repo.UnblockUser(userId)
	if err != nil {
		return err
	}
	return nil
}
