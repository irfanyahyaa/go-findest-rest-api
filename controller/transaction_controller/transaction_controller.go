package transactioncontroller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-findest-rest-api/dto"
	"go-findest-rest-api/model"
	"go-findest-rest-api/repository"
	"go-findest-rest-api/util"
	"gorm.io/gorm"
	"strconv"
	"strings"
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
			util.NotFound(c, "user not found", nil)
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

func (tc *TransactionController) GetTransactions(c *gin.Context) {
	// bind payload into json
	var payload dto.GetTransactionsQuery
	if err := c.ShouldBindQuery(&payload); err != nil {
		util.InternalServerError(c, err.Error(), nil)
		return
	}

	// map payload into filters
	var filters = map[string]interface{}{}
	if payload.UserID != 0 {
		filters["user_id"] = payload.UserID
	}
	if payload.Status != "" {
		filters["status"] = payload.Status
	}

	// build filter string
	var conditions []string
	for k, v := range filters {
		var valueStr string
		switch v := v.(type) {
		case string:
			valueStr = fmt.Sprintf("'%s'", v)
		case uint:
			valueStr = strconv.FormatUint(uint64(v), 10)
		}

		conditions = append(conditions, fmt.Sprintf("%s = %s", k, valueStr))
	}

	// join conditions with " OR "
	var filter string
	if conditions != nil && len(conditions) > 0 {
		filter = "WHERE " + strings.Join(conditions, " OR ")
	}

	// find all transactions
	transactions, findErr := tc.TransactionRepo.Find(filter)
	if findErr != nil {
		util.NotFound(c, "transaction not found", []dto.TransactionResponse{})
	}

	// build response
	res := dto.Pagination[dto.TransactionResponse]{
		TotalRecords: len(transactions),
		Data:         []dto.TransactionResponse{},
	}
	if len(transactions) > 0 {
		for _, transaction := range transactions {
			res.Data = append(res.Data, dto.TransactionResponse{
				ID:        transaction.ID,
				UserID:    transaction.UserID,
				Amount:    transaction.Amount,
				Status:    transaction.Status,
				CreatedAt: transaction.CreatedAt,
			})
		}
	}

	// return response
	util.Success(c, "transaction(s) fetched successfully", res)
}

func (tc *TransactionController) GetTransactionById(c *gin.Context) {
	// get param from context
	id := c.Param("id")

	// check if transaction exist
	transaction, firstErr := tc.TransactionRepo.First(id)
	if firstErr != nil {
		if errors.Is(firstErr, gorm.ErrRecordNotFound) {
			util.NotFound(c, "transaction not found", nil)
			return
		}

		util.InternalServerError(c, firstErr.Error(), nil)
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
	util.Success(c, "transaction fetched successfully", res)
}
