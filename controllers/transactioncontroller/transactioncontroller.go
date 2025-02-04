package transactioncontroller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-findest-rest-api/dtos"
	"go-findest-rest-api/models"
	"go-findest-rest-api/utils"
	"gorm.io/gorm"
)

func CreateTransaction(c *gin.Context) {
	// bind payload into json
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		utils.InternalServerError(c, err.Error(), nil)
		return
	}

	// check if user exist
	var user models.User
	if err := models.Database.First(&user, transaction.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.NotFound(c, "User not found", nil)
			return
		}

		utils.InternalServerError(c, err.Error(), nil)
		return
	}

	// insert transaction into database
	if err := models.Database.Create(&transaction).Error; err != nil {
		utils.InternalServerError(c, err.Error(), nil)
	}

	// build response
	res := dtos.TransactionResponse{
		ID:        transaction.ID,
		UserID:    transaction.UserID,
		Amount:    transaction.Amount,
		Status:    transaction.Status,
		CreatedAt: transaction.CreatedAt,
	}

	// return response
	utils.Created(c, "transaction created successfully", res)
}
