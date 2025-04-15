package main

import (
	"bufio"
	"fmt"
	"gotracker/cli"
	"gotracker/structs"
	"gotracker/utils"
	"os"
	"strconv"
	"strings"
	"time"
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
				fmt.Println("  - search: Search for a food item")
				fmt.Println("  - add: Add a consumed food item")
				fmt.Println("  - report: Generate a daily nutrition report")

			case "bodyfat":
				fmt.Printf("Body Fat: %.2f%%\n", user.BodyFat)

			case "imc":
				// Calculate IMC (Body Mass Index)
				imc := float64(user.Weight) / (float64(user.Height) / 100 * float64(user.Height) / 100)
				fmt.Printf("IMC: %.2f\n", imc)

			case "search":
				fmt.Print("Enter food name to search: ")
				var foodName string
				fmt.Scanln(&foodName)
				food, err := db.SearchFood(database, foodName)
				if err != nil {
					fmt.Println("Error:", err)
				} else {
					fmt.Printf("Found food: %+v\n", food)
				}

			case "add":
				fmt.Print("Enter food name: ")
				reader := bufio.NewReader(os.Stdin)
				foodName, _ := reader.ReadString('\n')
				foodName = strings.TrimSpace(foodName)

				fmt.Print("Enter quantity (grams): ")
				quantityInput, _ := reader.ReadString('\n')
				quantityInput = strings.TrimSpace(quantityInput)

				quantity, err := strconv.ParseFloat(quantityInput, 64)
				if err != nil {
					fmt.Println("Invalid quantity. Please enter a number.")
					break
				}

				err = db.AddFoodConsumption(database, user.ID, foodName, quantity, time.Now())
				if err != nil {
					fmt.Println("Error:", err)
				} else {
					fmt.Println("Food consumption added successfully.")
				}

			case "report":
				report, err := db.GenerateDailyReport(database, user.ID, time.Now())
				if err != nil {
					fmt.Println("Error:", err)
				} else {
					fmt.Println("Daily Nutrition Report:")
					fmt.Println(report)
				}

			case "exit":
				fmt.Println("Shutting down the application...")
				return

			default:
				fmt.Println("Unknown command. Type 'help' for a list of commands.")
			}
		}
	}
}
