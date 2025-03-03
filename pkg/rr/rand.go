package rr

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandString(n int, charSet string) (string, error) {
	if charSet == "" {
		charSet = letters
	}
	var result strings.Builder
	charsetLength := big.NewInt(int64(len(charSet)))
	for i := 0; i < n; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", nil
		}
		result.WriteByte(charSet[randomIndex.Int64()])
	}
	return result.String(), nil
}

func GenString(n int) string {
	data, _ := RandString(n, "")
	return data
}
