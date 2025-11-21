package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// // Ensures gofmt doesn't remove the "fmt" and "os" imports in stage 1 (feel free to remove this!)
var _ = fmt.Fprint
var _ = os.Stdout

func main() {

	for {
		// Start a shell prompt
		_, _ = fmt.Fprint(os.Stdout, "$ ")

		//	Read user input
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')

		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		if len(strings.TrimSpace(input)) == 0 {
			continue
		}

		// Custom parser to handle quotes and escapes
		argv := parseInput(input)

		// Get redirection arguments if any
		outFile, cleanArgv := getRedirection(argv)

		output, err := RunCommand(cleanArgv)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		printOutput(output, outFile)
	}
}

func parseInput(input string) []string {
	var argv []string
	var currentArg strings.Builder
	var quoteChar rune
	escapeNext := false

	for _, char := range input {
		if escapeNext {
			currentArg.WriteRune(char)
			escapeNext = false
			continue
		}

		switch char {
		case '\\':
			escapeNext = true
		case '\'', '"':
			if quoteChar == 0 {
				// Starting a new quote
				quoteChar = char
			} else if quoteChar == char {
				// Ending the current quote
				quoteChar = 0
			} else {
				// Inside one quote type, treat the other as literal
				currentArg.WriteRune(char)
			}
		case ' ', '\n', '\t':
			if quoteChar != 0 {
				currentArg.WriteRune(char)
			} else if currentArg.Len() > 0 {
				argv = append(argv, currentArg.String())
				currentArg.Reset()
			}
		default:
			currentArg.WriteRune(char)
		}
	}

	if currentArg.Len() > 0 {
		argv = append(argv, currentArg.String())
	}

	return argv
}

func getRedirection(argv []string) (outFile string, cleanedArgv []string) {
	cleanedArgv = []string{}
	i := 0
	for i < len(argv) {
		if (argv[i] == ">" || argv[i] == "1>") && i+1 < len(argv) {
			outFile = argv[i+1]
			i += 2 // Skip the redirection and file name
		} else {
			cleanedArgv = append(cleanedArgv, argv[i])
			i++
		}
	}
	return
}

func printOutput(output, outFile string) error {
	if outFile != "" {
		// Handle output redirection
		file, err := os.Create(outFile)
		if err != nil {
			return fmt.Errorf("error creating file %s: %v", outFile, err)
		}

		defer file.Close()

		_, err = file.WriteString(output)
		if err != nil {
			return fmt.Errorf("error writing to file %s: %v", outFile, err)
		}
	} else {
		// Print to standard output
		fmt.Println(output)
	}

	return nil
}
