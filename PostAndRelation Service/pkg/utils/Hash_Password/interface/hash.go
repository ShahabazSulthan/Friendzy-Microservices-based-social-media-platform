package interface_hash

type IhashPassword interface {
	HashPassword(password string) string
	CompairPassword(hashedPassword string, plainPassword string) error
}
