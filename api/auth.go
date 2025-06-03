package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

func LoginUser(username string, password string) (string, string, error) {
	err:=godotenv.Load() 
	if err != nil {
        log.Fatalf("Error loading .env file")
    }
	BACKEND_URL := os.Getenv("BACKEND_URL")
	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"username": username,
			"password": password,
		}).
		Post(BACKEND_URL+"/login")

	if err != nil {
		return "", "", fmt.Errorf("request failed: %v", err)
	}

	// Structs to unmarshal response
	var successResp struct {
		Message   string `json:"message"`
		PublicKey string `json:"publicKey"`
		Token     string `json:"token"`
		Success   bool   `json:"success"`
	}

	var errorResp struct {
		Error   string `json:"error"`
		Success bool   `json:"success"`
	}

	// Handle response based on success/failure (200 to 299 success codes)
	if resp.IsSuccess() {
		if err := json.Unmarshal(resp.Body(), &successResp); err != nil {
			return "", "", fmt.Errorf("parsing success response failed: %v", err)
		}
		fmt.Println("✅", successResp.Message)
		return successResp.PublicKey, successResp.Token, nil
	} else {
		if err := json.Unmarshal(resp.Body(), &errorResp); err != nil {
			return "", "", fmt.Errorf("parsing error response failed: %v", err)
		}
		return "", "", errors.New(errorResp.Error)
	}
}
func RegisterUser(username string, password string) (string, string, string, bool) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	BACKEND_URL := os.Getenv("BACKEND_URL")
	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"username": username,
			"password": password,
		}).
		Post(BACKEND_URL+"/register")

	if err != nil {
		return "", "", "request failed: " + err.Error(), false
	}

	var successResp struct {
		Message    string `json:"message"`
		PublicKey  string `json:"publicKey"`
		PrivateKey string `json:"privateKey"`
		Success    bool   `json:"success"`
	}

	var errorResp struct {
		Error   string `json:"error"`
		Success bool   `json:"success"`
	}

	if resp.IsSuccess() {
		if err := json.Unmarshal(resp.Body(), &successResp); err != nil {
			return "", "", "parsing success response failed: " + err.Error(), false
		}
		fmt.Println("✅", successResp.Message)
		return successResp.PublicKey, successResp.PrivateKey, successResp.Message, true
	} else {
		if err := json.Unmarshal(resp.Body(), &errorResp); err != nil {
			return "", "", "parsing error response failed: " + err.Error(), false
		}
		return "", "", errorResp.Error, false
	}
}
