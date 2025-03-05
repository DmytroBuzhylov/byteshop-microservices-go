package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	CHAT_SERVICE_URL    string
	MONGO_URL           string
	AUTH_SERVICE_URL    string
	PRODUCT_SERVICE_URL string
}

var ServiceConfig = loadConfig()

func loadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}

	return &Config{
		CHAT_SERVICE_URL:    os.Getenv("CHAT_SERVICE_URL"),
		MONGO_URL:           os.Getenv("MONGO_URL"),
		AUTH_SERVICE_URL:    os.Getenv("AUTH_SERVICE_URL"),
		PRODUCT_SERVICE_URL: os.Getenv("PRODUCT_SERVICE_URL"),
	}
}
