package hashedpassword

import (
	"errors"
	"fmt"

	interface_hash "github.com/ShahabazSulthan/Friendzy_Auth/pkg/utils/hashed_password/interfaces"
	"golang.org/x/crypto/bcrypt"
)

type HashUtil struct{}

func NewHashUtil() interface_hash.IHashPassword {
	return &HashUtil{}
}

func (h *HashUtil) HashedPassword(password string) string {
	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error in Hashing Password: ", err)
	}

	fmt.Println(HashedPassword)
	return string(HashedPassword)
}

func (h *HashUtil) ComparePassword(hashedPassword string, plainPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))

	if err != nil {
		return errors.New("password does not match")
	}

	return nil
}
