package main

import (
	"github.com/gin-gonic/gin"
	"go-findest-rest-api/controllers/transactioncontroller"
	"go-findest-rest-api/database"
	"go-findest-rest-api/models"
	"go-findest-rest-api/repository"
	"go-findest-rest-api/seeders"
)

func main() {
	r := gin.Default()

	// Initialize database connection
	database.InitDb()
	db := database.Database

	// seed user into database
	seeders.SeedUsers(database.Database)

	// Create repositories
	transactionRepo := repository.NewDatabaseRepository[models.Transaction](db)
	userRepo := repository.NewDatabaseRepository[models.User](db)

	// Inject repositories into the controller
	transactionController := transactioncontroller.NewTransactionController(transactionRepo, userRepo)

	// routes
	r.POST("/api/transaction", transactionController.CreateTransaction)

	r.Run()
}
