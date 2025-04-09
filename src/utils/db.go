package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// ConnectDB establishes a connection to the postgres database
func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "user=postgres dbname=postgres sslmode=disable host=localhost port=5432 password=postgres")
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
			body_fat FLOAT,
			imc FLOAT,
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

	// Create food table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS food (
			id SERIAL PRIMARY KEY,
			name VARCHAR(50),
			calories FLOAT
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create food table: %w", err)
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
			food_id INT REFERENCES food(id),
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
			food_id INT REFERENCES food(id),
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

	return nil
}