package service

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func GetRandomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GetHash(plainText string) string {
	h := sha256.New()
	h.Write([]byte(plainText))
	return fmt.Sprintf("%x", h.Sum(nil))
}
