package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func Connect(dsn string) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Ошибка получения SQL соединения: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Ошибка пинга БД: %v", err)
	}

	log.Println("Успешное подключение к БД!")
	DB = db

}
