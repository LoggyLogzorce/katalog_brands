package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var database *gorm.DB

var connect = "host=localhost port=5432 user=postgres password=1234 dbname=catalog_of_cosmetic_brands sslmode=disable"

func Connect() *gorm.DB {
	db, err := gorm.Open(postgres.Open(connect), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(err)
	}

	log.Println("Connected to the database")
	return db
}
