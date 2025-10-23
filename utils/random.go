package utils

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

// usage example token := utils.GenerateRandomString(16)
func GenerateRandomString(nSize int) string {
	byteData := make([]byte, nSize)
	_, err := rand.Read(byteData)
	if err != nil {
		log.Fatalf("Failed to generate random string: %v", err)
	}
	return base64.URLEncoding.EncodeToString(byteData)
}
