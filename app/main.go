package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

var _ = fmt.Print
var commands = []string{"exit", "echo", "type"}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")
		input, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("Error reading command:", err)
			continue
		}

		input = strings.TrimSpace(input)

		parts := strings.Fields(input)

		switch parts[0] {
		case "exit":
			return
		case "echo":
			// Implement echo command
			fmt.Println(strings.Join(parts[1:], " "))
		case "type":
			output := checkCommand(parts[1])
			fmt.Println(output)
		default:
			cmd := exec.Command(parts[0], parts[1:]...)
			err := cmd.Run()
			if err != nil {
				fmt.Println(input + ": command not found")
			}
		}
	}
}

func checkCommand(input string) string {
	if slices.Contains(commands, input) {
		return input + " is a shell builtin"
	}
	return checkExecutable(input)
}

func checkExecutable(input string) string {
	path, err := exec.LookPath(input)
	if err == nil {
		return input + " is " + path
	}

	return input + ": not found"
}
