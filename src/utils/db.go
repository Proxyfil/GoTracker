package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type Config struct {
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBHost     string `json:"db_host"`
	DBPort     string `json:"db_port"`
	DBName     string `json:"db_name"`
}

// LoadConfig reads the configuration from a JSON file
func LoadConfig() (*Config, error) {
	file, err := os.Open("./vars/config.json")
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return &config, nil
}

// ConnectDB establishes a connection to the postgres database
func ConnectDB() (*sql.DB, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	dsn := fmt.Sprintf("user=%s dbname=%s sslmode=disable host=%s port=%s password=%s",
		config.DBUser, config.DBName, config.DBHost, config.DBPort, config.DBPassword)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func Migrate(db *sql.DB) error {
	// Create the users table if it doesn't exist
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			firstname VARCHAR(50),
			lastname VARCHAR(50),
			age INT,
			weight INT,
			height INT,
			target_weight INT
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Create body_fat_history table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS body_fat_history (
			id SERIAL PRIMARY KEY,
			user_id INT REFERENCES users(id),
			date DATE,
			body_fat FLOAT
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create body_fat_history table: %w", err)
	}

	// Create imc_history table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS imc_history (
			id SERIAL PRIMARY KEY,
			user_id INT REFERENCES users(id),
			date DATE,
			imc FLOAT
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create imc_history table: %w", err)
	}

	// Create weight_history table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS weight_history (
			id SERIAL PRIMARY KEY,
			user_id INT REFERENCES users(id),
			date DATE,
			weight INT
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create weight_history table: %w", err)
	}

	// Create meal table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS meal (
			id SERIAL PRIMARY KEY,
			name VARCHAR(50),
			type VARCHAR(50)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create meal table: %w", err)
	}

	// Create meal_food table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS meal_food (
			id SERIAL PRIMARY KEY,
			meal_id INT REFERENCES meal(id),
			food_id INT,
			quantity FLOAT
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create meal_food table: %w", err)
	}

	// Create food_history table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS food_history (
			id SERIAL PRIMARY KEY,
			user_id INT REFERENCES users(id),
			food_id INT,
			date DATE,
			quantity FLOAT
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create food_history table: %w", err)
	}

	// Create day_preset table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS day_preset (
			id SERIAL PRIMARY KEY,
			user_id INT REFERENCES users(id),
			name VARCHAR(50)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create day_preset table: %w", err)
	}

	// Create day_preset_meal table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS day_preset_meal (
			id SERIAL PRIMARY KEY,
			day_preset_id INT REFERENCES day_preset(id),
			meal_id INT REFERENCES meal(id),
			quantity FLOAT
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create day_preset_meal table: %w", err)
	}

	// Create target table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS target (
			id SERIAL PRIMARY KEY,
			user_id INT REFERENCES users(id),
			date DATE,
			fat FLOAT,
			saturated_fat FLOAT,
			trans_fat FLOAT,
			cholesterol FLOAT,
			sodium FLOAT,
			carbohydrates FLOAT,
			fiber FLOAT,
			sugars FLOAT,
			protein FLOAT,
			calcium FLOAT,
			iron FLOAT,
			potassium FLOAT,
			calories FLOAT
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create target table: %w", err)
	}

	return nil
}

func CreateIMCHistory(db *sql.DB, userID int, date string, imc float64) error {
	_, err := db.Exec(`
		INSERT INTO imc_history (user_id, date, imc)
		VALUES ($1, $2, $3)
	`, userID, date, imc)
	if err != nil {
		return fmt.Errorf("failed to insert IMC history: %w", err)
	}
	return nil
}

func CreateBodyFatHistory(db *sql.DB, userID int, date string, bodyFat float64) error {
	_, err := db.Exec(`
		INSERT INTO body_fat_history (user_id, date, body_fat)
		VALUES ($1, $2, $3)
	`, userID, date, bodyFat)
	if err != nil {
		return fmt.Errorf("failed to insert body fat history: %w", err)
	}
	return nil
}

func CreateWeightHistory(db *sql.DB, userID int, date string, weight int) error {
	_, err := db.Exec(`
		INSERT INTO weight_history (user_id, date, weight)
		VALUES ($1, $2, $3)
	`, userID, date, weight)
	if err != nil {
		return fmt.Errorf("failed to insert weight history: %w", err)
	}
	return nil
}

func CreateUser(db *sql.DB, firstname string, lastname string, age int, weight int, height int, targetWeight int) (int, error) {
	var userID int
	err := db.QueryRow(`
		INSERT INTO users (firstname, lastname, age, weight, height, target_weight)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, firstname, lastname, age, weight, height, targetWeight).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}
	return userID, nil
}

func GetUser(db *sql.DB, userID int) (int, string, string, int, int, int, int, error) {
	var firstname, lastname string
	var id, age, weight, height, targetWeight int
	err := db.QueryRow(`
		SELECT id, firstname, lastname, age, weight, height, target_weight
		FROM users
		WHERE id = $1
	`, userID).Scan(&id, &firstname, &lastname, &age, &weight, &height, &targetWeight)
	if err != nil {
		return 0, "", "", 0, 0, 0, 0, fmt.Errorf("failed to get user: %w", err)
	}
	return id, firstname, lastname, age, weight, height, targetWeight, nil
}

func AddFoodHistory(db *sql.DB, userID int, foodID int, date string, quantity int) error {
	_, err := db.Exec(`
		INSERT INTO food_history (user_id, food_id, date, quantity)
		VALUES ($1, $2, $3, $4)
	`, userID, foodID, date, quantity)
	if err != nil {
		return fmt.Errorf("failed to insert food history: %w", err)
	}
	return nil
}

func GetFoodWithMeal(db *sql.DB, mealID int) ([][2]int, error) {
	rows, err := db.Query(`
		SELECT food_id, quantity
		FROM meal_food
		WHERE meal_id = $1
	`, mealID)
	if err != nil {
		return nil, fmt.Errorf("failed to get food with meal: %w", err)
	}
	defer rows.Close()

	var foods [][2]int
	for rows.Next() {
		var foodID int
		var quantity int
		if err := rows.Scan(&foodID, &quantity); err != nil {
			return nil, fmt.Errorf("failed to scan food ID: %w", err)
		}
		foods = append(foods, [2]int{foodID, quantity})
	}

	return foods, nil
}

func GetMealWithDayPreset(db *sql.DB, dayPresetID int) ([][2]int, error) {
	rows, err := db.Query(`
		SELECT meal_id, quantity
		FROM day_preset_meal
		WHERE day_preset_id = $1
	`, dayPresetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get meal with day preset: %w", err)
	}
	defer rows.Close()

	var meals [][2]int
	for rows.Next() {
		var mealID int
		var quantity int
		if err := rows.Scan(&mealID, &quantity); err != nil {
			return nil, fmt.Errorf("failed to scan meal ID: %w", err)
		}
		meals = append(meals, [2]int{mealID, quantity})
	}

	return meals, nil
}

func CreateMeal(db *sql.DB, name string, mealType string) (int, error) {
	var mealID int
	err := db.QueryRow(`
		INSERT INTO meal (name, type)
		VALUES ($1, $2)
		RETURNING id
	`, name, mealType).Scan(&mealID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert meal: %w", err)
	}
	return mealID, nil
}

func CreateDayPreset(db *sql.DB, userID int, name string) (int, error) {
	var dayPresetID int
	err := db.QueryRow(`
		INSERT INTO day_preset (user_id, name)
		VALUES ($1, $2)
		RETURNING id
	`, userID, name).Scan(&dayPresetID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert day preset: %w", err)
	}
	return dayPresetID, nil
}

func GetAllMeals(db *sql.DB) ([]string, error) {
	rows, err := db.Query(`
		SELECT id, name, type
		FROM meal
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get all meals: %w", err)
	}
	defer rows.Close()

	var meals []string
	for rows.Next() {
		var id int
		var name, mealType string
		if err := rows.Scan(&id, &name, &mealType); err != nil {
			return nil, fmt.Errorf("failed to scan meal: %w", err)
		}
		meals = append(meals, fmt.Sprintf("ID: %d | Name: %s | Type: %s", id, name, mealType))
	}

	return meals, nil
}

func GetAllDays(db *sql.DB) ([]string, error) {
	rows, err := db.Query(`
		SELECT id, name
		FROM day_preset
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get all day presets: %w", err)
	}
	defer rows.Close()

	var dayPresets []string
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, fmt.Errorf("failed to scan day preset: %w", err)
		}
		dayPresets = append(dayPresets, fmt.Sprintf("ID: %d | Name: %s", id, name))
	}

	return dayPresets, nil
}

func LinkFoodToMeal(db *sql.DB, foodID int, mealID int, quantity int) error {
	_, err := db.Exec(`
		INSERT INTO meal_food (meal_id, food_id, quantity)
		VALUES ($1, $2, $3)
	`, mealID, foodID, quantity)
	if err != nil {
		return fmt.Errorf("failed to link food to meal: %w", err)
	}
	return nil
}

func LinkMealToDayPreset(db *sql.DB, mealID int, dayPresetID int, quantity int) error {
	_, err := db.Exec(`
		INSERT INTO day_preset_meal (day_preset_id, meal_id, quantity)
		VALUES ($1, $2, $3)
	`, dayPresetID, mealID, quantity)
	if err != nil {
		return fmt.Errorf("failed to link meal to day preset: %w", err)
	}
	return nil
}

func GetFoodHistory(db *sql.DB, userID int) ([][4]interface{}, error) {
	rows, err := db.Query(`
		SELECT food_id, date, quantity, id
		FROM food_history
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get food history: %w", err)
	}
	defer rows.Close()

	var foodHistory [][4]interface{}
	for rows.Next() {
		var foodID int
		var date string
		var quantity float64
		var entryID int
		if err := rows.Scan(&foodID, &date, &quantity, &entryID); err != nil {
			return nil, fmt.Errorf("failed to scan food history: %w", err)
		}
		foodHistory = append(foodHistory, [4]interface{}{foodID, date, quantity, entryID})
	}

	return foodHistory, nil
}

func GetWeightHistory(db *sql.DB, userID int) ([][3]interface{}, error) {
	rows, err := db.Query(`
		SELECT date, weight
		FROM weight_history
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get weight history: %w", err)
	}
	defer rows.Close()

	var weightHistory [][3]interface{}
	for rows.Next() {
		var date string
		var weight int
		if err := rows.Scan(&date, &weight); err != nil {
			return nil, fmt.Errorf("failed to scan weight history: %w", err)
		}
		weightHistory = append(weightHistory, [3]interface{}{date, weight})
	}

	return weightHistory, nil
}

func GetBodyFatHistory(db *sql.DB, userID int) ([][3]interface{}, error) {
	rows, err := db.Query(`
		SELECT date, body_fat
		FROM body_fat_history
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get body fat history: %w", err)
	}
	defer rows.Close()

	var bodyFatHistory [][3]interface{}
	for rows.Next() {
		var date string
		var bodyFat float64
		if err := rows.Scan(&date, &bodyFat); err != nil {
			return nil, fmt.Errorf("failed to scan body fat history: %w", err)
		}
		bodyFatHistory = append(bodyFatHistory, [3]interface{}{date, bodyFat})
	}

	return bodyFatHistory, nil
}

func GetIMCHistory(db *sql.DB, userID int) ([][3]interface{}, error) {
	rows, err := db.Query(`
		SELECT date, imc
		FROM imc_history
		WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get IMC history: %w", err)
	}
	defer rows.Close()

	var imcHistory [][3]interface{}
	for rows.Next() {
		var date string
		var imc float64
		if err := rows.Scan(&date, &imc); err != nil {
			return nil, fmt.Errorf("failed to scan IMC history: %w", err)
		}
		imcHistory = append(imcHistory, [3]interface{}{date, imc})
	}

	return imcHistory, nil
}

func DeleteFoodHistory(db *sql.DB, entryID int) error {
	_, err := db.Exec(`
		DELETE FROM food_history
		WHERE id = $1
	`, entryID)
	if err != nil {
		return fmt.Errorf("failed to delete food history: %w", err)
	}
	return nil
}

func UpdateUserFirstname(db *sql.DB, userID int, firstname string) error {
	_, err := db.Exec(`
		UPDATE users
		SET firstname = $1
		WHERE id = $2
	`, firstname, userID)
	if err != nil {
		return fmt.Errorf("failed to update user firstname: %w", err)
	}
	return nil
}
func UpdateUserLastname(db *sql.DB, userID int, lastname string) error {
	_, err := db.Exec(`
		UPDATE users
		SET lastname = $1
		WHERE id = $2
	`, lastname, userID)
	if err != nil {
		return fmt.Errorf("failed to update user lastname: %w", err)
	}
	return nil
}
func UpdateUserAge(db *sql.DB, userID int, age int) error {
	_, err := db.Exec(`
		UPDATE users
		SET age = $1
		WHERE id = $2
	`, age, userID)
	if err != nil {
		return fmt.Errorf("failed to update user age: %w", err)
	}
	return nil
}
func UpdateUserWeight(db *sql.DB, userID int, weight int) error {
	_, err := db.Exec(`
		UPDATE users
		SET weight = $1
		WHERE id = $2
	`, weight, userID)
	if err != nil {
		return fmt.Errorf("failed to update user weight: %w", err)
	}
	return nil
}
func UpdateUserHeight(db *sql.DB, userID int, height int) error {
	_, err := db.Exec(`
		UPDATE users
		SET height = $1
		WHERE id = $2
	`, height, userID)
	if err != nil {
		return fmt.Errorf("failed to update user height: %w", err)
	}
	return nil
}
func UpdateUserTargetWeight(db *sql.DB, userID int, targetWeight int) error {
	_, err := db.Exec(`
		UPDATE users
		SET target_weight = $1
		WHERE id = $2
	`, targetWeight, userID)
	if err != nil {
		return fmt.Errorf("failed to update user target weight: %w", err)
	}
	return nil
}