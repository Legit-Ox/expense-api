package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"expense-api/database"
	"expense-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Migrate tables
	err = db.AutoMigrate(&models.Category{}, &models.Transaction{})
	assert.NoError(t, err)

	// Seed test categories
	testCategories := []models.Category{
		{Name: "Food", Type: "expense"},
		{Name: "Salary", Type: "income"},
	}

	for _, category := range testCategories {
		err = db.Create(&category).Error
		assert.NoError(t, err)
	}

	// Set the test database
	database.DB = db
	return db
}

func TestCreateTransaction(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	app := fiber.New()
	app.Post("/transactions", CreateTransaction)

	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		checkResponse  bool
	}{
		{
			name: "Valid expense transaction",
			payload: map[string]interface{}{
				"amount":       50.0,
				"type":         "expense",
				"category_id":  1,
				"description":  "Lunch",
				"date":         time.Now().Format(time.RFC3339),
			},
			expectedStatus: 201,
			checkResponse:  true,
		},
		{
			name: "Valid income transaction",
			payload: map[string]interface{}{
				"amount":       1000.0,
				"type":         "income",
				"category_id":  2,
				"description":  "Monthly salary",
				"date":         time.Now().Format(time.RFC3339),
			},
			expectedStatus: 201,
			checkResponse:  true,
		},
		{
			name: "Invalid transaction type",
			payload: map[string]interface{}{
				"amount":       50.0,
				"type":         "invalid",
				"category_id":  1,
				"description":  "Test",
				"date":         time.Now().Format(time.RFC3339),
			},
			expectedStatus: 400,
			checkResponse:  false,
		},
		{
			name: "Invalid category ID",
			payload: map[string]interface{}{
				"amount":       50.0,
				"type":         "expense",
				"category_id":  999,
				"description":  "Test",
				"date":         time.Now().Format(time.RFC3339),
			},
			expectedStatus: 400,
			checkResponse:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/transactions", bytes.NewReader(payloadBytes))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.checkResponse {
				var response models.Transaction
				err = json.NewDecoder(resp.Body).Decode(&response)
				assert.NoError(t, err)
				assert.NotZero(t, response.ID)
				assert.Equal(t, tt.payload["amount"], response.Amount)
				assert.Equal(t, tt.payload["type"], response.Type)
				assert.Equal(t, tt.payload["description"], response.Description)
			}
		})
	}
}

func TestGetTransactions(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// Create test transactions
	testTransactions := []models.Transaction{
		{
			Amount:      50.0,
			Type:        "expense",
			CategoryID:  1,
			Description: "Lunch",
			Date:        models.FlexibleDate{Time: time.Now()},
		},
		{
			Amount:      1000.0,
			Type:        "income",
			CategoryID:  2,
			Description: "Salary",
			Date:        models.FlexibleDate{Time: time.Now()},
		},
	}

	for _, transaction := range testTransactions {
		err := db.Create(&transaction).Error
		assert.NoError(t, err)
	}

	app := fiber.New()
	app.Get("/transactions", GetTransactions)

	tests := []struct {
		name           string
		query          string
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "Get all transactions",
			query:          "",
			expectedStatus: 200,
			expectedCount:  2,
		},
		{
			name:           "Filter by expense type",
			query:          "?type=expense",
			expectedStatus: 200,
			expectedCount:  1,
		},
		{
			name:           "Filter by income type",
			query:          "?type=income",
			expectedStatus: 200,
			expectedCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/transactions"+tt.query, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var response []models.TransactionResponse
			err = json.NewDecoder(resp.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Len(t, response, tt.expectedCount)
		})
	}
} 