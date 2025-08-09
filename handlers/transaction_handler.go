package handlers

import (
	"time"

	"expense-api/database"
	"expense-api/models"

	"github.com/gofiber/fiber/v2"
)

// CreateTransaction handles POST /transactions
func CreateTransaction(c *fiber.Ctx) error {
	var transaction models.Transaction

	if err := c.BodyParser(&transaction); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate transaction type
	if transaction.Type != "expense" && transaction.Type != "income" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Type must be either 'expense' or 'income'",
		})
	}

	// Set default date if not provided
	if transaction.Date.IsZero() {
		transaction.Date = time.Now()
	}

	// Verify category exists and matches type
	var category models.Category
	if err := database.DB.First(&category, transaction.CategoryID).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	if category.Type != transaction.Type {
		return c.Status(400).JSON(fiber.Map{
			"error": "Category type does not match transaction type",
		})
	}

	if err := database.DB.Create(&transaction).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create transaction",
		})
	}

	// Load category for response
	database.DB.Preload("Category").First(&transaction, transaction.ID)

	return c.Status(201).JSON(transaction)
}

// GetTransactions handles GET /transactions
func GetTransactions(c *fiber.Ctx) error {
	var transactions []models.Transaction
	query := database.DB.Preload("Category")

	// Apply type filter if provided
	if transactionType := c.Query("type"); transactionType != "" {
		if transactionType != "expense" && transactionType != "income" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Type must be either 'expense' or 'income'",
			})
		}
		query = query.Where("type = ?", transactionType)
	}

	if err := query.Find(&transactions).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch transactions",
		})
	}

	// Convert to response format
	var response []models.TransactionResponse
	for _, t := range transactions {
		response = append(response, models.TransactionResponse{
			ID:          t.ID,
			Amount:      t.Amount,
			Type:        t.Type,
			CategoryID:  t.CategoryID,
			Category:    t.Category.Name,
			Description: t.Description,
			Date:        t.Date,
			CreatedAt:   t.CreatedAt,
		})
	}

	return c.JSON(response)
}

// GetTransactionsAggregate handles GET /transactions/aggregate
func GetTransactionsAggregate(c *fiber.Ctx) error {
	var transactions []models.Transaction

	if err := database.DB.Preload("Category").Find(&transactions).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch transactions",
		})
	}

	// Calculate aggregates
	categories := make(map[string]float64)
	var totalIncome, totalExpenses float64

	for _, t := range transactions {
		categoryName := t.Category.Name
		categories[categoryName] += t.Amount

		if t.Type == "income" {
			totalIncome += t.Amount
		} else {
			totalExpenses += t.Amount
		}
	}

	response := models.AggregateResponse{
		Categories:    categories,
		TotalIncome:   totalIncome,
		TotalExpenses: totalExpenses,
		NetAmount:     totalIncome - totalExpenses,
	}

	return c.JSON(response)
}

// GetTransaction handles GET /transactions/:id
func GetTransaction(c *fiber.Ctx) error {
	id := c.Params("id")

	var transaction models.Transaction
	if err := database.DB.Preload("Category").First(&transaction, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	response := models.TransactionResponse{
		ID:          transaction.ID,
		Amount:      transaction.Amount,
		Type:        transaction.Type,
		CategoryID:  transaction.CategoryID,
		Category:    transaction.Category.Name,
		Description: transaction.Description,
		Date:        transaction.Date,
		CreatedAt:   transaction.CreatedAt,
	}

	return c.JSON(response)
}

// UpdateTransaction handles PUT /transactions/:id
func UpdateTransaction(c *fiber.Ctx) error {
	id := c.Params("id")

	var transaction models.Transaction
	if err := database.DB.First(&transaction, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	var updateData map[string]interface{}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate transaction type if provided
	if transactionType, exists := updateData["type"]; exists {
		if transactionType != "expense" && transactionType != "income" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Type must be either 'expense' or 'income'",
			})
		}
	}

	// Verify category exists and matches type if category_id is being updated
	if categoryID, exists := updateData["category_id"]; exists {
		var category models.Category
		if err := database.DB.First(&category, categoryID).Error; err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Category not found",
			})
		}

		// Check if type is being updated and matches category type
		transactionType := transaction.Type
		if t, exists := updateData["type"]; exists {
			transactionType = t.(string)
		}

		if category.Type != transactionType {
			return c.Status(400).JSON(fiber.Map{
				"error": "Category type does not match transaction type",
			})
		}
	}

	if err := database.DB.Model(&transaction).Updates(updateData).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update transaction",
		})
	}

	// Load updated transaction with category
	database.DB.Preload("Category").First(&transaction, transaction.ID)

	return c.JSON(transaction)
}

// DeleteTransaction handles DELETE /transactions/:id
func DeleteTransaction(c *fiber.Ctx) error {
	id := c.Params("id")

	var transaction models.Transaction
	if err := database.DB.First(&transaction, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	if err := database.DB.Delete(&transaction).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete transaction",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Transaction deleted successfully",
	})
}

// GetTransactionsByDateRange handles GET /transactions/date-range
func GetTransactionsByDateRange(c *fiber.Ctx) error {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Both start_date and end_date query parameters are required",
		})
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid start_date format. Use YYYY-MM-DD",
		})
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid end_date format. Use YYYY-MM-DD",
		})
	}

	// Set end date to end of day
	endDate = endDate.Add(24*time.Hour - time.Second)

	var transactions []models.Transaction
	query := database.DB.Preload("Category").Where("date BETWEEN ? AND ?", startDate, endDate)

	// Apply type filter if provided
	if transactionType := c.Query("type"); transactionType != "" {
		if transactionType != "expense" && transactionType != "income" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Type must be either 'expense' or 'income'",
			})
		}
		query = query.Where("type = ?", transactionType)
	}

	if err := query.Find(&transactions).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch transactions",
		})
	}

	// Convert to response format
	var response []models.TransactionResponse
	for _, t := range transactions {
		response = append(response, models.TransactionResponse{
			ID:          t.ID,
			Amount:      t.Amount,
			Type:        t.Type,
			CategoryID:  t.CategoryID,
			Category:    t.Category.Name,
			Description: t.Description,
			Date:        t.Date,
			CreatedAt:   t.CreatedAt,
		})
	}

	return c.JSON(response)
}

// GetSummary handles GET /transactions/summary
func GetSummary(c *fiber.Ctx) error {
	// Get total counts
	var totalTransactions int64
	var totalExpenses int64
	var totalIncome int64

	database.DB.Model(&models.Transaction{}).Count(&totalTransactions)
	database.DB.Model(&models.Transaction{}).Where("type = ?", "expense").Count(&totalExpenses)
	database.DB.Model(&models.Transaction{}).Where("type = ?", "income").Count(&totalIncome)

	// Get total amounts
	var expenseSum float64
	var incomeSum float64

	database.DB.Model(&models.Transaction{}).Where("type = ?", "expense").Select("COALESCE(SUM(amount), 0)").Scan(&expenseSum)
	database.DB.Model(&models.Transaction{}).Where("type = ?", "income").Select("COALESCE(SUM(amount), 0)").Scan(&incomeSum)

	// Get recent transactions (last 5)
	var recentTransactions []models.Transaction
	database.DB.Preload("Category").Order("created_at DESC").Limit(5).Find(&recentTransactions)

	// Convert to response format
	var recentResponse []models.TransactionResponse
	for _, t := range recentTransactions {
		recentResponse = append(recentResponse, models.TransactionResponse{
			ID:          t.ID,
			Amount:      t.Amount,
			Type:        t.Type,
			CategoryID:  t.CategoryID,
			Category:    t.Category.Name,
			Description: t.Description,
			Date:        t.Date,
			CreatedAt:   t.CreatedAt,
		})
	}

	summary := fiber.Map{
		"overview": fiber.Map{
			"total_transactions": totalTransactions,
			"total_expenses":     totalExpenses,
			"total_income":       totalIncome,
		},
		"totals": fiber.Map{
			"total_expense_amount": expenseSum,
			"total_income_amount":  incomeSum,
			"net_amount":           incomeSum - expenseSum,
		},
		"recent_transactions": recentResponse,
	}

	return c.JSON(summary)
}
