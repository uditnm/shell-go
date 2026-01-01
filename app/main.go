package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"strings"
)

var _ = fmt.Print

func main() {
	reader := bufio.NewReader(os.Stdin)

mainloop:
	for {
		fmt.Print("$ ")
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
		appendOutput := false
		for i := 1; i < len(parts); i++ {
			if slices.Contains(RedirectOutputCommands(), parts[i]) {
				if i+1 >= len(parts) {
					fmt.Print("Error: No file given")
					continue mainloop
				}

				redirectStdError = parts[i] == ErrorRedirect || parts[i] == ErrorAppend
				appendOutput = parts[i] == OutputAppend || parts[i] == StandardOutputAppend || parts[i] == ErrorAppend

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

			//writing output to a buffer
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

		writeOutput(fileName, output, errOutput, redirectStdError, appendOutput)
	}
}

func writeOutput(fileName string, output string, errOutput string, redirectStdError bool, appendOutput bool) {
	if output == "" && errOutput == "" {
		return
	}

	isWrittenToTerminal := false

	if redirectStdError {
		isWrittenToTerminal = write(errOutput, fileName, appendOutput)
		if output != "" {
			isWrittenToTerminal = true
			fmt.Print(output)
		}
	} else {
		isWrittenToTerminal = write(output, fileName, appendOutput)
		if errOutput != "" {
			isWrittenToTerminal = true
			fmt.Print(errOutput)
		}
	}

	if strings.HasSuffix(output, "\n") || strings.HasSuffix(errOutput, "\n") || !isWrittenToTerminal {
		return
	}

	fmt.Println()
}

func write(data string, fileName string, appendOutput bool) bool {
	if fileName == "" {
		fmt.Print(data)
		return true
	}

	writeMode := os.O_TRUNC
	if appendOutput {
		writeMode = os.O_APPEND
	}

	file, err := os.OpenFile(fileName,
		writeMode|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	if err != nil {
		fmt.Print(err)
		return true
	}
	defer file.Close()

	//extra logic to make sure append is in newline
	if appendOutput {
		rf, _ := os.Open(fileName)
		defer rf.Close()
		info, _ := rf.Stat()
		if info.Size() > 0 {
			rf.Seek(-1, io.SeekEnd)
			b := make([]byte, 1)
			rf.Read(b)
			if b[0] != '\n' {
				data = "\n" + data
			}
		}
	}

	_, writeErr := file.WriteString(data)
	if writeErr != nil {
		fmt.Print(writeErr)
		return true
	}

	return false
}
