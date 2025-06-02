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

func Login(username string, password string) (string, string, error) {
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
		Post(BACKEND_URL)

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
