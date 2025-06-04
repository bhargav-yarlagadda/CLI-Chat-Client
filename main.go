package main

import (
	commands "cli-chat-client/cmd"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/chzyer/readline"
)

func clearScreen() {
	// ANSI escape code to clear screen and move cursor to top-left
	fmt.Print("\033[H\033[2J")

	// Fallback: for older terminals or if ANSI is ignored
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func main() {

	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()
	fmt.Println("Welcome to Lockline CLI!")
	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "exit" || line == "quit" {
			fmt.Println("Exiting out of lockline, clearing session data.")
			os.Clearenv()
			rl.SetPrompt("") // Remove prompt
			rl.Clean()       // Clear readline buffer
			clearScreen()    // Now clear the terminal
			os.Exit(0)       // Exit completely
		}
		args := strings.Fields(line) // splitting the arguments into strings
		command := strings.ToLower(args[0])
		switch command {
		case "login":
			commands.Login(args[1:])
			
		case "register":
			commands.Register(args[1:])
		case "add":
			commands.AddFriend(args)
		default:
			fmt.Println("Command unknown: ", command)
		}
	}

}
