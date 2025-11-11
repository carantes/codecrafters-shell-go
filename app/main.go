package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" and "os" imports in stage 1 (feel free to remove this!)
var _ = fmt.Fprint
var _ = os.Stdout

var builtInCommands = []string{"exit", "echo", "type", "pwd", "cd"}

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
		case "pwd":
			pwdCommand(argv)
		case "cd":
			cdCommand(argv)
		default:
			runExternalCommand(argv)
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

func pwdCommand(_ []string) {
	pwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	fmt.Println(pwd)
}

func cdCommand(argv []string) {
	if len(argv) < 2 {
		// No argument, print current directory
		pwdCommand([]string{})
		return
	}

	currentDir, _ := os.Getwd()
	fullPath := filepath.Join(currentDir, argv[1])

	if exists := findDirectory(fullPath); exists {
		os.Chdir(fullPath)
	} else {
		fmt.Fprintf(os.Stdout, "cd: %s: No such file or directory\n", fullPath)
	}
}

func runExternalCommand(argv []string) {
	if len(argv) < 1 {
		return
	}

	cmd := argv[0]
	args := argv[1:]

	if _, exists := findFileInPath(cmd); exists {
		// Run the command and share standard output
		cmd := exec.Command(cmd, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running command: %s\n", err)
		}
	} else {
		// Command not found
		fmt.Println(cmd + ": command not found")
	}
}

func findFileInPath(command string) (string, bool) {
	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, string(os.PathListSeparator))

	for _, dir := range paths {
		fullPath := dir + "/" + command
		if fileInfo, err := os.Stat(fullPath); err == nil && isExecutable(fileInfo) {
			return fullPath, true
		}
	}

	return "", false
}

func findDirectory(targetDir string) bool {
	if _, err := os.ReadDir(targetDir); err == nil {
		return true
	}

	return false
}

func isExecutable(fileInfo os.FileInfo) bool {
	mode := fileInfo.Mode()
	return mode.IsRegular() && mode&0111 != 0
}
