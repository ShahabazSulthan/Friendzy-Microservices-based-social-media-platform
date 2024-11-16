package responsemodels

type SignUpResponse struct {
	Name        string
	UserName    string
	Email       string
	OTP         string
	Token       string
	IsUserExist string
}

type LoginResponse struct {
	AccessToken   string
	RefereshToken string
}

type OTPVerificationResponse struct {
	OTP           string
	AccessToken   string
	RefereshToken string
}

type UserProfileResponse struct {
	UserId            uint `gorm:"column:id"`
	Name              string
	UserName          string
	Bio               string
	Links             string
	UserProfileImgUrl string `gorm:"column:profile_img_url"`
	PostCount         uint
	FollowCount       uint
	FollowingCount    uint
	FollowedBy        string
	FollowingStatus   bool
	BlueTickVerified  string `json:"blue_tick_verified"`
}

type UserLiteResponse struct {
	UserName          string
	UserProfileImgUrl string `gorm:"column:profile_img_url"`
}

type UserListResponse struct {
	Id                uint32
	Name              string
	UserName          string
	UserProfileImgUrl string `gorm:"column:profile_img_url"`
}

type UserAdminResponse struct {
	ID              uint32 `json:"id"`
	Name            string `json:"name"`
	UserName        string `json:"user_name"`
	Email           string `json:"email"`
	Bio             string `json:"bio"`
	ProfileImageURL string `json:"profile_img_url"`
	Links           string `json:"links"`
	Status          string `json:"status"`
}

type AdminLoginres struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Result   string `json:"result,omitempty"`
	Token    string `json:"token,omitempty"`
}

type OnlinePayment struct {
	UserID          uint   `json:"user_id"`
	PaymentStatus   string `json:"payment_status"`   // e.g., success, pending, failed
	VerificationFee uint   `json:"verification_fee"` // Fixed fee for verification (e.g., 600)
}

type BlueTickResponse struct {
	ID               uint32 `json:"id"`
	Name             string `json:"name"`
	UserName         string `json:"user_name"`
	Email            string `json:"email"`
	Bio              string `json:"bio"`
	ProfileImageURL  string `json:"profile_img_url"`
	Links            string `json:"links"`
	Status           string `json:"status"`
	BlueTickVerified string `json:"blue_tick_verified"`
}
