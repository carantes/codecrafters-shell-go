package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

type CommandFunc func([]string) (string, error)

// AllCommands maps command names to their corresponding functions.
var AllCommands = map[string]CommandFunc{
	"exit": exitCommand,
	"echo": echoCommand,
	"type": typeCommand,
	"cat":  catCommand,
	"pwd":  pwdCommand,
	"cd":   cdCommand,
}

// Built-in command names
// TODO: This list should be generated dynamically from AllCommands map
var builtInCommands = []string{
	"exit",
	"echo",
	"type",
	"cat",
	"pwd",
	"cd",
}

func RunCommand(argv []string) (string, error) {
	if len(argv) < 1 {
		return "", nil
	}

	cmd := argv[0]

	if commandFunc, exists := AllCommands[cmd]; exists {
		return commandFunc(argv)
	}

	return externalCommand(argv)
}

// Built-in command implementations

func exitCommand(_ []string) (string, error) {
	code := 0
	os.Exit(code)
	return "", nil
}

func echoCommand(argv []string) (string, error) {
	if len(argv) > 1 {
		return strings.Join(argv[1:], " "), nil
	} else {
		return "", nil
	}
}

func typeCommand(argv []string) (string, error) {
	if len(argv) < 2 {
		return "", fmt.Errorf(`usage: type <command>`)
	}

	arg := argv[1]

	// Check if the command is a shell builtin
	if slices.Contains(builtInCommands, arg) {
		return fmt.Sprintf("%s is a shell builtin", arg), nil
	}

	// Find file in PATH
	if file, exists := findFileInPath(arg); exists {
		return fmt.Sprintf("%s is %s", file, arg), nil
	}

	// Command not found
	return fmt.Sprintf("type: %s: not found", arg), nil
}

func catCommand(argv []string) (string, error) {
	if len(argv) < 2 {
		return "", fmt.Errorf("usage: cat <path>")
	}

	output := ""

	for _, arg := range argv {
		if arg == "cat" {
			continue
		}

		content, err := os.ReadFile(arg)

		if err != nil {
			return "", fmt.Errorf("cat: %s: No such file or directory", arg)
		}

		output += string(content)
	}

	// Handle output redirection
	return output, nil
}

func pwdCommand(_ []string) (string, error) {
	dir, err := os.Getwd()

	if err != nil {
		return "", err
	}

	return dir, nil
}

func cdCommand(argv []string) (string, error) {
	if len(argv) < 2 {
		// No argument, print current directory
		return pwdCommand(argv)
	}

	targetDir := argv[1]

	var fullPath string

	if strings.HasPrefix(targetDir, "~") {
		homeDir, _ := os.UserHomeDir()
		fullPath = strings.Replace(targetDir, "~", homeDir, 1)
	} else if strings.HasPrefix(targetDir, "/") {
		fullPath = targetDir
	} else {
		cwd, _ := os.Getwd()
		fullPath = filepath.Join(cwd, targetDir)
	}

	stat, err := os.Stat(fullPath)

	if err != nil || !stat.IsDir() {
		return fmt.Sprintf("cd: %s: No such file or directory", fullPath), nil
	}

	os.Chdir(fullPath)
	return "", nil
}

// External commands

func externalCommand(argv []string) (string, error) {
	if len(argv) < 1 {
		return "", nil
	}

	cmd := argv[0]
	args := argv[1:]

	if _, exists := findFileInPath(cmd); exists {
		// Run the command and share standard output
		cmd := exec.Command(cmd, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("error running command: %s\n", err)
		}
	} else {
		// Command not found
		return "", fmt.Errorf("unknown command: %s\n", cmd)
	}

	return "", nil
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

func isExecutable(fileInfo os.FileInfo) bool {
	mode := fileInfo.Mode()
	return mode.IsRegular() && mode&0111 != 0
}
