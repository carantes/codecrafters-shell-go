package main

import (
	"bufio"
	"fmt"
	"os"
)

// Ensures gofmt doesn't remove the "fmt" and "os" imports in stage 1 (feel free to remove this!)
var _ = fmt.Fprint
var _ = os.Stdout

func main() {

	for {
		// Start a shell prompt
		_, _ = fmt.Fprint(os.Stdout, "$ ")

		//	Read user input
		cmd, err := bufio.NewReader(os.Stdin).ReadString('\n')

		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}

		//	Print command not found message
		fmt.Println(cmd[:len(cmd)-1] + ": command not found")
	}
}
