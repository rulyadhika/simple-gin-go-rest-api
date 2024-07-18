package helper

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	mathRand "math/rand"
)

func GenerateRandomHashString() string {
	// create 16 bytes of byte slice
	randomBytes := make([]byte, 16)

	// fill the slice with random data
	rand.Read(randomBytes)

	// compute md5 hash
	hash := md5.Sum(randomBytes)

	// convert hexadecimal to string
	hashString := hex.EncodeToString(hash[:])

	return hashString
}

func GenerateFixedLengthRandomNumber(length int) string {
	digits := "0123456789"
	randNumber := make([]byte, length)

	for i, _ := range randNumber {
		randNumber[i] = digits[mathRand.Intn(len(digits))]
	}

	return string(randNumber)
}
