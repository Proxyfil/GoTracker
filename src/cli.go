package cli

import (
	"fmt"
	"bufio"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Command : ")
	command, _ := reader.ReadString('\n')
	fmt.Print(command)
}