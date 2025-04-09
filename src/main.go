package main

import (
	"gotracker/cli"
	"gotracker/utils"
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

	// Initialize the CLI
	go cli.Open()

	for { }
}