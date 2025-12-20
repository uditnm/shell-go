package main

import (
	"fmt"
	"bufio"
	"os"
)

var _ = fmt.Print

func main() {
	fmt.Print("$ ")
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')

	if(err != nil){
		fmt.Println("Error reading command:", err)
		return
	}

	fmt.Print(input[:len(input) - 1] + ": command not found")
}
