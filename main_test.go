package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/polarisjrex0406/federico-app/database"
	"github.com/polarisjrex0406/federico-app/dto"
	"github.com/polarisjrex0406/federico-app/routes"
	"github.com/polarisjrex0406/federico-app/utils"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Connect to test database
	database.Connect()

	// Setup dependency injections
	userHandler := setupDependencyInjections(database.DB)

	// Setup router
	r := routes.SetupRouter(
		gin.Default(),
		userHandler,
	)

	return r
}

func TestTransactionIntegration(t *testing.T) {
	// Setup test router
	router := setupTestRouter()

	// Test cases
	testCases := []struct {
		name               string
		userId             uint
		transactionRequest dto.UserDoTransactionRequest
		expectedStatus     int
		expectedCode       string
	}{
		{
			name:   "Successful Transaction",
			userId: 1,
			transactionRequest: dto.UserDoTransactionRequest{
				Amount:        "100.00",
				TransactionID: "Test Transaction",
				State:         "win",
			},
			expectedStatus: http.StatusOK,
			expectedCode:   dto.CODE_SUCCESS,
		},
		{
			name:   "Invalid Transaction Amount",
			userId: 1,
			transactionRequest: dto.UserDoTransactionRequest{
				Amount:        "50.00",
				TransactionID: "Invalid Transaction",
				State:         "loss",
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   dto.CODE_FAILED_REQUEST_BODY_NOT_VALID,
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
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Source-type", "test")

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
		})
	}
}

func TestGetBalanceIntegration(t *testing.T) {
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
