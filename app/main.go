package main

import (
	"bufio"
	"fmt"
	"os"
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

		fmt.Println(input[:len(input)-1] + ": command not found")
	}
}
