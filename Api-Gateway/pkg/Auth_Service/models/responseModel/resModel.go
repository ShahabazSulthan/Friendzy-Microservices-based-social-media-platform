package responsemodel_auth

type SignupData struct {
	Name            string `json:"name,omitempty"`
	UserName        string `json:"username,omitempty"`
	Email           string `json:"email,omitempty"`
	Password        string `json:"password,omitempty"`
	OTP             string `json:"otp,omitempty"`
	Token           string `json:"token,omitempty"`
	ConfirmPassword string `json:"confirmPassword,omitempty"`
	IsUserExist     string `json:"isUserExist,omitempty"`
}

type OtpVerifResult struct {
	Email        string `json:"email,omitempty"`
	Otp          string `json:"otp,omitempty"`
	Result       string `json:"result,omitempty"`
	Token        string `json:"token,omitempty"`
	AccessToken  string `json:"accesstoken,omitempty"`
	RefreshToken string `json:"refreshtoken,omitempty"`
}

type UserLoginRes struct {
	Email        string `json:"email,omitempty"`
	Password     string `json:"password,omitempty"`
	AccessToken  string `json:"accesstoken,omitempty"`
	RefreshToken string `json:"refreshtoken,omitempty"`
}

type ForgotPasswordRes struct {
	Email string `json:"email,omitempty"`
	Token string `json:"token,omitempty"`
}

type ForgotPasswordData struct {
	Token           string `json:"token,omitempty"`
	Otp             string `json:"otp,omitempty"`
	Password        string `json:"password,omitempty"`
	ConfirmPassword string `json:"confirmPassword,omitempty"`
}

type EditUserProfileResp struct {
	Name     string `json:"name,omitempty" `
	UserName string `json:"username,omitempty" `
	Bio      string `json:"bio,omitempty"`
	Links    string `json:"links,omitempty"`
}

type UserProfileA struct {
	UserId uint `json:"UserId"  gorm:"column:id"`

	Name              string `json:"Name"`
	UserName          string `json:"UserName"`
	Bio               string `json:"Bio"`
	Links             string `json:"Links"`
	UserProfileImgURL string `json:"UserProfileImageURL"`

	PostsCount     uint `json:"PostsCount"`
	FollowersCount uint `json:"FollowersCount"`
	FollowingCount uint `json:"FollowingCount"`
}
type UserProfileB struct {
	UserId uint `json:"UserId"  gorm:"column:id"`

	Name              string `json:"Name"`
	UserName          string `json:"UserName"`
	Bio               string `json:"Bio"`
	Links             string `json:"Links"`
	UserProfileImgURL string `json:"UserProfileImageURL"`

	PostsCount     uint `json:"PostsCount"`
	FollowersCount uint `json:"FollowersCount"`
	FollowingCount uint `json:"FollowingCount"`

	//for userB only
	FollowedBy      string `json:"Followedby,omitempty"`
	FollowingStatus bool   `json:"Following_status"`
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
