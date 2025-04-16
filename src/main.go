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

				case "add":
					if(user.ID == 0){
						fmt.Println("Please login or register first.")
						fmt.Print("Command: ")
						continue
					}
					// We get arguments from the command
					addArgs := strings.Split(msg.Command, " ")
					if len(addArgs) < 2 {
						fmt.Println("Please provide type, can be 'food', 'meal' or 'day'")
						fmt.Print("Command: ")
						continue
					}
					addType := addArgs[1]
				
					switch addType {
						case "food":
							if len(addArgs) < 4 {
								fmt.Println("Please provide a food ID and quantity to add.")
								fmt.Print("Command: ")
								continue
							}
							foodID, err := strconv.Atoi(addArgs[2])
							if err != nil {
								fmt.Println("Error converting food ID to int:", err)
								fmt.Print("Command: ")
								continue
							}
							quantity, err := strconv.Atoi(addArgs[3])
							if err != nil {
								fmt.Println("Error converting quantity to int:", err)
								fmt.Print("Command: ")
								continue
							}

							// Save the food history to the database
							current_time := time.Now().Local()
							err = db.AddFoodHistory(database, user.ID, foodID, current_time.Format("2006-01-02"), quantity)
							if err != nil {
								fmt.Println("Error saving food history:", err)
								fmt.Print("Command: ")
								continue
							} else {
								fmt.Println("Food history saved successfully.")
							}
						
						case "meal":
							if len(addArgs) < 3 {
								fmt.Println("Please provide a meal ID to add.")
								fmt.Print("Command: ")
								continue
							}
							mealID, err := strconv.Atoi(addArgs[2])
							if err != nil {
								fmt.Println("Error converting meal ID to int:", err)
								fmt.Print("Command: ")
								continue
							}
							// Get food IDs associated with the meal
							foods, err := db.GetFoodWithMeal(database, mealID)
							if err != nil {
								fmt.Println("Error fetching food IDs with meal:", err)
								fmt.Print("Command: ")
								continue
							}
							if len(foods) == 0 {
								fmt.Println("No food IDs found for the specified meal.")
								fmt.Print("Command: ")
								continue
							}
							// Save the food history to the database
							current_time := time.Now().Local()
							for _, food := range foods{
								err = db.AddFoodHistory(database, user.ID, food[0], current_time.Format("2006-01-02"), food[1])
								if err != nil {
									fmt.Println("Error saving food history:", err)
									fmt.Print("Command: ")
									continue
								} else {
									fmt.Printf("Food history for food ID %d saved successfully.\n", food[1])
								}
							}

						case "day":
							if len(addArgs) < 3 {
								fmt.Println("Please provide a day ID to add.")
								fmt.Print("Command: ")
								continue
							}
							dayID, err := strconv.Atoi(addArgs[2])
							if err != nil {
								fmt.Println("Error converting day ID to int:", err)
								fmt.Print("Command: ")
								continue
							}
							// Get meal IDs associated with the day
							meals, err := db.GetMealWithDayPreset(database, dayID)
							if err != nil {
								fmt.Println("Error fetching meal IDs with day:", err)
								fmt.Print("Command: ")
								continue
							}
							if len(meals) == 0 {
								fmt.Println("No meal IDs found for the specified day.")
								fmt.Print("Command: ")
								continue
							}
							// Save the food history to the database
							current_time := time.Now().Local()
							for _, meal := range meals {
								foods, err := db.GetFoodWithMeal(database, meal[0])
								if err != nil {
									fmt.Println("Error fetching food IDs with meal:", err)
									fmt.Print("Command: ")
									continue
								}
								if len(foods) == 0 {
									fmt.Println("No food IDs found for the specified meal.")
									fmt.Print("Command: ")
									continue
								}
								for _, food := range foods {
									err = db.AddFoodHistory(database, user.ID, food[0], current_time.Format("2006-01-02"), food[1]*meal[1])
									if err != nil {
										fmt.Println("Error saving food history:", err)
										fmt.Print("Command: ")
										continue
									} else {
										fmt.Printf("Food history for food ID %d saved successfully.\n", food[0])
									}
								}
							}

						default:
							fmt.Println("Unknown add type. Please use 'food', 'meal' or 'day'.")
					}

				case "create":
					if(user.ID == 0){
						fmt.Println("Please login or register first.")
						fmt.Print("Command: ")
						continue
					}
					// We get arguments from the command
					createArgs := strings.Split(msg.Command, " ")
					if len(createArgs) < 2 {
						fmt.Println("Please provide a type, can be 'meal' or 'day'")
						fmt.Print("Command: ")
						continue
					}
					createType := createArgs[1]
					switch createType {
						case "meal":
							if len(createArgs) < 4 {
								fmt.Println("Please provide a meal name and type to create.")
								fmt.Print("Command: ")
								continue
							}
							mealName := createArgs[2]
							mealType := createArgs[3]
							// Create a new meal in the database
							mealID, err := db.CreateMeal(database, mealName, mealType)
							if err != nil {
								fmt.Println("Error creating meal:", err)
								fmt.Print("Command: ")
								continue
							}
							fmt.Printf("Meal '%s' created successfully with ID: %d\n", mealName, mealID)
						case "day":
							if len(createArgs) < 3 {
								fmt.Println("Please provide a day name to create.")
								fmt.Print("Command: ")
								continue
							}
							dayName := createArgs[2]
							// Create a new day in the database
							dayID, err := db.CreateDayPreset(database, user.ID, dayName)
							if err != nil {
								fmt.Println("Error creating day:", err)
								fmt.Print("Command: ")
								continue
							}
							fmt.Printf("Day '%s' created successfully with ID: %d\n", dayName, dayID)
						default:
							fmt.Println("Unknown create type. Please use 'meal' or 'day'.")
					}
				case "list":
					// We get arguments from the command
					listArgs := strings.Split(msg.Command, " ")
					if len(listArgs) < 2 {
						fmt.Println("Please provide a type, can be 'meal' or 'day'")
						fmt.Print("Command: ")
						continue
					}
					listType := listArgs[1]
					switch listType {
						case "meal":
							// List all meals in the database
							meals, err := db.GetAllMeals(database)
							if err != nil {
								fmt.Println("Error fetching meals:", err)
								fmt.Print("Command: ")
								continue
							}
							if len(meals) == 0 {
								fmt.Println("No meals found.")
							} else {
								fmt.Println("Meals:")
								for _, meal := range meals {
									fmt.Println(" -", meal)
								}
							}
						case "day":
							// List all days in the database
							days, err := db.GetAllDays(database)
							if err != nil {
								fmt.Println("Error fetching days:", err)
								fmt.Print("Command: ")
								continue
							}
							if len(days) == 0 {
								fmt.Println("No days found.")
							} else {
								fmt.Println("Days:")
								for _, day := range days {
									fmt.Println(" -", day)
								}
							}
						default:
							fmt.Println("Unknown list type. Please use 'meal' or 'day'.")
					}

				case "link":
					if(user.ID == 0){
						fmt.Println("Please login or register first.")
						fmt.Print("Command: ")
						continue
					}
					// We get arguments from the command
					linkArgs := strings.Split(msg.Command, " ")
					if len(linkArgs) < 2 {
						fmt.Println("Please provide a type, can be 'food_to_meal' or 'meal_to_day'")
						fmt.Print("Command: ")
						continue
					}
					linkType := linkArgs[1]
					switch linkType {
						case "food_to_meal":
							if len(linkArgs) < 5 {
								fmt.Println("Please provide a food ID, meal ID and quantity to link.")
								fmt.Print("Command: ")
								continue
							}
							foodID, err := strconv.Atoi(linkArgs[2])
							if err != nil {
								fmt.Println("Error converting food ID to int:", err)
								fmt.Print("Command: ")
								continue
							}
							mealID, err := strconv.Atoi(linkArgs[3])
							if err != nil {
								fmt.Println("Error converting meal ID to int:", err)
								fmt.Print("Command: ")
								continue
							}
							quantity, err := strconv.Atoi(linkArgs[4])
							if err != nil {
								fmt.Println("Error converting quantity to int:", err)
								fmt.Print("Command: ")
								continue
							}
							// Link food to meal in the database
							err = db.LinkFoodToMeal(database, foodID, mealID, quantity)
							if err != nil {
								fmt.Println("Error linking food to meal:", err)
								fmt.Print("Command: ")
								continue
							} else {
								fmt.Printf("Food ID %d linked to Meal ID %d successfully.\n", foodID, mealID)
							}
						case "meal_to_day":
							if len(linkArgs) < 5 {
								fmt.Println("Please provide a meal ID, day ID and quantity to link.")
								fmt.Print("Command: ")
								continue
							}
							mealID, err := strconv.Atoi(linkArgs[2])
							if err != nil {
								fmt.Println("Error converting meal ID to int:", err)
								fmt.Print("Command: ")
								continue
							}
							dayID, err := strconv.Atoi(linkArgs[3])
							if err != nil {
								fmt.Println("Error converting day ID to int:", err)
								fmt.Print("Command: ")
								continue
							}
							quantity, err := strconv.Atoi(linkArgs[4])
							if err != nil {
								fmt.Println("Error converting quantity to int:", err)
								fmt.Print("Command: ")
								continue
							}
							// Link meal to day in the database
							err = db.LinkMealToDayPreset(database, mealID, dayID, quantity)
							if err != nil {
								fmt.Println("Error linking meal to day:", err)
								fmt.Print("Command: ")
								continue
							} else {
								fmt.Printf("Meal ID %d linked to Day ID %d successfully.\n", mealID, dayID)
							}
						default:
							fmt.Println("Unknown link type. Please use 'food_to_meal' or 'meal_to_day'.")
					}

				case "history":
					if(user.ID == 0){
						fmt.Println("Please login or register first.")
						fmt.Print("Command: ")
						continue
					}
					// We get arguments from the command
					historyArgs := strings.Split(msg.Command, " ")
					if len(historyArgs) < 2 {
						fmt.Println("Please provide a type, can be 'food', 'weight', 'imc' or 'bodyfat'")
						fmt.Print("Command: ")
						continue
					}
					historyType := historyArgs[1]
					switch historyType {
						case "food":
							// Get food history for the user
							foodHistory, err := db.GetFoodHistory(database, user.ID)
							if err != nil {
								fmt.Println("Error fetching food history:", err)
								fmt.Print("Command: ")
								continue
							}
							if len(foodHistory) == 0 {
								fmt.Println("No food history found.")
							} else {
								fmt.Println("Food History:")
								for _, food := range foodHistory {
									fmt.Println(" - Date:", food[1], "| Food ID:", food[0], "| Quantity:", food[2], "| Entry ID:", food[3])
								}
							}
						case "weight":
							// Get weight history for the user
							weightHistory, err := db.GetWeightHistory(database, user.ID)
							if err != nil {
								fmt.Println("Error fetching weight history:", err)
								fmt.Print("Command: ")
								continue
							}
							if len(weightHistory) == 0 {
								fmt.Println("No weight history found.")
							} else {
								fmt.Println("Weight History:")
								for _, weight := range weightHistory {
									fmt.Println(" - Date:", weight[0], "| Weight:", weight[1])
								}
							}
						case "imc":
							// Get IMC history for the user
							imcHistory, err := db.GetIMCHistory(database, user.ID)
							if err != nil {
								fmt.Println("Error fetching IMC history:", err)
								fmt.Print("Command: ")
								continue
							}
							if len(imcHistory) == 0 {
								fmt.Println("No IMC history found.")
							} else {
								fmt.Println("IMC History:")
								for _, imc := range imcHistory {
									fmt.Println(" - Date:", imc[0], "| IMC:", imc[1])
								}
							}
						case "bodyfat":
							// Get body fat history for the user
							bodyFatHistory, err := db.GetBodyFatHistory(database, user.ID)
							if err != nil {
								fmt.Println("Error fetching body fat history:", err)
								fmt.Print("Command: ")
								continue
							}
							if len(bodyFatHistory) == 0 {
								fmt.Println("No body fat history found.")
							} else {
								fmt.Println("Body Fat History:")
								for _, bodyFat := range bodyFatHistory {
									fmt.Println(" - Date:", bodyFat[0], "| Body Fat:", bodyFat[1])
								}
							}
						default:
							fmt.Println("Unknown history type. Please use 'food', 'weight', 'imc' or 'bodyfat'.")
					}
				case "delete":
					if(user.ID == 0){
						fmt.Println("Please login or register first.")
						fmt.Print("Command: ")
						continue
					}
					// We get arguments from the command
					deleteArgs := strings.Split(msg.Command, " ")
					if len(deleteArgs) < 2 {
						fmt.Println("Please provide a type, can be 'food'.")
						fmt.Print("Command: ")
						continue
					}
					deleteType := deleteArgs[1]
					switch deleteType {
						case "food":
							if len(deleteArgs) < 3 {
								fmt.Println("Please provide a Entry ID to delete.")
								fmt.Print("Command: ")
								continue
							}
							entryID, err := strconv.Atoi(deleteArgs[2])
							if err != nil {
								fmt.Println("Error converting entry ID to int:", err)
								fmt.Print("Command: ")
								continue
							}
							// Delete food history from the database
							err = db.DeleteFoodHistory(database, entryID)
							if err != nil {
								fmt.Println("Error deleting food history:", err)
								fmt.Print("Command: ")
								continue
							} else {
								fmt.Printf("Food history with Entry ID %d deleted successfully.\n", entryID)
							}
						default:
							fmt.Println("Unknown delete type. Please use 'food'.")
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
