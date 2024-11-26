package database

import (
	"fmt"
	"os"
	"sync"

	"github.com/yuta_2710/go-clean-arc-reviews/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresDatabase struct {
	Db *gorm.DB
}

var (
	once       sync.Once
	dbInstance *postgresDatabase
)

func NewPostgresDatabase(conf *config.Config) Database {
	once.Do(func() {
		host := os.Getenv("POSTGRES_HOST")
		user := os.Getenv("POSTGRES_USER")
		password := os.Getenv("POSTGRES_PASSWORD")
		dbname := os.Getenv("POSTGRES_DB")
		port := os.Getenv("POSTGRES_PORT")
		sslMode := os.Getenv("POSTGRES_SSL_MODE")
		timezone := os.Getenv("POSTGRES_TIMEZONE")

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			host,
			user,
			password,
			dbname,
			port,
			sslMode,
			timezone,
		)
		fmt.Println("Concac")
		fmt.Println(dsn)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

		if err != nil {
			panic("Failed to connect postgres database")
		} else {
			fmt.Println("Successfully connected postgres database")
		}

		dbInstance = &postgresDatabase{db}
	})

	return dbInstance
}

func (p *postgresDatabase) GetDb() *gorm.DB {
	return dbInstance.Db
}
