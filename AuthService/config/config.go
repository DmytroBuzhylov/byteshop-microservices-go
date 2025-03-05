package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	AUTH_SERVICE_URL string
	DBUrl            string
	JWTSecret        string
}

var ServiceConfig = LoadConfig()

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Файл .env не найден")
	}

	return &Config{
		AUTH_SERVICE_URL: os.Getenv("AUTH_SERVICE_URL"),
		DBUrl:            os.Getenv("DATABASE_URL"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
	}
}
