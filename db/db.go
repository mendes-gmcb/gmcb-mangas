package db

import (
	"fmt"
	"log"
	"trabalho/env"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Load() {
	Connect()
	Migrate()
}

func Connect() (*gorm.DB, error) {
	dsn := buildDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	DB = db

	return db, nil
}

func buildDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		env.DB_CONFIG.Host,
		env.DB_CONFIG.Port,
		env.DB_CONFIG.User,
		env.DB_CONFIG.Pass,
		env.DB_CONFIG.Database,
	)
}
