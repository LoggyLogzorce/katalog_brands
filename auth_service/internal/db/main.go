package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var database *gorm.DB

// var cfg *configs.Config
var connect = "host=localhost port=5432 user=postgres password=1234 dbname=catalog_of_cosmetic_brands sslmode=disable"

//func init() {
//	cfg = configs.Get()
//	connect = "host=" + cfg.HostDb +
//		" port=" + cfg.PortDb +
//		" user= " + cfg.User +
//		" password=" + cfg.Password +
//		" dbname=" + cfg.DbName +
//		" sslmode=" + cfg.SslMode
//}

func DB() *gorm.DB {
	return database
}

func Connect() {
	db, err := gorm.Open(postgres.Open(connect), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(err)
	}

	database = db
	log.Println("Connected to the database")
}
