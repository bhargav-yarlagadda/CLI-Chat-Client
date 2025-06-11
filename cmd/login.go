package commands

import (
	
	"cli-chat-client/api"
	"cli-chat-client/data"
	"fmt"
	
	"strings"
)

func Login(arguments []string) {
	if len(arguments) == 0 {
		fmt.Println("Use --help for login usage")
		return
	}

	for _, arg := range arguments {
		if strings.EqualFold(arg, "--help") {
			fmt.Println("Usage: login --username:your_username --password:your_password")
			return
		}
	}
	if data.USERNAME != "" || data.PUBLIC_KEY != "" {
		fmt.Println("You are already logged in please close current session to login again.")
		return
	}
	if len(arguments) < 2 {
		fmt.Println("Not enough arguments. Use --help for usage.")
		return
	}

	usernameToken := arguments[0]
	passwordToken := arguments[1]

	if !strings.HasPrefix(usernameToken, "--username:") || !strings.HasPrefix(passwordToken, "--password:") {
		fmt.Println("Invalid format. Use: login --username:your_username --password:your_password")
		return
	}

	username := strings.TrimPrefix(usernameToken, "--username:")
	password := strings.TrimPrefix(passwordToken, "--password:")

	publicKey, token, err := api.LoginUser(username, password)
	if err != nil {
		fmt.Println("❌", err)
		return
	}
	data.USERNAME = username
	data.PUBLIC_KEY = publicKey
	data.JWT_TOKEN = token
	fmt.Println("🔑 Public Key:", publicKey)
	fmt.Println("🔐 Token:", token)
	SetKey(make([]string,0))
}
