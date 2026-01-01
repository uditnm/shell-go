package main

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

func getTokens(input string) ([]string, error) {
	var tokens []string
	var currentString strings.Builder

	singleQuote := false
	doubleQuote := false
	isBackSlash := false

	for i := 0; i < len(input); i++ {
		char := input[i]

		if isBackSlash {
			currentString.WriteByte(char)
			isBackSlash = !isBackSlash
			continue
		}

		switch char {
		case '\'':
			if !doubleQuote {
				singleQuote = !singleQuote
			} else {
				currentString.WriteByte(char)
			}
		case '"':
			if !singleQuote {
				doubleQuote = !doubleQuote
			} else {
				currentString.WriteByte(char)
			}
		case '\\':
			if singleQuote {
				currentString.WriteByte(char)
			} else if doubleQuote {
				if i < len(input)-1 && (input[i+1] == '"' || input[i+1] == '\\') {
					isBackSlash = !isBackSlash
				} else {
					currentString.WriteByte(char)
				}
			} else {
				isBackSlash = !isBackSlash
			}
		case ' ':
			if singleQuote || doubleQuote {
				currentString.WriteByte(char)
			} else if currentString.Len() > 0 {
				tokens = append(tokens, currentString.String())
				currentString.Reset()
			}
		default:
			currentString.WriteByte(char)
		}
	}

	if currentString.Len() > 0 {
		tokens = append(tokens, currentString.String())
	}

	if singleQuote || doubleQuote || isBackSlash {
		return nil, fmt.Errorf("invalid input")
	}

	return tokens, nil
}

func getPresentWorkingDirectory() (string, string) {
	path, err := os.Getwd()
	if err != nil {
		return "", err.Error()
	}

	return path, ""
}

func changeDirectory(path string) string {
	if path == HomeDirectory {
		home, _ := os.UserHomeDir()
		path = home
	}

	err := os.Chdir(path)
	if err != nil {
		return path + ": No such file or directory"
	}

	return ""
}

func checkCommand(input string) (string, string) {
	if slices.Contains(Commands(), input) {
		return input + " is a shell builtin", ""
	}

	path, err := exec.LookPath(input)
	if err == nil {
		return input + " is " + path, ""
	}

	return "", input + ": not found"
}
