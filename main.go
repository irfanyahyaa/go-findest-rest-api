package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go-findest-rest-api/controller/transaction_controller"
	"go-findest-rest-api/database"
	"go-findest-rest-api/model"
	"go-findest-rest-api/repository"
	"go-findest-rest-api/seeder"
	"log"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	r := gin.Default()

	// initialize database connection
	database.InitDb(
		os.Getenv("HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	db := database.Database

	// seed user into database
	seeder.SeedUsers(database.Database)

	// create repositories
	transactionRepo := repository.NewDatabaseRepository[model.Transaction](db)
	userRepo := repository.NewDatabaseRepository[model.User](db)

	// inject repositories into the controller
	transactionController := transactioncontroller.NewTransactionController(transactionRepo, userRepo)

	// routes
	r.POST("/api/transaction", transactionController.CreateTransaction)
	r.GET("/api/transactions", transactionController.GetTransactions)

	r.Run()
}
