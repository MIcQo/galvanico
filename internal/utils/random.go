package utils

import (
	"cmp"
	"crypto/rand"
	"math/big"

	"golang.org/x/exp/constraints"
)

var characterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// RandomString generates a random string of n length
func RandomString(n int) string {
	b := make([]rune, n)

	for i := range b {
		var r, err = rand.Int(rand.Reader, big.NewInt(int64(len(characterRunes))))
		if err != nil {
			panic(err)
		}
		b[i] = characterRunes[r.Int64()]
	}
	return string(b)
}

// RandomIndex returns value of random index
func RandomIndex[E cmp.Ordered](slice []E) (val E, err error) {
	randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(slice))))
	if err != nil {
		return
	}
	return slice[randIndex.Int64()], nil
}

func RandomNumber[T constraints.Signed](minN, maxN T) (int64, error) {
	bg := big.NewInt(int64(maxN - minN))
	n, err := rand.Int(rand.Reader, bg)
	if err != nil {
		return 0, err
	}

	return n.Int64() + int64(minN), nil
}
