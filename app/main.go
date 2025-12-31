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

		parts, quoteErr := getTokens(input)

		if quoteErr != nil {
			fmt.Println("Error reading command:", quoteErr)
			continue
		}

		fileName := ""
		lastIdx := -1
		for i := 1; i < len(parts); i++ {
			if parts[i] == ">" || parts[i] == "1>" {
				if i+1 >= len(parts) {
					fmt.Println("Error: No file given")
				}

				lastIdx = i
				fileName = strings.Join(parts[(i+1):], " ")
				break
			}
		}

		if lastIdx != -1 {
			parts = parts[:lastIdx]
		}

		output := ""

		switch parts[0] {
		case Exit:
			return
		case Echo:
			output = strings.Join(parts[1:], " ")
		case Type:
			output = checkCommand(parts[1])
		case Pwd:
			output = getPresentWorkingDirectory()
		case Cd:
			changeDirectory(parts[1])
		default:
			_, err := exec.LookPath(parts[0])
			if err != nil {
				fmt.Println(parts[0] + ": command not found")
				continue
			}

			cmd := exec.Command(parts[0], parts[1:]...)
			cmd.Args = parts
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr

			if fileName == "" {
				cmd.Stdout = os.Stdout
			} else {
				outFile, err := os.OpenFile(
					fileName,
					os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
					0644,
				)
				if err != nil {
					fmt.Println("failed to open output file: ", err)
					continue
				}
				defer outFile.Close()

				cmd.Stdout = outFile
			}

			execErr := cmd.Run()
			if execErr != nil {
				if _, ok := execErr.(*exec.ExitError); !ok {
					fmt.Println("Execution error: ", execErr)
				}
			}

			continue
		}

		writeOutput(fileName, output)
	}
}

func writeOutput(fileName string, output string) {
	if output == "" {
		return
	}

	if fileName == "" {
		fmt.Println(output)
	} else {
		err := os.WriteFile(fileName, []byte(output+"\n"), 0644)
		if err != nil {
			fmt.Println(err)
		}
	}
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

func getPresentWorkingDirectory() string {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	return path
}

func changeDirectory(path string) {
	if path == HomeDirectory {
		home, _ := os.UserHomeDir()
		path = home
	}

	err := os.Chdir(path)
	if err != nil {
		fmt.Println(path + ": No such file or directory")
	}
}

func checkCommand(input string) string {
	if slices.Contains(commands, input) {
		return input + " is a shell builtin"
	}

	path, err := exec.LookPath(input)
	if err == nil {
		return input + " is " + path
	}

	fmt.Println(input + ": not found")
	return ""
}
