package main

import (
	"fmt"
	"gotracker/cli"
	"gotracker/utils"
	"gotracker/structs"
)

func main() {
	// Connect to the database
	database, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}

	// Migrate the database
	err = db.Migrate(database)
	if err != nil {
		panic(err)
	}

	// Create the user
	user := suser.SUser{
		ID:           1,
		Firstname:    "John",
		Lastname:     "Doe",
		Age:          30,
		Weight:       70,  // kg
		Height:       175, // cm
		BodyFat:      15.5,
		TargetWeight: 65,
	}

	// Create a channel to send commands
	commandChannel := make(chan cli.CommandMessage)

	// Start the CLI in a goroutine
	go cli.Open(commandChannel)

	// Process commands from the channel
	for {
		select {
		case msg := <-commandChannel:
			switch msg.Command {
			case "help":
				fmt.Println("Available commands:")
				fmt.Println("  - exit: Exit the CLI")
				fmt.Println("  - help: Show this help message")
				fmt.Println("  - bodyfat: Display the user's body fat percentage")
				fmt.Println("  - imc: Display the user's Body Mass Index (IMC)")

			case "bodyfat":
				fmt.Printf("Body Fat: %.2f%%\n", user.BodyFat)

			case "imc":
				// Calculate IMC (Body Mass Index)
				imc := float64(user.Weight) / (float64(user.Height)/100 * float64(user.Height)/100)
				fmt.Printf("IMC: %.2f\n", imc)

			case "exit":
				fmt.Println("Shutting down the application...")
				return

			default:
				fmt.Println("Unknown command. Type 'help' for a list of commands.")
			}
		}
	}
}
