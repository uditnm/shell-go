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
		if input == "exit" {
			return
		}

		fmt.Println(input + ": command not found")
	}
}
