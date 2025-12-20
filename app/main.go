package main

import (
	"fmt",
	"bufio",
	"os"
)

var _ = fmt.Print

func main() {
	fmt.Print("$ ")
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')

	fmt.Print(input[:len(input) - 1] + ": command not found")
}
