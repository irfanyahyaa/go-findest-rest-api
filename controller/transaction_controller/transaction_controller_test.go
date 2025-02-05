package transactioncontroller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/mock"
	"go-findest-rest-api/controller/transaction_controller"
	"go-findest-rest-api/dto"
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
			mockBody: &dto.TransactionCreate{
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
			mockBody: &dto.TransactionCreate{
				UserID: 1,
				Amount: 1,
				Status: "pending",
			},
			mockFirstErr:   []any{(*model.User)(nil), gorm.ErrRecordNotFound},
			expectedStatus: http.StatusNotFound,
		},
		"error user internal server error": {
			mockBody: &dto.TransactionCreate{
				UserID: 1,
				Amount: 1,
				Status: "pending",
			},
			mockFirstErr:   []any{(*model.User)(nil), errors.New("")},
			expectedStatus: http.StatusInternalServerError,
		},
		"error cannot insert transaction into database": {
			mockBody: &dto.TransactionCreate{
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

			req, _ := http.NewRequest(http.MethodGet, test.testURL, nil)
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatus, w.Code)
		})
	}
}

func TestGetTransactionById(t *testing.T) {
	testCases := map[string]struct {
		testURL        string
		mockFirstErr   []any
		expectedStatus int
	}{
		"successfully get transaction by id": {
			testURL:        "/api/transaction/1",
			mockFirstErr:   []any{&model.Transaction{ID: 1}, nil},
			expectedStatus: http.StatusOK,
		},
		"error transaction not found": {
			testURL:        "/api/transaction/10",
			mockFirstErr:   []any{nil, gorm.ErrRecordNotFound},
			expectedStatus: http.StatusNotFound,
		},
		"error internal server error": {
			testURL:        "/api/transaction/wrong-format",
			mockFirstErr:   []any{nil, errors.New("")},
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

			mockTransactionRepo.On("First", mock.Anything, mock.Anything).Return(test.mockFirstErr...).Once()

			router := setUpRouter()
			router.GET("/api/transaction/:id", controller.GetTransactionById)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest(http.MethodGet, test.testURL, nil)
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatus, w.Code)
		})
	}
}

func TestUpdateTransaction(t *testing.T) {
	testCases := map[string]struct {
		testURL        string
		mockBody       any
		mockFirstErr   []any
		mockSaveErr    []any
		expectedStatus int
	}{
		"successfully updated transaction status": {
			testURL: "/api/transaction/1",
			mockBody: &dto.TransactionUpdate{
				Status: "success",
			},
			mockFirstErr: []any{&model.Transaction{ID: 1}, nil},
			mockSaveErr: []any{&model.Transaction{
				Status: "success",
			}, nil},
			expectedStatus: http.StatusOK,
		},
		"error cannot bind payload into json": {
			testURL:        "/api/transaction/1",
			mockBody:       "wrong-format",
			expectedStatus: http.StatusInternalServerError,
		},
		"error transaction not found": {
			testURL:        "/api/transaction/1",
			mockBody:       &dto.TransactionUpdate{},
			mockFirstErr:   []any{(*model.Transaction)(nil), gorm.ErrRecordNotFound},
			expectedStatus: http.StatusNotFound,
		},
		"error transaction internal server error": {
			testURL:        "/api/transaction/1",
			mockBody:       &dto.TransactionUpdate{},
			mockFirstErr:   []any{(*model.Transaction)(nil), errors.New("")},
			expectedStatus: http.StatusInternalServerError,
		},
		"error cannot updated transaction status into database": {
			testURL: "/api/transaction/1",
			mockBody: &dto.TransactionUpdate{
				Status: "pending",
			},
			mockFirstErr:   []any{&model.Transaction{ID: 1}, nil},
			mockSaveErr:    []any{(*model.Transaction)(nil), errors.New("")},
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

			mockTransactionRepo.On("First", mock.Anything, mock.Anything).Return(test.mockFirstErr...).Once()
			mockTransactionRepo.On("Save", mock.Anything, mock.Anything).Return(test.mockSaveErr...)

			router := setUpRouter()
			router.PUT("/api/transaction/:id", controller.UpdateTransaction)

			w := httptest.NewRecorder()

			body, _ := json.Marshal(test.mockBody)
			req, _ := http.NewRequest(http.MethodPut, test.testURL, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatus, w.Code)
		})
	}
}
