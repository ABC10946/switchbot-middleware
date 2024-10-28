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
	"github.com/labstack/echo"
	yaml "gopkg.in/yaml.v3"
)

type SwitchBotConfiguration struct {
	SwitchBotConfiguration []SwitchBotConfigurationItem `yaml:"switchbot-configuration"`
}

type SwitchBotConfigurationItem struct {
	Name      string   `yaml:"name"`
	Type      string   `yaml:"type"`
	Path      string   `yaml:"path"`
	DeviceIds []string `yaml:"deviceIds"`
}

func main() {
	// read switchbot configuration from yaml
	filePath := "/app/switchbot-configuration.yaml"

	switchbotConfiguration, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	t := SwitchBotConfiguration{}
	err = yaml.Unmarshal(switchbotConfiguration, &t)
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, world!")
	})

	var toggles map[string]bool

	for _, item := range t.SwitchBotConfiguration {
		fmt.Println(item.Name)
		fmt.Println(item.Type)
		fmt.Println(item.Path)
		fmt.Println(item.DeviceIds)
		fmt.Println("")
		if item.Type == "toggle" {
			toggles = make(map[string]bool)
			toggles[item.Name] = false

			e.GET(item.Path, func(c echo.Context) error {
				for _, deviceId := range item.DeviceIds {
					apiHeader := makeAPIHeader()
					if toggles[item.Name] {
						SwitchBotControl(apiHeader, deviceId, "turnOff", "default", "command")
						toggles[item.Name] = false
					} else {
						SwitchBotControl(apiHeader, deviceId, "turnOn", "default", "command")
						toggles[item.Name] = true
					}
				}
				fmt.Println(toggles)
				return c.String(http.StatusOK, fmt.Sprintf("Toggle %s devices", item.Type))
			})
		} else {
			e.GET(item.Path, func(c echo.Context) error {
				for _, deviceId := range item.DeviceIds {
					apiHeader := makeAPIHeader()
					SwitchBotControl(apiHeader, deviceId, item.Type, "default", "command")
				}
				return c.String(http.StatusOK, fmt.Sprintf("Turn %s devices", item.Type))
			})
		}
	}

	e.Logger.Fatal(e.Start(":8080"))
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
