package cli

import (
	"fmt"
	"bufio"
	"os"
	"gotracker/structs"
)

func Open() {
	// Créer un utilisateur fictif
	user := suser.SUser{
		ID:           1,
		Firstname:    "John",
		Lastname:     "Doe",
		Age:          30,
		Weight:       70,  // Poids en kg
		Height:       175, // Taille en cm
		BodyFat:      15.5,
		TargetWeight: 65,
	}

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
			fmt.Println("  - bodyfat: Display the user's body fat percentage")
            fmt.Println("  - imc: Display the user's Body Mass Index (IMC)")
		
		case "bodyfat\n":
            fmt.Printf("Body Fat: %.2f%%\n", user.GetBodyFat())

        case "imc\n":
            fmt.Printf("IMC: %.2f\n", user.GetIMC())
	
		default:
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}
	}

	fmt.Println("Exiting CLI...")
	os.Exit(0)
}