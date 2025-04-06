package user

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var firstWords = []string{
	"Spinning",
	"Cotton",
	"Power",
	"Steam",
	"Telegraph",
	"Telephone",
	"Locomotive",
	"Phonograph",
	"Sewing",
	"Electric",
	"X-ray",
	"Airplane",
	"Mechanical",
}

var secondWords = []string{
	"Jenny",
	"Loom",
	"Engine",
	"Gin",
	"Typewriter",
	"machine",
	"Bessemer",
	"light", "bulb",
	"combustion",
	"engine",
	"Bicycle",
	"Dynamite",
	"Elevator",
	"Motion",
	"picture",
	"camera",
}

const (
	minNumber = 100
	maxNumber = 999
)

func UsernameGenerator() (string, error) {
	var b strings.Builder

	for i := range 2 {
		result, err := randomWord(i)
		if err != nil {
			return "", err
		}

		b.WriteString(result)
		b.WriteString("_")
	}

	var randomNumber, err = randomInt(minNumber, maxNumber)
	if err != nil {
		return "", err
	}

	var caser = cases.Title(language.English)
	var wrdBld strings.Builder
	var wrds = strings.Split(b.String(), "_")
	wrdBld.WriteString(caser.String(wrds[0]))
	for _, word := range wrds[1:] {
		wrdBld.WriteString(caser.String(word))
	}

	wrdBld.WriteString(strconv.FormatInt(randomNumber, 10))

	return wrdBld.String(), nil
}

func randomWord(i int) (string, error) {
	if i == 0 {
		return randomFirstWord()
	}

	return randomSecondWord()
}

func randomInt(minN, maxN int64) (int64, error) {
	bg := big.NewInt(maxN - minN)
	n, err := rand.Int(rand.Reader, bg)
	if err != nil {
		return 0, err
	}

	return n.Int64() + minN, nil
}

func randomFirstWord() (string, error) {
	var firstWord, err = rand.Int(rand.Reader, big.NewInt(int64(len(firstWords))))
	if err != nil {
		return "", err
	}
	return firstWords[firstWord.Int64()], nil
}
func randomSecondWord() (string, error) {
	var secondWord, err = rand.Int(rand.Reader, big.NewInt(int64(len(secondWords))))
	if err != nil {
		return "", err
	}
	return secondWords[secondWord.Int64()], nil
}
