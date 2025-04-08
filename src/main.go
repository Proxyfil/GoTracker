package main

import (
	"gotracker/cli"
)

func main() {
	go cli.Open()

	for { }
}