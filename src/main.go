package main

import (
	"fmt"
	"gotracker/cli"
	"gotracker/fdcnal"
	"gotracker/utils"
	"gotracker/structs"
	"strings"
	"time"
	"strconv"
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

	// Initialize the user variable
	var user suser.SUser
	// Set the user ID to 0 (not logged in)
	user.ID = 0

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
					fmt.Println("  - search: Search for food or meal")
					fmt.Println("  - report: Generate a report")
					fmt.Println("  - register: Register a new user")
					fmt.Println("  - login: Login as an existing user")

				case "bodyfat":
					if(user.ID == 0){
						fmt.Println("Please login or register first.")
						fmt.Print("Command: ")
						continue
					}
					fmt.Printf("Body Fat: %.2f%%\n", user.BodyFat)

				case "imc":
					if(user.ID == 0){
						fmt.Println("Please login or register first.")
						fmt.Print("Command: ")
						continue
					}
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
							foods, err := fdcnal.GetFoodByName(foodName)
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
				case "report":
					if(user.ID == 0){
						fmt.Println("Please login or register first.")
						fmt.Print("Command: ")
						continue
					}
					// We get arguments from the command
					reportArgs := strings.Split(msg.Command, " ")
					if len(reportArgs) < 2 {
						fmt.Println("Please provide a report type, can be 'imc', 'weight' or 'bodyfat'")
						fmt.Print("Command: ")
						continue
					}

					reportType := reportArgs[1]
					switch reportType {
						case "imc":
							// Generate IMC report
							imc := user.GetIMC()
							fmt.Printf("IMC Report: %.2f\n", imc)
							// Save IMC history to the database
							current_time := time.Now().Local()
							err := db.CreateIMCHistory(database, user.ID, current_time.Format("2006-01-02"), imc)
							if err != nil {
								fmt.Println("Error saving IMC history:", err)
							} else {
								fmt.Println("IMC history saved successfully.")
							}

						case "bodyfat":
							// Generate body fat report
							bodyFat := user.GetBodyFat()
							fmt.Printf("Body Fat Report: %.2f%%\n", bodyFat)
							// Save body fat history to the database
							current_time := time.Now().Local()
							err := db.CreateBodyFatHistory(database, user.ID, current_time.Format("2006-01-02"), bodyFat)
							if err != nil {
								fmt.Println("Error saving body fat history:", err)
							} else {
								fmt.Println("Body fat history saved successfully.")
							}
						case "weight":
							// Generate weight report
							weight := user.Weight
							fmt.Printf("Weight Report: %d kg\n", weight)
							// Save weight history to the database
							current_time := time.Now().Local()
							err := db.CreateWeightHistory(database, user.ID, current_time.Format("2006-01-02"), weight)
							if err != nil {
								fmt.Println("Error saving weight history:", err)
							} else {
								fmt.Println("Weight history saved successfully.")
							}
						default:
							fmt.Println("Unknown report type. Please use 'imc', 'weight' or 'bodyfat'.")
					}
				case "register":
					// We get arguments from the command
					registerArgs := strings.Split(msg.Command, " ")
					if( len(registerArgs) < 7) {
						fmt.Println("Please provide a firstname, lastname, age, weight, height and target weight")
						fmt.Print("Command: ")
						continue
					}
					
					firstname := registerArgs[1]
					lastname := registerArgs[2]
					age := registerArgs[3]
					weight := registerArgs[4]
					height := registerArgs[5]
					targetWeight := registerArgs[6]

					// Convert age, weight, height and target weight to int
					ageInt, err := strconv.Atoi(age)
					if err != nil {
						fmt.Println("Error converting age to int:", err)
						fmt.Print("Command: ")
						continue
					}
					weightInt, err := strconv.Atoi(weight)
					if err != nil {
						fmt.Println("Error converting weight to int:", err)
						fmt.Print("Command: ")
						continue
					}
					heightInt, err := strconv.Atoi(height)
					if err != nil {
						fmt.Println("Error converting height to int:", err)
						fmt.Print("Command: ")
						continue
					}
					targetWeightInt, err := strconv.Atoi(targetWeight)
					if err != nil {
						fmt.Println("Error converting target weight to int:", err)
						fmt.Print("Command: ")
						continue
					}
					// Create a new user
					newUser := suser.SUser{
						Firstname:    firstname,
						Lastname:     lastname,
						Age:          ageInt,
						Weight:       weightInt,
						Height:       heightInt,
						TargetWeight: targetWeightInt,
					}
					// Save the user to the database
					userID, err := db.CreateUser(database, newUser.Firstname, newUser.Lastname, newUser.Age, newUser.Weight, newUser.Height, newUser.TargetWeight)
					if err != nil {
						fmt.Println("Error creating user:", err)
						fmt.Print("Command: ")
						continue
					}
					fmt.Printf("User created successfully with ID: %d\n", userID)
					// Set the user ID
					newUser.ID = userID
					// Set the user in the main user variable
					user = newUser
					fmt.Printf("User %s %s registered successfully with ID: %d\n", user.Firstname, user.Lastname, user.ID)

				case "login":
					// We get arguments from the command
					loginArgs := strings.Split(msg.Command, " ")
					if len(loginArgs) < 2 {
						fmt.Println("Please provide a user ID to login")
						fmt.Print("Command: ")
						continue
					}
					userID, err := strconv.Atoi(loginArgs[1])
					if err != nil {
						fmt.Println("Error converting user ID to int:", err)
						fmt.Print("Command: ")
						continue
					}
					// Check if the user exists in the database
					id, firstname, lastname, age, weight, height, target_weight, err := db.GetUser(database, userID)
					if err != nil {
						fmt.Println("Error checking user existence:", err)
						fmt.Print("Command: ")
						continue
					}

					if id == 0 {
						fmt.Println("User not found. Please register first.")
						fmt.Print("Command: ")
						continue
					}

					user = suser.SUser{
						ID:           id,
						Firstname:    firstname,
						Lastname:     lastname,
						Age:          age,
						Weight:       weight,
						Height:       height,
						TargetWeight: target_weight,
					}
					fmt.Printf("User %s %s logged in successfully with ID: %d\n", user.Firstname, user.Lastname, user.ID)

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
