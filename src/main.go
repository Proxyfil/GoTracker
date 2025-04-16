package main

import (
	"fmt"
	"gotracker/cli"
	"gotracker/fdcnal"
	"gotracker/utils"
	"gotracker/structs"
	"strings"
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
			switch strings.Split(msg.Command," ")[0] {
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

				case "search":
					// We get arguments from the command
					searchArgs := strings.Split(msg.Command, " ")
					if len(searchArgs) < 2 {
						fmt.Println("Please provide a search type, can be 'meal', 'food' or 'day'")
						fmt.Print("Command: ")
						continue
					}
					searchType := searchArgs[1]

					switch searchType {
						case "food":
							if len(searchArgs) < 3 {
								fmt.Println("Please provide a food name to search for.")
								fmt.Print("Command: ")
								continue
							}
							foodName := searchArgs[2]
							foods, err := api.GetFoodByName(foodName)
							if err != nil {
								fmt.Println("Error fetching food data:", err)
								fmt.Print("Command: ")
								continue
							}
							fmt.Println("Food Search Results:")
							for _, food := range foods {
								fmt.Println(" -", food)
							}
							
						default:
							fmt.Println("Unknown search type. Please use 'meal', 'food' or 'day'.")
					}


				case "exit":
					fmt.Println("Shutting down the application...")
					return

				default:
					fmt.Println("Unknown command. Type 'help' for a list of commands.")
				}
		}
		fmt.Print("Command: ")
	}
}
