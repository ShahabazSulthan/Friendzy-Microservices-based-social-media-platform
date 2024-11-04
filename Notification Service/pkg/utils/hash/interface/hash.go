package interface_hash

type Ihash interface {
	HashPassword(password string) string
	CompairPassword(hashedPassword string, plainPassword string) error
}
