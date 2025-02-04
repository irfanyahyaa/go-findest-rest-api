package database

import (
	"fmt"
	"go-findest-rest-api/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

func InitDb(host string, dbUser string, dbPassword string, dbName string, dbPort string) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, dbUser, dbPassword, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&model.Transaction{})
	db.AutoMigrate(&model.User{})

	Database = db
}
