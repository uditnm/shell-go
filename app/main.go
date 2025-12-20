package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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
			fmt.Println(input + ": command not found")
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
	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, string(os.PathListSeparator))

	for _, dir := range paths {
		fullPath := filepath.Join(dir, input)

		_, err := os.Stat(fullPath)
		if err == nil && isExecutable(fullPath) {
			return input + " is " + fullPath
		}
	}

	return input + ": not found"
}

func isExecutable(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".exe", ".bat", ".cmd", ".com":
		return true
	default:
		return false
	}
}
