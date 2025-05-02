package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Env Struct to hold environment variables
type Env struct {
	Environment   string
	ServerPort    string
	ServerHost    string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	RedisHost     string
	RedisPort     string
	RedisDB       string
	RedisPassword string
	RedisURL      string
	CORS          string
}

// LoadEnv loads environment variables from a .env file
func LoadEnv() Env {
	// Load .env file
	// load specific to environment using string interpolation
	environment := os.Getenv("ENV")
	if environment == "" {
		// default to development if ENV is not set
		environment = "development"
	}
	// load .env file based on ENV variable
	envFile := fmt.Sprintf(".env.%s", environment)
	err := godotenv.Load(envFile)

	if err != nil {
		// throw error if .env file is not found
		log.Fatalf("Error loading .env.%s file", environment)
		return Env{}
	}

	return Env{
		Environment:   environment,
		ServerPort:    os.Getenv("SERVER_PORT"),
		ServerHost:    os.Getenv("SERVER_HOST"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		DBSSLMode:     os.Getenv("DB_SSL_MODE"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisDB:       os.Getenv("REDIS_DB"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisURL:      os.Getenv("REDIS_URL"),
	}
}

func IsDevelopment() bool {
	return os.Getenv("ENV") == "development"
}

func IsProduction() bool {
	return os.Getenv("ENV") == "production"
}

func IsStaging() bool {
	return os.Getenv("ENV") == "staging"
}
