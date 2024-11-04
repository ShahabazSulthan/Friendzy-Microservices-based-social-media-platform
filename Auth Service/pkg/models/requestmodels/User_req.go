package requestmodels

type UserSignUpRequest struct {
	Name            string
	UserName        string
	Email           string
	Password        string
	ConformPassword string
}

type UserLoginRequest struct {
	Email    string
	Password string
}

type ForgotPasswordRequest struct {
	OTP            string
	Password       string
	ConformPassword string
}

type EditUserProfileRequest struct {
	Name     string
	UserName string
	Bio      string
	Links    string
	UserId   string
}

type AdminLoginData struct {
	Email    string 
	Password string 
}
