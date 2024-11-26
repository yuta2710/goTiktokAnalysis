package main

import (
	"fmt"
	"log"

	"github.com/subosito/gotenv"
	"github.com/yuta_2710/go-clean-arc-reviews/config"
	"github.com/yuta_2710/go-clean-arc-reviews/database"
	"github.com/yuta_2710/go-clean-arc-reviews/modules/users/entities"
	"github.com/yuta_2710/go-clean-arc-reviews/shared"
)

func init() {
	gotenv.Load()
}

func main() {
	conf := config.GetConfig()
	db := database.NewPostgresDatabase(conf)
	didTriggerMigration := TriggerUserMigrate(db)
	if didTriggerMigration {
		fmt.Println("Migrating successfully")
	}
}

func TriggerUserMigrate(db database.Database) bool {
	db.GetDb().Migrator().CreateTable(&entities.User{})
	batches := []entities.User{
		{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@gmail.com",
			Password:  "password123",
			Role:      "admin",
			IsActive:  false,
			IsBlocked: false,
			IsAdmin:   true,
		},
		{
			FirstName: "Jane",
			LastName:  "Smith",
			Email:     "jane.smith@gmail.com",
			Password:  "securepass456",
			Role:      "user",
			IsActive:  false,
			IsBlocked: false,
			IsAdmin:   false,
		},
		{
			FirstName: "Alice",
			LastName:  "Johnson",
			Email:     "alice.johnson@gmail.com",
			Password:  "alice2023",
			Role:      "user",
			IsActive:  false,
			IsBlocked: false,
			IsAdmin:   false,
		},
		{
			FirstName: "Bob",
			LastName:  "Brown",
			Email:     "bob.brown@gmail.com",
			Password:  "bobby789",
			Role:      "user",
			IsActive:  false,
			IsBlocked: false,
			IsAdmin:   false,
		},
		{
			FirstName: "Emma",
			LastName:  "Williams",
			Email:     "emma.williams@gmail.com",
			Password:  "emma_pw2023",
			Role:      "admin",
			IsActive:  false,
			IsBlocked: false,
			IsAdmin:   true,
		},
	}

	for i := range batches {
		hashVal, err := shared.HashPassword(batches[i].Password)

		if err != nil {
			log.Printf("Failed to hash password for user %s %s: %v", batches[i].FirstName, batches[i].LastName, err)
		}

		batches[i].Password = hashVal
	}

	db.GetDb().CreateInBatches(batches, 5)

	return true
}
