package interface_usecase

type IJwt interface {
	VerifyAccessToken(token *string) (*string, error)
	AccessRegenerator(accessToken string, refreshToken string) (string, error)
	VerifyAdminToken(adminToken *string) (*string, error)
}


