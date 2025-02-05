package dashboardcontroller_test

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	dashboardcontroller "go-findest-rest-api/controller/dashboard_controller"
	"go-findest-rest-api/dto"
	mocks "go-findest-rest-api/mock"
	"go-findest-rest-api/model"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setUpRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	return r
}

func TestGetTransactions(t *testing.T) {
	testCases := map[string]struct {
		mockFindSuccessfulErr []any
		mockAvgTransactionErr []any
		mockFindLatestErr     []any
		expectedStatus        int
	}{
		"successfully get dashboard summary": {
			mockFindSuccessfulErr: []any{[]model.Transaction{
				{ID: 1,
					UserID:    1,
					Amount:    1,
					Status:    "pending",
					CreatedAt: time.Now(),
				},
			}, nil},
			mockAvgTransactionErr: []any{[]dto.AverageTransactionAttr{
				{
					UserId:         1,
					AvgTransaction: 1,
				},
			}, nil},
			mockFindLatestErr: []any{[]model.Transaction{
				{ID: 1,
					UserID:    1,
					Amount:    1,
					Status:    "pending",
					CreatedAt: time.Now(),
				},
			}, nil},
			expectedStatus: http.StatusOK,
		},
		"successfully get empty dashboard summary": {
			mockFindSuccessfulErr: []any{[]model.Transaction{}, nil},
			mockAvgTransactionErr: []any{[]dto.AverageTransactionAttr{}, nil},
			mockFindLatestErr:     []any{[]model.Transaction{}, nil},
			expectedStatus:        http.StatusOK,
		},
		"error internal server error": {
			mockFindSuccessfulErr: []any{nil, errors.New("")},
			mockAvgTransactionErr: []any{nil, errors.New("")},
			mockFindLatestErr:     []any{nil, errors.New("")},
			expectedStatus:        http.StatusInternalServerError,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			mockTransactionRepo := new(mocks.MockDatabaseRepository[model.Transaction])
			mockUserRepo := new(mocks.MockDatabaseRepository[model.User])

			controller := dashboardcontroller.NewDashboardController(
				mockTransactionRepo,
				mockUserRepo,
			)

			mockTransactionRepo.On("Find", mock.Anything).Return(test.mockFindSuccessfulErr...).Once()
			mockTransactionRepo.On("AverageTransaction", mock.Anything).Return(test.mockAvgTransactionErr...).Once()
			mockTransactionRepo.On("Find", mock.Anything).Return(test.mockFindLatestErr...).Once()

			router := setUpRouter()
			router.GET("/api/dashboard/summary", controller.GetDashboardSummary)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodGet, "/api/dashboard/summary", nil)
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatus, w.Code)
		})
	}
}
