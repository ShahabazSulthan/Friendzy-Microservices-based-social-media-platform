package jwt

import (
	"errors"
	"fmt"
	"time"

	interface_jwt "github.com/ShahabazSulthan/Friendzy_Auth/pkg/utils/JWT/Interface"
	"github.com/golang-jwt/jwt"
)

type jwtUtil struct{}

func NewjwtUtil() interface_jwt.Ijwt {
	return &jwtUtil{}
}

func (j *jwtUtil) TempTokenForOtpVerification(securityKey string, email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenstring, err := token.SignedString([]byte(securityKey))
	if err != nil {
		fmt.Println("Error at jwt toke : ", err)
	}
	return tokenstring, err
}

func (j *jwtUtil) GenerateRefreshToken(secrutKey string) (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Unix() + 604800,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secrutKey))
	if err != nil {
		fmt.Println("error Occured While Creating Token :", err)
		return "", err
	}

	return signedToken, nil
}

func (j *jwtUtil) GenerateAccessToken(securityKey string, id string) (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Unix() + 36000,
		"id":  id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(securityKey))

	if err != nil {
		fmt.Println("error creating acesss token :", err)
		return "", err
	}

	return signedToken, nil
}

func (j *jwtUtil) UnbindEmailFromClaim(tokenstring string, tempVerification string) (string, error) {
	secret := []byte(tempVerification)
	parsedToken, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil || !parsedToken.Valid {
		fmt.Println(err)
		return "", err
	}

	claims := parsedToken.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	return email, nil
}

func (j *jwtUtil) VerifyRefreshToken(AccessToken string, secretKey string) error {
	key := []byte(secretKey)

	_, err := jwt.Parse(AccessToken, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		NewResp := err.Error() + ":RefreshToken"
		fmt.Println("NewResp = ", NewResp)
		return errors.New(NewResp)
	}

	return nil
}

func (j *jwtUtil) VerifyAccessToken(token string, secretKey string) (string, error) {
	key := []byte(secretKey)
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	// If parsing failed, check the specific error and handle accordingly
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				// Token is malformed
				fmt.Println("malformed token")
				return "", errors.New("malformed token")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				claims, ok := parsedToken.Claims.(jwt.MapClaims)
				if !ok {
					fmt.Println("failed to extract token claims")
					return "", errors.New("failed to extract claims")
				}

				id, ok := claims["id"].(string)
				if !ok {
					fmt.Println("id calim not found or not a string")
					return "", errors.New("ID claim not found or not a string")
				}

				// Token is expired or not valid yet
				fmt.Println("token expired")
				return id, errors.New("expired token")
			} else {
				// Other validation errors
				fmt.Println("validation error")
				return "", errors.New("validation error")
			}
		} else {
			// Other parsing errors
			fmt.Println("other error:", err)
			return "", err
		}
	}

	// If the token is valid, extract claims and return the ID
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("token valid,but failed to extract claims")
		return "", errors.New("failed to extract claims")
	}

	id, ok := claims["id"].(string)
	if !ok {
		fmt.Println("token valid,id claim not found or not a string")
		return "", errors.New("ID claim not found or not a string")
	}

	return id, nil
}

// GenerateAdminToken creates an access token for an admin
func (j *jwtUtil) GenerateAdminToken(securityKey string, email string) (string, error) {
	claims := jwt.MapClaims{
		"exp":   time.Now().Unix() + 36000, // 10 hours
		"email": email,                     // Include email as claim
		"role":  "admin",                   // Specify that this token is for an admin
	}

	fmt.Println("=====", securityKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(securityKey))
	if err != nil {
		fmt.Println("Error creating admin token: ", err)
		return "", err
	}

	return signedToken, nil
}

func (j *jwtUtil) VerifyAdminToken(tokenString, secretKey string) (string, error) {
	key := []byte(secretKey)

	fmt.Println(tokenString)

	fmt.Println("---", secretKey)
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	fmt.Println("parsedToken:", parsedToken)
	fmt.Println("parsedToken.Valid:", parsedToken.Valid)

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				fmt.Println("Malformed token")
				return "", errors.New("malformed token")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				fmt.Println("Expired or not valid yet token")
				return "", errors.New("expired or not valid yet token")
			} else {
				fmt.Println("Validation error")
				return "", errors.New("validation error")
			}
		} else {
			fmt.Println("Parsing error:", err)
			return "", err
		}
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		email, ok := claims["email"].(string)
		if !ok {
			fmt.Println("Email claim not found or not a string")
			return "", errors.New("email claim not found or not a string")
		}

		role, ok := claims["role"].(string)
		if !ok || role != "admin" {
			fmt.Println("Role claim not found or not an admin")
			return "", errors.New("role claim not found or not an admin")
		}

		exp, ok := claims["exp"].(float64)
		if !ok || time.Now().Unix() > int64(exp) {
			fmt.Println("Token expired")
			return "", errors.New("token expired")
		}

		fmt.Println("Token verified successfully")
		return email, nil
	}

	fmt.Println("Failed to extract claims")
	return "", errors.New("failed to extract claims")
}
