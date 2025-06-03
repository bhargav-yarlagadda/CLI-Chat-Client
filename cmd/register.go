package commands

import (
	"cli-chat-client/api"
	"fmt"
	"strings"
)

func Register(arguments []string) {
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

	if len(arguments) < 2 {
		fmt.Println("Not enough arguments. Use --help for usage.")
		return
	}

	usernameToken := arguments[0]
	passwordToken := arguments[1]

	if !strings.HasPrefix(usernameToken, "--username:") || !strings.HasPrefix(passwordToken, "--password:") {
		fmt.Println("Invalid format. Use: register --username:your_username --password:your_password")
		return
	}

	username := strings.TrimPrefix(usernameToken, "--username:")
	password := strings.TrimPrefix(passwordToken, "--password:") 
	if(username == "" || password == ""){
		fmt.Println("Invalid format-username and password cannot be empty. Use: register --username:your_username --password:your_password")
		return
	}
	publicKey,privateKey,message,success := api.RegisterUser(username,password) 
	if(success){
		fmt.Println(message)
		fmt.Println("Here are your public and private keys: ")
		fmt.Println("Public Key: ",publicKey)
		fmt.Println("Private Key: ",privateKey)
		fmt.Println("\nPlease securely store your private key this will not be generate again or not stored in our database,if this key is lost you cant perform further actions")
		fmt.Println("\nPlease login after storing the keys.")
	}else{
		fmt.Println("Errror in Registering user: ",message)
		return 
	}
}