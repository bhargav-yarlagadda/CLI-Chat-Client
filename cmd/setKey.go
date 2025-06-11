package commands

import (
	"bufio"
	"cli-chat-client/data"
	"fmt"
	"os"
	"strings"
)


func SetKey(arguments []string) bool {
	if len(arguments) > 1 {
		for _, arg := range arguments {
			if strings.EqualFold(arg, "--help") {
				fmt.Println("Usage: set --key:your_private_key")
				return false
			}
			if strings.HasPrefix(arg, "--key:") {
				parts := strings.SplitN(arg, ":", 2)
				if len(parts) == 2 && parts[1] != "" {
					data.PRIVATE_KEY = parts[1]
					fmt.Println("Private key set successfully.")
					return true
				}
			}
		}
		fmt.Println("Invalid format. Use --key:your_private_key")
		return false
	}

	// Fallback: ask user to input manually
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter your private key. We do not store your keys; they are cleared when the session ends.")

	for {
		fmt.Print("Enter Private Key: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			fmt.Println("Private Key cannot be empty. Please try again.")
		} else {
			data.PRIVATE_KEY = input
			fmt.Println("Private key set successfully.")
			return true
		}
	}
}
