package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var _ = fmt.Print

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
		default:
			fmt.Println(input + ": command not found")
		}
	}
}
