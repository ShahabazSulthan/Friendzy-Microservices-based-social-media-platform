package interface_jwt

type Ijwt interface {
	TempTokenForOtpVerification(securityKey string, email string) (string, error)
	GenerateRefreshToken(secrutKey string) (string, error)
	GenerateAccessToken(securityKey string, id string) (string, error)
	UnbindEmailFromClaim(tokenstring string, tempVerification string) (string, error)
	VerifyRefreshToken(AccessToken string, secretKey string) error
	VerifyAccessToken(token string, secretKey string) (string, error)
	GenerateAdminToken(securityKey string, email string) (string, error)
	VerifyAdminToken(token string, secretKey string) (string, error)
}
