package interface_regex

type IRegex interface {
	IsValidUsername(username string) (bool, string)
	IsValidPassword(password string) (bool, string)
}
