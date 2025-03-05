package database

import (
	"OrderService/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB = Connect(config.ServiceConfig.DATABASE_URL)

func Connect(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Connect error:", dsn)
	}
	return db
}
