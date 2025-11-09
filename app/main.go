package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" and "os" imports in stage 1 (feel free to remove this!)
var _ = fmt.Fprint
var _ = os.Stdout

var builtInCommands = []string{"exit", "echo", "type"}

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

		argv := strings.Fields(input)
		cmd := argv[0]

		switch cmd {
		case "exit":
			exitCommand(argv)
		case "echo":
			echoCommand(argv)
		case "type":
			typeCommand(argv)
		default:
			fmt.Println(cmd + ": command not found")
		}

	}
}

func exitCommand(_ []string) {
	code := 0
	os.Exit(code)
}

func echoCommand(argv []string) {
	if len(argv) > 1 {
		fmt.Println(strings.Join(argv[1:], " "))
	} else {
		fmt.Println()
	}
}

func typeCommand(argv []string) {
	// Not enough arguments
	if len(argv) < 2 {
		return
	}

	arg := argv[1]

	// Check if the command is a shell builtin
	if slices.Contains(builtInCommands, arg) {
		fmt.Println(arg + " is a shell builtin")
		return
	}

	// Find file in PATH
	if file, exists := findFileInPath(arg); exists {
		fmt.Fprintf(os.Stdout, "%s is %s\n", arg, file)
		return
	}

	// Command not found
	fmt.Println(arg + ": not found")
}

func findFileInPath(command string) (string, bool) {
	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, string(os.PathListSeparator))

	for _, dir := range paths {
		fullPath := dir + "/" + command
		if fileInfo, err := os.Stat(fullPath); err == nil && !fileInfo.IsDir() {
			return fullPath, true
		}
	}

	return "", false
}
