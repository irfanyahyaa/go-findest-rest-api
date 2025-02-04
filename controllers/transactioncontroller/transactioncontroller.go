package transactioncontroller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-findest-rest-api/dtos"
	"go-findest-rest-api/models"
	"go-findest-rest-api/repository"
	"go-findest-rest-api/utils"
	"gorm.io/gorm"
)

type TransactionController struct {
	TransactionRepo repository.DatabaseRepository[models.Transaction]
	UserRepo        repository.DatabaseRepository[models.User]
}

func NewTransactionController(
	transactionRepo repository.DatabaseRepository[models.Transaction],
	userRepo repository.DatabaseRepository[models.User],
) *TransactionController {
	return &TransactionController{
		TransactionRepo: transactionRepo,
		UserRepo:        userRepo,
	}
}

func (tc *TransactionController) CreateTransaction(c *gin.Context) {
	// bind payload into json
	var payload dtos.TransactionDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.InternalServerError(c, err.Error(), nil)
		return
	}

	// check if user exist
	_, firstErr := tc.UserRepo.First(payload.UserID)
	if firstErr != nil {
		if errors.Is(firstErr, gorm.ErrRecordNotFound) {
			utils.NotFound(c, "User not found", nil)
			return
		}

		utils.InternalServerError(c, firstErr.Error(), nil)
		return
	}

	// insert transaction into database
	transaction, createErr := tc.TransactionRepo.Create(
		&models.Transaction{
			UserID: payload.UserID,
			Amount: payload.Amount,
			Status: payload.Status,
		},
	)
	if createErr != nil {
		utils.InternalServerError(c, createErr.Error(), nil)
		return
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
