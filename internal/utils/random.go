package utils

import (
	"crypto/rand"
	"math/big"
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
