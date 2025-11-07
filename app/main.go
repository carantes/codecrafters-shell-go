package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" and "os" imports in stage 1 (feel free to remove this!)
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

		cmd := strings.Fields(input)[0]
		args := strings.Fields(input)[1:]

		switch cmd {
		case "exit":
			os.Exit(0)
		case "echo":
			if len(args) > 0 {
				fmt.Println(strings.Join(args, " "))
			}
		default:
			//	Print command not found message
			fmt.Println(cmd + ": command not found")
		}

	}
}
