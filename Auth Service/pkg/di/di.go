package di

import (
	"fmt"

	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/client"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/config"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/db"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/repository"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/server"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/usecase"
	jwt "github.com/ShahabazSulthan/Friendzy_Auth/pkg/utils/JWT"
	hashedpassword "github.com/ShahabazSulthan/Friendzy_Auth/pkg/utils/hashed_password"
	randomnumber "github.com/ShahabazSulthan/Friendzy_Auth/pkg/utils/randomNumber"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/utils/regex"
	"github.com/ShahabazSulthan/Friendzy_Auth/pkg/utils/smtp"
)

func InitializeAuthService(config *config.Config) (*server.AuthService, error) {
	hashUtil := hashedpassword.NewHashUtil()
	DB, err := db.ConnectDatabase(&config.DB, hashUtil)
	if err != nil {
		fmt.Println("Error in connecting database from Dependency Injection")
		return nil, err
	}
	smtpUtil := smtp.NewSmtpCreadintial(&config.Smtp)
	// JWT Utility
	jwtutil := jwt.NewjwtUtil()
	if jwtutil == nil {
		fmt.Println("Error initializing JWT utility")
		return nil, fmt.Errorf("JWT utility is nil")
	}
	fmt.Println("JWT Util Initialized")

	randNumber := randomnumber.NewRandomNumUtil()
	regexUtil := regex.NewRegexUtil()
	postANDClient, err := client.InitPostClientService(config)
	if err != nil {
		return nil, err
	}
	UserRepo := repository.NewUserRepo(DB)
	UseruseCase := usecase.NewUserCase(UserRepo, smtpUtil, jwtutil, randNumber, regexUtil, &config.Token, hashUtil, postANDClient, &config.Razopay)
	// JWT Use Case
	jwtUsecase := usecase.NewJwtUseCase(&config.Token, jwtutil, UserRepo)
	if jwtUsecase == nil {
		fmt.Println("Error initializing JWT UseCase")
		return nil, fmt.Errorf("JWT UseCase is nil")
	}
	fmt.Println("JWT Use Case Initialized")

	adminRepo := repository.NewAdminRepo(DB)
	adminUsecase := usecase.NewAdminUseCase(adminRepo, jwtutil, regexUtil, &config.Token, hashUtil)

	paymentRepo := repository.NewPaymentRepo(DB)
	paymentUsecase := usecase.NewPaymenyUsecase(paymentRepo, &config.Razopay)

	EmbeddingStruct := server.NewAuthServer(UseruseCase, jwtUsecase, adminUsecase, paymentUsecase)
	return EmbeddingStruct, nil
}
