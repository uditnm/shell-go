package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

var _ = fmt.Print

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("$ ")

mainloop:
	for {

		input, err := reader.ReadString('\n')

		if err != nil {
			fmt.Print("Error reading command:", err)
			continue
		}

		input = strings.TrimSpace(input)

		parts, quoteErr := getTokens(input)

		if quoteErr != nil {
			fmt.Print("Error reading command:", quoteErr)
			continue
		}

		fileName := ""
		lastIdx := -1
		redirectStdError := false
		for i := 1; i < len(parts); i++ {
			if parts[i] == Redirect || parts[i] == StandardRedirect || parts[i] == ErrorRedirect {
				if i+1 >= len(parts) {
					fmt.Print("Error: No file given")
					continue mainloop
				}

				redirectStdError = parts[i] == ErrorRedirect

				lastIdx = i
				fileName = strings.Join(parts[(i+1):], " ")
				break
			}
		}

		if lastIdx != -1 {
			parts = parts[:lastIdx]
		}

		output := ""
		errOutput := ""

		switch parts[0] {
		case Exit:
			return
		case Echo:
			output = strings.Join(parts[1:], " ")
		case Type:
			output, errOutput = checkCommand(parts[1])
		case Pwd:
			output, errOutput = getPresentWorkingDirectory()
		case Cd:
			errOutput = changeDirectory(parts[1])
		default:
			_, err := exec.LookPath(parts[0])
			if err != nil {
				errOutput = parts[0] + ": command not found"
				break
			}

			var stderrBuf, stdoutBuf bytes.Buffer

			cmd := exec.Command(parts[0], parts[1:]...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = &stdoutBuf
			cmd.Stderr = &stderrBuf

			execErr := cmd.Run()

			output = stdoutBuf.String()
			errOutput = stderrBuf.String()

			if execErr != nil {
				if _, ok := execErr.(*exec.ExitError); !ok {
					fmt.Println("Execution error: " + execErr.Error())
					continue
				}
			}
		}

		writeOutput(fileName, output, errOutput, redirectStdError)
		fmt.Print("$ ")
	}
}

func writeOutput(fileName string, output string, errOutput string, redirectStdError bool) {
	if output == "" && errOutput == "" {
		return
	}

	isWritten := false

	if redirectStdError {
		isWritten = write(errOutput, fileName)
		if output != "" {
			isWritten = true
			fmt.Print(output)
		}
	} else {
		isWritten = write(output, fileName)
		if errOutput != "" {
			isWritten = true
			fmt.Print(errOutput)
		}
	}

	if strings.HasSuffix(output, "\n") || strings.HasSuffix(errOutput, "\n") || !isWritten {
		return
	}

	fmt.Println()
}

func write(data string, fileName string) bool {
	if fileName == "" {
		fmt.Print(data)
		return true
	}

	err := os.WriteFile(fileName, []byte(data), 0644)
	if err != nil {
		fmt.Print(err)
		return true
	}
	return false
}

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
	if slices.Contains(commands, input) {
		return input + " is a shell builtin", ""
	}

	path, err := exec.LookPath(input)
	if err == nil {
		return input + " is " + path, ""
	}

	return "", input + ": not found"
}
