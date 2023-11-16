package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type dbConfig struct {
	Host     string
	Port     string
	User     string
	Pass     string
	Database string
}

type awsConfig struct {
	AcessKey  string
	SecretKey string
	Region    string
}

func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func getDB() dbConfig {
	cfg := dbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Pass:     os.Getenv("DB_PASS"),
		Database: os.Getenv("DB_NAME"),
	}

	return cfg
}

func getAWS() awsConfig {
	cfg := awsConfig{
		AcessKey:  os.Getenv("AWS_ACCESS_KEY"),
		SecretKey: os.Getenv("AWS_SECRET_KEY"),
		Region:    os.Getenv("AWS_REGION"),
	}

	return cfg
}
