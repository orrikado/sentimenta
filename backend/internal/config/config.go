package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	POSTGRES_HOST     string
	POSTGRES_PORT     string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_DB       string

	JWT_COOKIE_NAME string
	JWT_SECRET      string

	GOOGLE_CLIENT_ID       string
	GOOGLE_CLIENT_SECRET   string
	GOOGLE_CLIENT_CALLBACK string
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("не удалось загрузить .env: %v\n", err)
	}

	return &Config{
		POSTGRES_HOST:     os.Getenv("POSTGRES_HOST"),
		POSTGRES_PORT:     os.Getenv("POSTGRES_PORT"),
		POSTGRES_USER:     os.Getenv("POSTGRES_USER"),
		POSTGRES_PASSWORD: os.Getenv("POSTGRES_PASSWORD"),
		POSTGRES_DB:       os.Getenv("POSTGRES_DB"),

		JWT_COOKIE_NAME: os.Getenv("JWT_COOKIE_NAME"),
		JWT_SECRET:      os.Getenv("JWT_SECRET"),

		GOOGLE_CLIENT_ID:       os.Getenv("GOOGLE_CLIENT_ID"),
		GOOGLE_CLIENT_SECRET:   os.Getenv("GOOGLE_CLIENT_SECRET"),
		GOOGLE_CLIENT_CALLBACK: os.Getenv("GOOGLE_CLIENT_CALLBACK"),
	}
}

//  = 5432
//  = "access_token"
