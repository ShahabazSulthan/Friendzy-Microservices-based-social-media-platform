package interface_hash

type IHashPassword interface {
	HashedPassword(password string) string
	ComparePassword(hashedPassword string, plainPassword string) error
}
