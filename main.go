package main

import (
	"log"

	"github.com/polarisjrex0406/federico-app/cmd"
	"github.com/polarisjrex0406/federico-app/config"
	"github.com/polarisjrex0406/federico-app/database"
)

func main() {
	_, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	database.Connect()

	if !cmd.Commands(database.DB) {
		log.Fatalf("Error while running commands: %v", err)
	}
}
