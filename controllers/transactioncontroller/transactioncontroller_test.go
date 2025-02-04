package transactioncontroller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/mock"
	"go-findest-rest-api/controllers/transactioncontroller"
	"go-findest-rest-api/mocks"
	"go-findest-rest-api/models"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
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
			mockBody: &models.Transaction{
				UserID: 1,
				Amount: 1,
				Status: "pending",
			},
			mockFirstErr: []any{&models.User{ID: 1}, nil},
			mockCreateErr: []any{&models.Transaction{
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
			mockBody: &models.Transaction{
				UserID: 1,
				Amount: 1,
				Status: "pending",
			},
			mockFirstErr:   []any{(*models.User)(nil), gorm.ErrRecordNotFound},
			expectedStatus: http.StatusNotFound,
		},
		"error user internal server error": {
			mockBody: &models.Transaction{
				UserID: 1,
				Amount: 1,
				Status: "pending",
			},
			mockFirstErr:   []any{(*models.User)(nil), errors.New("")},
			expectedStatus: http.StatusInternalServerError,
		},
		"error cannot insert transaction into database": {
			mockBody: &models.Transaction{
				UserID: 1,
				Amount: 1,
				Status: "pending",
			},
			mockFirstErr:   []any{&models.User{ID: 1}, nil},
			mockCreateErr:  []any{(*models.Transaction)(nil), errors.New("")},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			mockTransactionRepo := new(mocks.MockDatabaseRepository[models.Transaction])
			mockUserRepo := new(mocks.MockDatabaseRepository[models.User])

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
