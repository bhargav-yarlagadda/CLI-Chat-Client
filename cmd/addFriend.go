package commands

import (
	"cli-chat-client/api"
	"cli-chat-client/data"
	"fmt"
	"strings"
)

func AddFriend(arguments []string){
	if len(arguments) == 0 {
		fmt.Println("Use --help for login usage")
		return 
	}
	for _, arg := range arguments {
		if strings.EqualFold(arg, "--help") {
			fmt.Println("Usage: send --to:`username`")
			return 
		}
	}
	if data.USERNAME == "" && data.PUBLIC_KEY == "" {
		fmt.Println("Please Login in , to send request")
		return 
		
	}
	if len(arguments) < 2 {
		fmt.Println("Not enough arguments. Use --help for usage.")
		return 
	}
	oppositePartyToken := arguments[1]
	if(!strings.HasPrefix(oppositePartyToken,"--to:")){
		fmt.Println("Invalid Usage of the command refer `--help`")
		return 
	}
	oppositeParty :=strings.TrimPrefix(oppositePartyToken,"--to:")
	api.SendRequestNotification(oppositeParty)
}