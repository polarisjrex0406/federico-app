package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/polarisjrex0406/federico-app/database"
	"github.com/polarisjrex0406/federico-app/dto"
	"github.com/polarisjrex0406/federico-app/migrations"
	"github.com/polarisjrex0406/federico-app/routes"
	"github.com/polarisjrex0406/federico-app/utils"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Connect to test database
	database.Connect()

	if err := migrations.Migrate(database.DB); err != nil {
		log.Fatalf("error migration: %v", err)
	}
	if err := migrations.Seeder(database.DB); err != nil {
		log.Fatalf("error migration seeder: %v", err)
	}

	// Setup dependency injections
	userHandler := setupDependencyInjections(database.DB)

	// Setup router
	r := routes.SetupRouter(
		gin.Default(),
		userHandler,
	)

	return r
}

func TestTransaction(t *testing.T) {
	// Setup test router
	router := setupTestRouter()

	// Test cases
	testCases := []struct {
		name               string
		contentType        string
		sourceTypeKey      string
		sourceTypeValue    string
		userId             uint
		transactionRequest dto.UserDoTransactionRequest
		expectedStatus     int
		expectedCode       string
		expectedMessage    string
		expectedResponse   dto.UserDoTransactionResponse
	}{
		{
			name:            "Successful Win Transaction With Source-type game",
			contentType:     "application/json",
			sourceTypeKey:   "Source-type",
			sourceTypeValue: "game",
			userId:          1,
			transactionRequest: dto.UserDoTransactionRequest{
				Amount:        "100.00",
				TransactionID: uuid.New().String(),
				State:         "win",
			},
			expectedStatus:  http.StatusOK,
			expectedCode:    dto.CODE_SUCCESS,
			expectedMessage: dto.MESSAGE_SUCCESS_USER_DO_TRANSACTION,
			expectedResponse: dto.UserDoTransactionResponse{
				UserID:  1,
				Balance: "100.00",
			},
		},
		{
			name:            "Successful Loss Transaction With Source-type server",
			contentType:     "application/json",
			sourceTypeKey:   "Source-type",
			sourceTypeValue: "server",
			userId:          1,
			transactionRequest: dto.UserDoTransactionRequest{
				Amount:        "50.00",
				TransactionID: uuid.New().String(),
				State:         "loss",
			},
			expectedStatus:  http.StatusOK,
			expectedCode:    dto.CODE_SUCCESS,
			expectedMessage: dto.MESSAGE_SUCCESS_USER_DO_TRANSACTION,
			expectedResponse: dto.UserDoTransactionResponse{
				UserID:  1,
				Balance: "50.00",
			},
		},
		{
			name:            "Successful Loss Transaction With Source-type payment",
			contentType:     "application/json",
			sourceTypeKey:   "Source-type",
			sourceTypeValue: "payment",
			userId:          1,
			transactionRequest: dto.UserDoTransactionRequest{
				Amount:        "50.00",
				TransactionID: uuid.New().String(),
				State:         "loss",
			},
			expectedStatus:  http.StatusOK,
			expectedCode:    dto.CODE_SUCCESS,
			expectedMessage: dto.MESSAGE_SUCCESS_USER_DO_TRANSACTION,
			expectedResponse: dto.UserDoTransactionResponse{
				UserID:  1,
				Balance: "0.00",
			},
		},
		{
			name:            "Invalid Content-Type",
			contentType:     "application/xml",
			sourceTypeKey:   "Source-type",
			sourceTypeValue: "payment",
			userId:          1,
			transactionRequest: dto.UserDoTransactionRequest{
				Amount:        "50.00",
				TransactionID: uuid.New().String(),
				State:         "win",
			},
			expectedStatus:  http.StatusUnsupportedMediaType,
			expectedCode:    dto.CODE_FAILED_HEADER_CONTENT_TYPE_WRONG,
			expectedMessage: dto.MESSAGE_FAILED_HEADER_CONTENT_TYPE_WRONG,
		},
		{
			name:            "No Source-type Header",
			contentType:     "application/json",
			sourceTypeKey:   "",
			sourceTypeValue: "",
			userId:          1,
			transactionRequest: dto.UserDoTransactionRequest{
				Amount:        "50.00",
				TransactionID: uuid.New().String(),
				State:         "win",
			},
			expectedStatus:  http.StatusBadRequest,
			expectedCode:    dto.CODE_FAILED_HEADER_SOURCE_TYPE_WRONG,
			expectedMessage: dto.MESSAGE_FAILED_HEADER_SOURCE_TYPE_WRONG,
		},
		{
			name:            "Invalid Source-type value",
			contentType:     "application/json",
			sourceTypeKey:   "Source-type",
			sourceTypeValue: "xxx",
			userId:          1,
			transactionRequest: dto.UserDoTransactionRequest{
				Amount:        "50.00",
				TransactionID: uuid.New().String(),
				State:         "win",
			},
			expectedStatus:  http.StatusBadRequest,
			expectedCode:    dto.CODE_FAILED_HEADER_SOURCE_TYPE_WRONG,
			expectedMessage: dto.MESSAGE_FAILED_HEADER_SOURCE_TYPE_WRONG,
		},
		{
			name:            "Invalid amount format",
			contentType:     "application/json",
			sourceTypeKey:   "Source-type",
			sourceTypeValue: "game",
			userId:          1,
			transactionRequest: dto.UserDoTransactionRequest{
				Amount:        "50.123",
				TransactionID: uuid.New().String(),
				State:         "win",
			},
			expectedStatus:  http.StatusBadRequest,
			expectedCode:    dto.CODE_FAILED_REQUEST_BODY_NOT_VALID,
			expectedMessage: dto.MESSAGE_FAILED_GET_REQUEST_BODY,
		},
		{
			name:            "Invalid state",
			contentType:     "application/json",
			sourceTypeKey:   "Source-type",
			sourceTypeValue: "game",
			userId:          1,
			transactionRequest: dto.UserDoTransactionRequest{
				Amount:        "50.12",
				TransactionID: uuid.New().String(),
				State:         "yyy",
			},
			expectedStatus:  http.StatusBadRequest,
			expectedCode:    dto.CODE_FAILED_REQUEST_BODY_NOT_VALID,
			expectedMessage: dto.MESSAGE_FAILED_GET_REQUEST_BODY,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Prepare request body
			jsonBody, _ := json.Marshal(tc.transactionRequest)
			req, _ := http.NewRequest(
				"POST",
				fmt.Sprintf("/user/%d/transaction", tc.userId),
				bytes.NewBuffer(jsonBody),
			)
			req.Header.Set("Content-Type", tc.contentType)
			if tc.sourceTypeKey != "" {
				req.Header.Set(tc.sourceTypeKey, tc.sourceTypeValue)
			}

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(w, req)

			// Check the response status
			assert.Equal(t, tc.expectedStatus, w.Code)

			// Parse the response
			var response utils.Response
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Check the response code
			assert.Equal(t, tc.expectedCode, response.Code)

			// Check the response message
			assert.Equal(t, tc.expectedMessage, response.Message)

			// Check the response data
			if tc.expectedStatus == http.StatusOK {
				jsonData, err := json.Marshal(response.Data)
				assert.NoError(t, err)
				if err == nil {
					respData := dto.UserDoTransactionResponse{}
					err = json.Unmarshal(jsonData, &respData)
					assert.NoError(t, err)
					if err == nil {
						assert.Equal(t, tc.expectedResponse, respData)
					}
				}
			} else {
				assert.Empty(t, response.Data)
			}
		})
	}
}

func TestGetBalance(t *testing.T) {
	// Setup test router
	router := setupTestRouter()

	// Test cases for getting balance
	testCases := []struct {
		name           string
		userId         uint
		expectedStatus int
		expectedCode   string
	}{
		{
			name:           "Get Balance Successfully",
			userId:         1,
			expectedStatus: http.StatusOK,
			expectedCode:   dto.CODE_SUCCESS,
		},
		{
			name:           "Invalid User ID",
			userId:         9999, // Non-existent user
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   dto.CODE_FAILED_INTERNAL_PROCESS,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request
			req, _ := http.NewRequest(
				"GET",
				fmt.Sprintf("/user/%d/balance", tc.userId),
				nil,
			)

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(w, req)

			// Check the response status
			assert.Equal(t, tc.expectedStatus, w.Code)

			// Parse the response
			var response utils.Response
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Check the response code
			assert.Equal(t, tc.expectedCode, response.Code)

			// Check the response data
			if tc.expectedStatus == http.StatusOK {
				jsonData, err := json.Marshal(response.Data)
				assert.NoError(t, err)
				if err == nil {
					respData := dto.UserGetBalanceResponse{}
					err = json.Unmarshal(jsonData, &respData)
					assert.NoError(t, err)
					if err == nil {
						assert.True(t, regexp.MustCompile(`^\d+(\.\d{1,2})?$`).MatchString(respData.Balance))
					}
				}
			} else {
				assert.Nil(t, response.Data)
			}
		})
	}
}
