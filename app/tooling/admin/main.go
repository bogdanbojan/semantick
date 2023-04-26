package main

import (
	"fmt"

	"github.com/bogdanbojan/semantick/app/tooling/admin/commands"
	"github.com/bogdanbojan/semantick/business/sys/database"
)

func main() {
	dbConfig := database.Config{
		User:       "postgres",
		Password:   "postgres",
		Host:       "localhost",
		Name:       "postgres",
		DisableTLS: true,
	}

	if err := commands.Migrate(dbConfig); err != nil {
		fmt.Printf("migrating database: %v", err)
	}

	if err := commands.Seed(dbConfig); err != nil {
		fmt.Printf("seeding database: %v", err)
	}
}
