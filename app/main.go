package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"io"

	"github.com/google/uuid"
)

func main() {

	deviceId := "01-202304012328-87495896"

	// turn on switchbot color bulb
	apiHeader := makeAPIHeader()
	SwitchBotControl(apiHeader, deviceId, "turnOn", "default", "command")

	// turn off switchbot color bulb
	apiHeader = makeAPIHeader()
	SwitchBotControl(apiHeader, deviceId, "turnOff", "default", "command")
}

func makeAPIHeader() map[string]string {
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

	return apiHeader
}

func SwitchBotControl(apiHeader map[string]string, deviceId string, command string, parameter string, commandType string) {
	// Build API Request Body
	apiRequestBody := map[string]string{
		"command":     command,
		"parameter":   parameter,
		"commandType": commandType,
	}

	// convert request body to json
	requestBody, err := json.Marshal(apiRequestBody)
	if err != nil {
		fmt.Println(err)
	}

	// request switchbot api with api header
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.switch-bot.com/v1.1/devices/%s/commands", deviceId), bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println(err)
	}

	// add headers to request
	for key, value := range apiHeader {
		req.Header.Add(key, value)
	}

	// add request body to request
	client := &http.Client{}

	// send request

	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()

	// read all response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	// print response body
	fmt.Println(string(body))
}
