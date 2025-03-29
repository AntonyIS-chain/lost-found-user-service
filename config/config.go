package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ENV               string
	SERVER_PORT       string
	USER_TABLE        string
	SECRET_KEY        string
	POSTGRES_DB       string
	POSTGRES_USER     string
	POSTGRES_HOST     string
	POSTGRES_PORT     string
	POSTGRES_PASSWORD string
	DEBUG             bool
	TEST              bool
}

func NewConfig() (*Config, error) {
	ENV := os.Getenv("ENV")
	switch ENV {
	case "development":
		err := godotenv.Load(".env.development")
		if err != nil {
			return nil, err
		}
	}

	var (
		SECRET_KEY        = os.Getenv("SECRET_KEY")
		POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
		POSTGRES_USER     = os.Getenv("POSTGRES_USER")
		POSTGRES_DB       = os.Getenv("POSTGRES_DB")
		POSTGRES_HOST     = os.Getenv("POSTGRES_HOST")
		POSTGRES_PORT     = "5432"
		SERVER_PORT       = "8081"
		USER_TABLE        = "Users"
		DEBUG             = false
		TEST              = false
	)

	switch ENV {
	case "production":
		TEST = false
		DEBUG = false

	case "production_test":
		TEST = true
		DEBUG = true
		USER_TABLE = "ProductionTestUsers"

	case "development":
		TEST = true
		DEBUG = true
		POSTGRES_HOST = "localhost"
		USER_TABLE = "DevUsers"

	case "development_test":
		TEST = true
		DEBUG = true
		SECRET_KEY = "testsecret"
		POSTGRES_PASSWORD = "pass1234"
		POSTGRES_HOST = "localhost"
		USER_TABLE = "TestUsers"

	case "docker":
		TEST = true
		DEBUG = true
		USER_TABLE = "DockerUsers"

	case "docker_test":
		TEST = true
		DEBUG = true
		USER_TABLE = "DockerUsers"
	}

	config := Config{
		ENV:               ENV,
		SERVER_PORT:       SERVER_PORT,
		USER_TABLE:        USER_TABLE,
		SECRET_KEY:        SECRET_KEY,
		DEBUG:             DEBUG,
		TEST:              TEST,
		POSTGRES_DB:       POSTGRES_DB,
		POSTGRES_USER:     POSTGRES_USER,
		POSTGRES_HOST:     POSTGRES_HOST,
		POSTGRES_PORT:     POSTGRES_PORT,
		POSTGRES_PASSWORD: POSTGRES_PASSWORD,
	}

	return &config, nil
}
