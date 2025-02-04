package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

func ConnectDb() {
	dsn := "host=localhost user=postgres password=123 dbname=go_findest_rest_api_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Transaction{})
	db.AutoMigrate(&User{})

	Database = db
}
