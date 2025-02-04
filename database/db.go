package database

import (
	"go-findest-rest-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

func InitDb() {
	dsn := "host=localhost user=postgres password=123 dbname=go_findest_rest_api_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.Transaction{})
	db.AutoMigrate(&models.User{})

	Database = db
}
