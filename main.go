package main

import (
	"fmt"
	"log"

	"github.com/subosito/gotenv"
	"github.com/yuta_2710/go-clean-arc-reviews/config"
	"github.com/yuta_2710/go-clean-arc-reviews/database"
	"github.com/yuta_2710/go-clean-arc-reviews/server"
	"github.com/yuta_2710/go-clean-arc-reviews/shared"
)

// "github.com/yuta_2710/go-clean-arc-reviews/config"
// "github.com/yuta_2710/go-clean-arc-reviews/database"
// "github.com/yuta_2710/go-clean-arc-reviews/server"

func init() {
	gotenv.Load()
}

func main() {
	fmt.Println("Hello Docker, we're gotik")
	conf := config.GetConfig()
	postgres := database.NewPostgresDatabase(conf)

	if postgres == nil || postgres.GetDb() == nil {
		log.Fatal("Failed to initialize the PostgreSQL database")
	} else {
		fmt.Println("Success create db")
		authIdProvider := &shared.Base64AuthIdProvider{}
		if authIdProvider == nil {
			log.Fatal("AuthIdProvider initialization failed")
		}

		// common.LoadRelations(postgres)
		// common.LoadRelations(postgres)
	}

	// fmt.Println("Eeeee")
	server.NewEchoServer(conf, postgres).Start()
	// fmt.Println("Eeeeeeeeeeeeeeee")
}
