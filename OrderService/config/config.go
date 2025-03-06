package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DATABASE_URL      string
	AUTH_SERVICE_URL  string
	ORDER_SERVICE_URL string
	EMAIL             string
}

var ServiceConfig = LoadConfig()

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Файл .env не найден")
	}

	return &Config{
		DATABASE_URL:      os.Getenv("DATABASE_URL"),
		AUTH_SERVICE_URL:  os.Getenv("AUTH_SERVICE_URL"),
		ORDER_SERVICE_URL: os.Getenv("ORDER_SERVICE_URL"),
		EMAIL:             os.Getenv("EMAIL"),
	}
}
