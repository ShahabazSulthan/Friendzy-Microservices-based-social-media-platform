package requestmodel_auth

type UserSignUpReq struct {
	Name            string `json:"Name" validate:"required,gte=3,lte=30"`
	UserName        string `json:"UserName" validate:"required,gte=3,lte=30"`
	Email           string `json:"Email" validate:"required,email"`
	Password        string `json:"Password" validate:"required,gte=3,lte=30"`
	ConfirmPassword string `json:"ConfirmPassword" validate:"required,eqfield=Password"`
}

type OtpVerification struct {
	Otp string `json:"Otp" validate:"required,len=4,number"`
}

type UserLoginReq struct {
	Email    string `json:"Email"    validate:"required,email"`
	Password string `json:"Password" validate:"required,min=4,max=30"`
}

type ForgotPasswordReq struct {
	Email string `json:"Email"    validate:"required,email"`
}

type ForgotPasswordData struct {
	Otp             string `json:"Otp" validate:"required,len=4,number"`
	Password        string `json:"Password" validate:"required,gte=3,lte=30"`
	ConfirmPassword string `json:"ConfirmPassword" validate:"required,eqfield=Password"`
}

type EditUserProfile struct {
	Name     string `json:"Name" validate:"required,gte=3,lte=25"`
	UserName string `json:"UserName" validate:"required,gte=3,lte=30"`
	Bio      string `json:"Bio" validate:"lte=50"`
	Links    string `json:"Links" validate:"lte=25"`
	UserId string
}

type AdminLoginData struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"min=4"`
}
