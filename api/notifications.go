package api

import (
	"cli-chat-client/data"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

func SendRequestNotification(ToUser string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
		return
	}
	BACKEND_URL := os.Getenv("BACKEND_URL")
	if data.JWT_TOKEN == "" || data.PUBLIC_KEY == "" || data.USERNAME == "" {
		fmt.Println("Please login in to send Requests")
		return
	}
	if ToUser == "" {
		fmt.Println("Reciever cannot be empty,please specify a username")
		return
	}
	client := resty.New()
	resp, err := client.R().SetHeader("Content-Type", "application/json").SetAuthToken(data.JWT_TOKEN).SetBody(map[string]interface{}{
		"to": ToUser,
	}).Post(BACKEND_URL + "/friend-request/send")
	if err != nil {
		fmt.Println("Failed To send Request: ", err)
	}
	var successResp struct {
		Status  bool   `json:"success"`
		Message string `json:"message"`
	}
	var errorResp struct {
		Status       bool   `json:"success"`
		ErrorMessage string `json:"error"`
	}
	if resp.IsSuccess() {
		if err := json.Unmarshal(resp.Body(), &successResp); err != nil {
			fmt.Printf("parsing success response failed: %v", err)
		}
		fmt.Println("✅ send request succesfull:", successResp.Message)
		return
	} else {
		if err := json.Unmarshal(resp.Body(), &errorResp); err != nil {
			fmt.Printf("parsing success response failed: %v", err)
			return 
		}
		fmt.Println("Cannot send Request: ",errorResp.ErrorMessage)
	}
}
