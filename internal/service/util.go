package service

import (
	"crypto/sha256"
	"fmt"
)

func GetHash(plainText string) string {
	h := sha256.New()
	h.Write([]byte(plainText))
	return fmt.Sprintf("%x", h.Sum(nil))
}
