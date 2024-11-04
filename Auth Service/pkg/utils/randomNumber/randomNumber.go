package randomnumber

import (
	"math/rand"

	interface_random "github.com/ShahabazSulthan/Friendzy_Auth/pkg/utils/randomNumber/interface"
)

type RandomNum struct{}

func NewRandomNumUtil() interface_random.IRandGene {
	return &RandomNum{}
}

func (rn RandomNum) RandomNumber() int {
	randomInt := rand.Intn(9000) + 1000
	return randomInt
}
