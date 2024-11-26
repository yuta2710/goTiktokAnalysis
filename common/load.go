package common

import (
	"fmt"

	"github.com/yuta_2710/go-clean-arc-reviews/database"
	TodoSchema "github.com/yuta_2710/go-clean-arc-reviews/modules/todo/entities"
	TokenProviderSchema "github.com/yuta_2710/go-clean-arc-reviews/modules/token/entities"
	UserSchema "github.com/yuta_2710/go-clean-arc-reviews/modules/users/entities"
)

// func LoadRelations(repo database.Database) {
// 	fmt.Println("Load called2")
// 	if !repo.GetDb().Migrator().HasTable("users") {
// 		repo.GetDb().Migrator().CreateTable(&UserSchema.User{})
// 	}
// 	if !repo.GetDb().Migrator().HasTable("token_providers") {
// 		repo.GetDb().Migrator().CreateTable(&TokenProviderSchema.TokenProvider{})
// 	}

// 	if !repo.GetDb().Migrator().HasTable("todos") {
// 		repo.GetDb().Migrator().CreateTable(&TodoSchema.Todo{})
// 	}

// 	// if !repo.GetDb().Migrator().HasTable("todo_members") {
// 	// 	repo.GetDb().Migrator().CreateTable(&TodoSchema.TodoMember{})
// 	// }
// }

func LoadRelations(repo database.Database) {
	if repo == nil || repo.GetDb() == nil {
		panic("Database connection is nil, cannot load relations")
	}

	fmt.Println("Loading database schema...")

	// Debug the models being migrated
	fmt.Println("Migrating the following models:")
	fmt.Printf("User: %+v\n", UserSchema.User{})
	fmt.Printf("TokenProvider: %+v\n", TokenProviderSchema.TokenProvider{})
	fmt.Printf("Todo: %+v\n", TodoSchema.Todo{})
	fmt.Printf("TodoMember: %+v\n", TodoSchema.TodoMember{})

	// Migrate all models
	err := repo.GetDb().AutoMigrate(
		UserSchema.User{},
		TokenProviderSchema.TokenProvider{},
		TodoSchema.Todo{},
		TodoSchema.TodoMember{},
	)
	if err != nil {
		panic(fmt.Sprintf("Error migrating schema: %v", err))
	}

	fmt.Println("Database schema loaded successfully!")
}
