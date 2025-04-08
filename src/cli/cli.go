package cli

import (
	"fmt"
	"bufio"
	"os"
)

func Open() {
	// Define the reader to read from standard input
	reader := bufio.NewReader(os.Stdin)

	// While loop to keep the CLI running
	for {
		fmt.Print("Command : ");
		// Read the command from standard input
		command, _ := reader.ReadString('\n');

		// Switch between commands
		switch command {
		case "exit\n":
			fmt.Println("Exiting CLI...")
			os.Exit(0)

		case "help\n":
			fmt.Println("Available commands:")
			fmt.Println("  - exit: Exit the CLI")
			fmt.Println("  - help: Show this help message")

		default:
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}
	}

	// Close the reader
	defer reader.Close()
	fmt.Println("Exiting CLI...")
	os.Exit(0)
}