package main

import (
	"github.com/gin-gonic/gin"
	"go-findest-rest-api/controllers/transactioncontroller"
	"go-findest-rest-api/models"
	"go-findest-rest-api/seeders"
)

func main() {
	r := gin.Default()

	// connect database
	models.ConnectDb()

	// seed user into database
	seeders.SeedUsers(models.Database)

	// routes
	r.POST("/api/transaction", transactioncontroller.CreateTransaction)

	r.Run()
}
