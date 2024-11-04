package regex

import (
	"regexp"
	"strings"

	interface_regex "github.com/ShahabazSulthan/Friendzy_Auth/pkg/utils/regex/interface"
)

type RegexUtil struct{}

func NewRegexUtil() interface_regex.IRegex {
	return &RegexUtil{}
}

func (rgx *RegexUtil) IsValidUsername(username string) (bool, string) {
	minLength := 3
	maxLength := 30
	validChars := `[a-zA-Z0-9._]+$` // Allow uppercase and lowercase letters, numbers, dots, and underscores

	if len(username) < minLength {
		return false, "Username must be at least 3 characters long"
	}

	if len(username) > maxLength {
		return false, "Username cannot exceed 30 characters"
	}

	if username == "" || !regexp.MustCompile(validChars).MatchString(username) {
		return false, "Username contains invalid characters. Only letters, numbers, dots, and underscores are allowed."
	}

	if username[0] == '.' || username[len(username)-1] == '.' {
		return false, "Username cannot start or end with a dot (.)"
	}

	if strings.Contains(username, "..") {
		return false, "Username cannot contain consecutive dots (..)"
	}

	return true, ""
}

func (rgx *RegexUtil) IsValidPassword(password string) (bool, string) {
	minLength := 8
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasDigit = true
		default:
			hasSpecial = true
		}
	}

	if len(password) < minLength {
		return false, "Password must be at least 8 characters long"
	}

	if !hasUpper {
		return false, "Password must contain at least one uppercase letter"
	}

	if !hasLower {
		return false, "Password must contain at least one lowercase letter"
	}

	if !hasDigit {
		return false, "Password must contain at least one digit"
	}

	if !hasSpecial {
		return false, "Password must contain at least one special character"
	}

	return true, ""
}
