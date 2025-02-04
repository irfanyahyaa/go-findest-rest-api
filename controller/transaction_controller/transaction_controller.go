package transactioncontroller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-findest-rest-api/dto"
	"go-findest-rest-api/model"
	"go-findest-rest-api/repository"
	"go-findest-rest-api/util"
	"gorm.io/gorm"
)

type TransactionController struct {
	TransactionRepo repository.DatabaseRepository[model.Transaction]
	UserRepo        repository.DatabaseRepository[model.User]
}

func NewTransactionController(
	transactionRepo repository.DatabaseRepository[model.Transaction],
	userRepo repository.DatabaseRepository[model.User],
) *TransactionController {
	return &TransactionController{
		TransactionRepo: transactionRepo,
		UserRepo:        userRepo,
	}
}

func (tc *TransactionController) CreateTransaction(c *gin.Context) {
	// bind payload into json
	var payload dto.TransactionDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		util.InternalServerError(c, err.Error(), nil)
		return
	}

	// check if user exist
	_, firstErr := tc.UserRepo.First(payload.UserID)
	if firstErr != nil {
		if errors.Is(firstErr, gorm.ErrRecordNotFound) {
			util.NotFound(c, "User not found", nil)
			return
		}

		util.InternalServerError(c, firstErr.Error(), nil)
		return
	}

	// insert transaction into database
	transaction, createErr := tc.TransactionRepo.Create(
		&model.Transaction{
			UserID: payload.UserID,
			Amount: payload.Amount,
			Status: payload.Status,
		},
	)
	if createErr != nil {
		util.InternalServerError(c, createErr.Error(), nil)
		return
	}

	// build response
	res := dto.TransactionResponse{
		ID:        transaction.ID,
		UserID:    transaction.UserID,
		Amount:    transaction.Amount,
		Status:    transaction.Status,
		CreatedAt: transaction.CreatedAt,
	}

	// return response
	util.Created(c, "transaction created successfully", res)
}
