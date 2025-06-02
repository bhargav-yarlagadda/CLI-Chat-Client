package main

import (
	commands "cli-chat-client/cmd"
	"fmt"
	"strings"

	"github.com/chzyer/readline"
)
func main(){
	
	rl,err:=readline.New("> ")
	if err!=nil{
		panic(err)
	}
	defer rl.Close() 
	fmt.Println("Welcome to Lockline CLI!")
	for {
		line,err :=rl.Readline()
		if err!=nil{
			break
		}
		line = strings.TrimSpace(line)
		if line == ""{
			continue
		}
		if line == "exit" || line == "quit"{
			fmt.Println("Exiting out of lockline,clearing session data.")
			break
		}
		args := strings.Fields(line) // splitting the arguments into strings 
		command := strings.ToLower(args[0])
		switch command{
			case "login":
				commands.Login(args[1:])
			
			default:
				fmt.Println("Command unknown: ",command)
		}
	}

}