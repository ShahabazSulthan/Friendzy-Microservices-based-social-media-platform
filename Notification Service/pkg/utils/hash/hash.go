package hash

import (
	"errors"
	"fmt"

	interface_hash "github.com/ShahabazSulthan/Friendzy_Notification/pkg/utils/hash/interface"
	"golang.org/x/crypto/bcrypt"
)

type HashUtil struct{}

func NewHashUtil() interface_hash.Ihash {
	return &HashUtil{}
}

func (h *HashUtil) HashPassword(password string) string {

	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err, "problem at hashing ")
	}
	fmt.Println(HashedPassword)
	return string(HashedPassword)
}

func (h *HashUtil) CompairPassword(hashedPassword string, plainPassword string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))

	if err != nil {
		return errors.New("passwords does not match")
	}

	return nil
}
