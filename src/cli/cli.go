package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type CommandMessage struct {
	Command string
}

func Open(commands chan<- CommandMessage) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Command: ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		// Clean input (remove newline, spaces, etc.)
		input = strings.TrimSpace(input)

		// Send the command to the channel
		commands <- CommandMessage{Command: input}
	}
}
