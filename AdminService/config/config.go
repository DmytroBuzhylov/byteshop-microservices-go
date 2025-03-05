package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	ADMIN_SERVICE_URL   string
	PRODUCT_SERVICE_URL string
	AUTH_SERVICE_URL    string
}

var ServiceConfig = LoadConfig()

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Файл .env не найден:", err)
	}

	return &Config{
		ADMIN_SERVICE_URL:   os.Getenv("ADMIN_SERVICE_URL"),
		PRODUCT_SERVICE_URL: os.Getenv("PRODUCT_SERVICE_URL"),
		AUTH_SERVICE_URL:    os.Getenv("AUTH_SERVICE_URL"),
	}
}
