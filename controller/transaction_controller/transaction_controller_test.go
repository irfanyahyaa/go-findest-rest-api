package transactioncontroller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/mock"
	"go-findest-rest-api/controller/transaction_controller"
	mocks "go-findest-rest-api/mock"
	"go-findest-rest-api/model"
	"gorm.io/gorm"
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

func TestCreateTransaction(t *testing.T) {
	testCases := map[string]struct {
		mockBody       any
		mockFirstErr   []any
		mockCreateErr  []any
		expectedStatus int
	}{
		"successfully created transaction": {
			mockBody: &model.Transaction{
				UserID: 1,
				Amount: 1,
				Status: "pending",
			},
			mockFirstErr: []any{&model.User{ID: 1}, nil},
			mockCreateErr: []any{&model.Transaction{
				UserID: 1,
				Amount: 1,
				Status: "pending",
			}, nil},
			expectedStatus: http.StatusCreated,
		},
		"error cannot bind payload into json": {
			mockBody:       "wrong-format",
			expectedStatus: http.StatusInternalServerError,
		},
		"error user not found": {
			mockBody: &model.Transaction{
				UserID: 1,
				Amount: 1,
				Status: "pending",
			},
			mockFirstErr:   []any{(*model.User)(nil), gorm.ErrRecordNotFound},
			expectedStatus: http.StatusNotFound,
		},
		"error user internal server error": {
			mockBody: &model.Transaction{
				UserID: 1,
				Amount: 1,
				Status: "pending",
			},
			mockFirstErr:   []any{(*model.User)(nil), errors.New("")},
			expectedStatus: http.StatusInternalServerError,
		},
		"error cannot insert transaction into database": {
			mockBody: &model.Transaction{
				UserID: 1,
				Amount: 1,
				Status: "pending",
			},
			mockFirstErr:   []any{&model.User{ID: 1}, nil},
			mockCreateErr:  []any{(*model.Transaction)(nil), errors.New("")},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			mockTransactionRepo := new(mocks.MockDatabaseRepository[model.Transaction])
			mockUserRepo := new(mocks.MockDatabaseRepository[model.User])

			controller := transactioncontroller.NewTransactionController(
				mockTransactionRepo,
				mockUserRepo,
			)

			mockUserRepo.On("First", mock.Anything, mock.Anything).Return(test.mockFirstErr...).Once()
			mockTransactionRepo.On("Create", mock.Anything).Return(test.mockCreateErr...).Once()

			router := setUpRouter()
			router.POST("/api/transaction", controller.CreateTransaction)

			w := httptest.NewRecorder()

			body, _ := json.Marshal(test.mockBody)
			req, _ := http.NewRequest(http.MethodPost, "/api/transaction", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatus, w.Code)
		})
	}
}

func TestGetTransactions(t *testing.T) {
	testCases := map[string]struct {
		testURL        string
		mockFindErr    []any
		expectedStatus int
	}{
		"successfully get transaction with query": {
			testURL: "/api/transactions?userId=1&status=pending",
			mockFindErr: []any{[]model.Transaction{
				{ID: 1,
					UserID:    1,
					Amount:    1,
					Status:    "pending",
					CreatedAt: time.Now(),
				},
			}, nil},
			expectedStatus: http.StatusOK,
		},
		"successfully get transaction without query": {
			testURL: "/api/transactions",
			mockFindErr: []any{[]model.Transaction{
				{ID: 1,
					UserID:    1,
					Amount:    1,
					Status:    "pending",
					CreatedAt: time.Now(),
				},
				{ID: 2,
					UserID:    2,
					Amount:    2,
					Status:    "success",
					CreatedAt: time.Now(),
				},
			}, nil},
			expectedStatus: http.StatusOK,
		},
		"error get transactions": {
			testURL:        "/api/transactions",
			mockFindErr:    []any{nil, errors.New("")},
			expectedStatus: http.StatusNotFound,
		},
		"error cannot bind payload into json": {
			testURL:        "/api/transactions?userId=wrong-format",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			mockTransactionRepo := new(mocks.MockDatabaseRepository[model.Transaction])
			mockUserRepo := new(mocks.MockDatabaseRepository[model.User])

			controller := transactioncontroller.NewTransactionController(
				mockTransactionRepo,
				mockUserRepo,
			)

			mockTransactionRepo.On("Find", mock.Anything).Return(test.mockFindErr...).Once()

			router := setUpRouter()
			router.GET("/api/transactions", controller.GetTransactions)

			w := httptest.NewRecorder()

			query, _ := json.Marshal(test.testURL)
			req, _ := http.NewRequest(http.MethodGet, test.testURL, bytes.NewBuffer(query))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatus, w.Code)
		})
	}
}
