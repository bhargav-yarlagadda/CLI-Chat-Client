package api

import (
	"cli-chat-client/data"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

/* ─────────────────────────── CONSTANT ─────────────────────────── */

const backendBaseURL = "http://localhost:8000" // <<— single source of truth

/* ───────────────────────── SEND REQUEST ───────────────────────── */

func SendRequestNotification(toUser string) {
	if toUser == "" {
		fmt.Println("Receiver cannot be empty, please specify a username")
		return
	}
	if data.JWT_TOKEN == "" {
		fmt.Println("Please log in first")
		return
	}

	resp, err := resty.New().R().
		SetAuthToken(data.JWT_TOKEN).
		SetBody(map[string]interface{}{"to": toUser}).
		Post(backendBaseURL + "/friend-request/send")
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return
	}

	var ok struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	var bad struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
	}
	if resp.IsSuccess() {
		_ = json.Unmarshal(resp.Body(), &ok)
		fmt.Println("✅ request sent:", ok.Message)
	} else {
		_ = json.Unmarshal(resp.Body(), &bad)
		fmt.Println("Cannot send request:", bad.Error)
	}
}

/* ──────────────────────── NOTIFICATION TYPES ─────────────────── */

type FriendRequest struct {
	ID        string    `json:"id"`
	From      string    `json:"from"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type RequestsResponse struct {
	Success  bool            `json:"success"`
	Requests []FriendRequest `json:"requests"`
}

/* ─────────────────── GET ALL PENDING NOTIFICATIONS ───────────── */

func GetAllNotifications() ([]FriendRequest, error) {
	if data.JWT_TOKEN == "" {
		return nil, fmt.Errorf("please login to view friend requests")
	}

	url := backendBaseURL + "/friend-request?status=pending"

	var respData RequestsResponse
	resp, err := resty.New().R().
		SetAuthToken(data.JWT_TOKEN).
		SetResult(&respData).
		Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch friend requests: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("server returned %s", resp.Status())
	}

	return respData.Requests, nil
}
