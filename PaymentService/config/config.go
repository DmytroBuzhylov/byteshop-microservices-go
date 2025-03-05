package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	PAYPAL_CLIENT_ID    string
	PAYPAL_SECRET       string
	PAYMENT_SERVICE_URL string
	PAYPAL_API          string
	DATABASE_URL        string
	AUTH_SERVICE_URL    string
	PRODUCT_SERVICE_URL string
}

var ServiceConfig = LoadConfig()

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Файл .env не найден")
	}

	return &Config{
		PAYPAL_CLIENT_ID:    os.Getenv("PAYPAL_CLIENT_ID"),
		PAYPAL_SECRET:       os.Getenv("PAYPAL_SECRET"),
		PAYMENT_SERVICE_URL: os.Getenv("PAYMENT_SERVICE_URL"),
		PAYPAL_API:          os.Getenv("PAYPAL_API"),
		DATABASE_URL:        os.Getenv("DATABASE_URL"),
		AUTH_SERVICE_URL:    os.Getenv("AUTH_SERVICE_URL"),
		PRODUCT_SERVICE_URL: os.Getenv("PRODUCT_SERVICE_URL"),
	}
}
