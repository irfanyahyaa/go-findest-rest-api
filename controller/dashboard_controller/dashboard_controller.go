package dashboardcontroller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-findest-rest-api/dto"
	"go-findest-rest-api/model"
	"go-findest-rest-api/repository"
	"go-findest-rest-api/util"
	"time"
)

type DashboardController struct {
	TransactionRepo repository.DatabaseRepository[model.Transaction]
	UserRepo        repository.DatabaseRepository[model.User]
}

func NewDashboardController(
	transactionRepo repository.DatabaseRepository[model.Transaction],
	userRepo repository.DatabaseRepository[model.User],
) *DashboardController {
	return &DashboardController{
		TransactionRepo: transactionRepo,
		UserRepo:        userRepo,
	}
}

func (dc *DashboardController) GetDashboardSummary(c *gin.Context) {
	// build date filter
	today := time.Now().Format("2006-01-02")
	dateFilter := fmt.Sprintf("AND status = 'success' AND updated_at BETWEEN '%s 00:00:00' AND '%s 23:59:59'", today, today)

	// fetched data
	successfulTransactionsToday, err1 := dc.TransactionRepo.Find(dateFilter)
	averageTransactionPerUser, err2 := dc.TransactionRepo.AverageTransaction()
	latestTransactions, err3 := dc.TransactionRepo.Find("ORDER BY created_at DESC LIMIT 10")

	// handling error
	if err1 != nil || err2 != nil || err3 != nil {
		util.InternalServerError(c, "internal server error", nil)
		return
	}

	// build response
	res := dto.DashboardResponse{
		SuccessfulTransactionToday: buildTransactionPagination(successfulTransactionsToday),
		AverageTransactionPerUser:  averageTransactionPerUser,
		LatestTransaction:          buildTransactionPagination(latestTransactions),
	}

	// return response
	util.Success(c, "dashboard summary fetched successfully", res)
}

func buildTransactionPagination(transactions []model.Transaction) dto.DashboardPagination[dto.TransactionResponse] {
	mapped := make([]dto.TransactionResponse, 0, len(transactions))
	for _, t := range transactions {
		mapped = append(mapped, dto.TransactionResponse{
			ID:        t.ID,
			UserID:    t.UserID,
			Amount:    t.Amount,
			Status:    t.Status,
			CreatedAt: t.CreatedAt,
		})
	}

	return dto.DashboardPagination[dto.TransactionResponse]{
		TotalRecords: len(transactions),
		Transactions: mapped,
	}
}
