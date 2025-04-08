package cli

import (
	"fmt"
	"bufio"
	"os"
)

func Open() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Command : ");
		command, _ := reader.ReadString('\n');
		fmt.Print(command);

		if command == "exit\n" {
			break
		}
	}

	fmt.Println("Exiting CLI...")
	os.Exit(0)
}