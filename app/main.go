package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

func main() {
	// Declare empty header map
	apiHeader := make(map[string]string)

	// Get token and secret from environment variables
	token := os.Getenv("SWITCHBOT_TOKEN")
	secret := os.Getenv("SWITCHBOT_CLIENT_SECRET")

	// Generate nonce and timestamp
	nonce := uuid.New()
	t := time.Now().UnixNano() / int64(time.Millisecond)

	// Create string to sign
	stringToSign := fmt.Sprintf("%s%d%s", token, t, nonce.String())

	// Create HMAC-SHA256 signature
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	sign := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Print values
	fmt.Printf("Authorization: %s\n", token)
	fmt.Printf("t: %d\n", t)
	fmt.Printf("sign: %s\n", sign)
	fmt.Printf("nonce: %s\n", nonce.String())

	// Build API header map
	apiHeader["Authorization"] = token
	apiHeader["Content-Type"] = "application/json"
	apiHeader["charset"] = "utf8"
	apiHeader["t"] = fmt.Sprintf("%d", t)
	apiHeader["sign"] = sign
	apiHeader["nonce"] = nonce.String()

	// Print API header
	for key, value := range apiHeader {
		fmt.Printf("%s: %s\n", key, value)
	}
}
