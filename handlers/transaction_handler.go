package handlers

import (
	"time"

	"expense-api/database"
	"expense-api/models"

	"github.com/gofiber/fiber/v2"
)

// convertToTransactionResponse converts a Transaction model to TransactionResponse
func convertToTransactionResponse(t models.Transaction) models.TransactionResponse {
	response := models.TransactionResponse{
		ID:            t.ID,
		TransactionID: t.TransactionID,
		Amount:        t.Amount,
		Type:          t.Type,
		CategoryID:    t.CategoryID,
		BankAccountID: t.BankAccountID,
		BankAccount: models.BankAccountResponse{
			ID:            t.BankAccount.ID,
			Name:          t.BankAccount.Name,
			AccountNumber: t.BankAccount.AccountNumber,
			BankName:      t.BankAccount.BankName,
			AccountType:   t.BankAccount.AccountType,
			Balance:       t.BankAccount.Balance,
			IsActive:      t.BankAccount.IsActive,
		},
		DestinationBankAccountID: t.DestinationBankAccountID,
		Description:              t.Description,
		Date:                     t.Date.Time,
		CreatedAt:                t.CreatedAt,
	}

	// Set category name if category exists
	if t.CategoryID != nil {
		response.Category = t.Category.Name
	}

	// Set destination bank account if it exists
	if t.DestinationBankAccountID != nil {
		response.DestinationBankAccount = &models.BankAccountResponse{
			ID:            t.DestinationBankAccount.ID,
			Name:          t.DestinationBankAccount.Name,
			AccountNumber: t.DestinationBankAccount.AccountNumber,
			BankName:      t.DestinationBankAccount.BankName,
			AccountType:   t.DestinationBankAccount.AccountType,
			Balance:       t.DestinationBankAccount.Balance,
			IsActive:      t.DestinationBankAccount.IsActive,
		}
	}

	return response
}

// CreateTransaction handles POST /transactions
func CreateTransaction(c *fiber.Ctx) error {
	var transaction models.Transaction

	if err := c.BodyParser(&transaction); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate transaction type
	if transaction.Type != "expense" && transaction.Type != "income" && transaction.Type != "transfer" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Type must be either 'expense', 'income', or 'transfer'",
		})
	}

	// Set default date if not provided
	if transaction.Date.IsZero() {
		transaction.Date = models.FlexibleDate{Time: time.Now()}
	}

	// Validate bank account exists
	var bankAccount models.BankAccount
	if err := database.DB.First(&bankAccount, transaction.BankAccountID).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Bank account not found",
		})
	}

	// Validate based on transaction type
	if transaction.Type == "transfer" {
		// For transfers, category is not required but destination account is
		if transaction.DestinationBankAccountID == nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Destination bank account is required for transfers",
			})
		}

		// Validate destination bank account exists
		var destBankAccount models.BankAccount
		if err := database.DB.First(&destBankAccount, *transaction.DestinationBankAccountID).Error; err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Destination bank account not found",
			})
		}

		// Cannot transfer to the same account
		if transaction.BankAccountID == *transaction.DestinationBankAccountID {
			return c.Status(400).JSON(fiber.Map{
				"error": "Cannot transfer to the same bank account",
			})
		}

		// Set category to nil for transfers
		transaction.CategoryID = nil
	} else {
		// For expense/income, category is required
		if transaction.CategoryID == nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Category is required for expense and income transactions",
			})
		}

		// Verify category exists and matches type
		var category models.Category
		if err := database.DB.First(&category, *transaction.CategoryID).Error; err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Category not found",
			})
		}

		if category.Type != transaction.Type {
			return c.Status(400).JSON(fiber.Map{
				"error": "Category type does not match transaction type",
			})
		}

		// Destination account should not be set for expense/income
		transaction.DestinationBankAccountID = nil
	}

	if err := database.DB.Create(&transaction).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create transaction",
		})
	}

	// Load related data for response
	database.DB.Preload("Category").Preload("BankAccount").Preload("DestinationBankAccount").First(&transaction, transaction.ID)

	// Convert to response format
	response := convertToTransactionResponse(transaction)

	return c.Status(201).JSON(response)
}

// GetTransactions handles GET /transactions
func GetTransactions(c *fiber.Ctx) error {
	var transactions []models.Transaction
	query := database.DB.Preload("Category").Preload("BankAccount").Preload("DestinationBankAccount")

	// Apply type filter if provided
	if transactionType := c.Query("type"); transactionType != "" {
		if transactionType != "expense" && transactionType != "income" && transactionType != "transfer" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Type must be either 'expense', 'income', or 'transfer'",
			})
		}
		query = query.Where("type = ?", transactionType)
	}

	// Apply bank account filter if provided
	if bankAccountID := c.Query("bank_account_id"); bankAccountID != "" {
		query = query.Where("bank_account_id = ? OR destination_bank_account_id = ?", bankAccountID, bankAccountID)
	}

	if err := query.Order("date DESC").Find(&transactions).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch transactions",
		})
	}

	// Convert to response format
	var response []models.TransactionResponse
	for _, t := range transactions {
		response = append(response, convertToTransactionResponse(t))
	}

	return c.JSON(response)
}

// GetTransactionsAggregate handles GET /transactions/aggregate
func GetTransactionsAggregate(c *fiber.Ctx) error {
	var transactions []models.Transaction

	// Exclude transfers from aggregation
	if err := database.DB.Preload("Category").Where("type != ?", "transfer").Order("date DESC").Find(&transactions).Error; err != nil {
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
		} else if t.Type == "expense" {
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
	if err := database.DB.Preload("Category").Preload("BankAccount").Preload("DestinationBankAccount").First(&transaction, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	response := convertToTransactionResponse(transaction)
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

	if err := query.Order("date DESC").Find(&transactions).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch transactions",
		})
	}

	// Convert to response format
	var response []models.TransactionResponse
	for _, t := range transactions {
		response = append(response, models.TransactionResponse{
			ID:            t.ID,
			TransactionID: t.TransactionID,
			Amount:        t.Amount,
			Type:          t.Type,
			CategoryID:    t.CategoryID,
			Category:      t.Category.Name,
			Description:   t.Description,
			Date:          t.Date.Time,
			CreatedAt:     t.CreatedAt,
		})
	}

	return c.JSON(response)
}

// CreateBulkTransactions handles POST /transactions/bulk
func CreateBulkTransactions(c *fiber.Ctx) error {
	var request models.BulkTransactionRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if len(request.Transactions) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "At least one transaction is required",
		})
	}

	if len(request.Transactions) > 5000 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Maximum 5000 transactions allowed per bulk request",
		})
	}

	var response models.BulkTransactionResponse
	response.TotalCount = len(request.Transactions)

	// Process each transaction
	for i, transaction := range request.Transactions {
		// Set default date if not provided
		if transaction.Date.IsZero() {
			transaction.Date = models.FlexibleDate{Time: time.Now()}
		}

		// Validate transaction type
		if transaction.Type != "expense" && transaction.Type != "income" {
			response.Failed = append(response.Failed, models.BulkTransactionError{
				Index:       i,
				Transaction: transaction,
				Error:       "Type must be either 'expense' or 'income'",
			})
			continue
		}

		// Verify category exists and matches type
		var category models.Category
		if err := database.DB.First(&category, transaction.CategoryID).Error; err != nil {
			response.Failed = append(response.Failed, models.BulkTransactionError{
				Index:       i,
				Transaction: transaction,
				Error:       "Category not found",
			})
			continue
		}

		if category.Type != transaction.Type {
			response.Failed = append(response.Failed, models.BulkTransactionError{
				Index:       i,
				Transaction: transaction,
				Error:       "Category type does not match transaction type",
			})
			continue
		}

		// Create transaction
		if err := database.DB.Create(&transaction).Error; err != nil {
			response.Failed = append(response.Failed, models.BulkTransactionError{
				Index:       i,
				Transaction: transaction,
				Error:       "Failed to create transaction: " + err.Error(),
			})
			continue
		}

		// Load category for response
		database.DB.Preload("Category").First(&transaction, transaction.ID)

		// Add to success list
		response.Success = append(response.Success, models.TransactionResponse{
			ID:            transaction.ID,
			TransactionID: transaction.TransactionID,
			Amount:        transaction.Amount,
			Type:          transaction.Type,
			CategoryID:    transaction.CategoryID,
			Category:      transaction.Category.Name,
			Description:   transaction.Description,
			Date:          transaction.Date.Time,
			CreatedAt:     transaction.CreatedAt,
		})
	}

	response.SuccessCount = len(response.Success)
	response.FailedCount = len(response.Failed)

	// Return appropriate status code
	statusCode := 201
	if response.FailedCount > 0 {
		if response.SuccessCount == 0 {
			statusCode = 400 // All failed
		} else {
			statusCode = 207 // Partial success (Multi-Status)
		}
	}

	return c.Status(statusCode).JSON(response)
}

// UpdateTransactionCategory handles PATCH /transactions/:id/category
func UpdateTransactionCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	var transaction models.Transaction
	if err := database.DB.First(&transaction, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Transaction not found",
		})
	}

	var request struct {
		CategoryID uint `json:"category_id" validate:"required"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if request.CategoryID == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "category_id is required",
		})
	}

	// Verify new category exists and matches transaction type
	var newCategory models.Category
	if err := database.DB.First(&newCategory, request.CategoryID).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	if newCategory.Type != transaction.Type {
		return c.Status(400).JSON(fiber.Map{
			"error": "Category type does not match transaction type",
		})
	}

	// Update the category
	if err := database.DB.Model(&transaction).Update("category_id", request.CategoryID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update transaction category",
		})
	}

	// Load updated transaction with new category
	database.DB.Preload("Category").First(&transaction, transaction.ID)

	response := models.TransactionResponse{
		ID:            transaction.ID,
		TransactionID: transaction.TransactionID,
		Amount:        transaction.Amount,
		Type:          transaction.Type,
		CategoryID:    transaction.CategoryID,
		Category:      transaction.Category.Name,
		Description:   transaction.Description,
		Date:          transaction.Date.Time,
		CreatedAt:     transaction.CreatedAt,
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
			ID:            t.ID,
			TransactionID: t.TransactionID,
			Amount:        t.Amount,
			Type:          t.Type,
			CategoryID:    t.CategoryID,
			Category:      t.Category.Name,
			Description:   t.Description,
			Date:          t.Date.Time,
			CreatedAt:     t.CreatedAt,
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

// DeleteBulkTransactions handles DELETE /transactions/bulk
func DeleteBulkTransactions(c *fiber.Ctx) error {
	var request models.BulkDeleteRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if len(request.TransactionIDs) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "At least one transaction ID is required",
		})
	}

	if len(request.TransactionIDs) > 1000 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Maximum 1000 transaction IDs allowed per bulk delete request",
		})
	}

	var response models.BulkDeleteResponse
	response.TotalCount = len(request.TransactionIDs)

	// Process each transaction ID
	for _, transactionID := range request.TransactionIDs {
		// Check if transaction exists
		var transaction models.Transaction
		if err := database.DB.First(&transaction, transactionID).Error; err != nil {
			response.Failed = append(response.Failed, models.BulkDeleteError{
				TransactionID: transactionID,
				Error:         "Transaction not found",
			})
			continue
		}

		// Delete the transaction
		if err := database.DB.Delete(&transaction).Error; err != nil {
			response.Failed = append(response.Failed, models.BulkDeleteError{
				TransactionID: transactionID,
				Error:         "Failed to delete transaction: " + err.Error(),
			})
			continue
		}

		// Add to success list
		response.Deleted = append(response.Deleted, transactionID)
	}

	response.DeletedCount = len(response.Deleted)
	response.FailedCount = len(response.Failed)

	// Return appropriate status code
	statusCode := 200
	if response.FailedCount > 0 {
		if response.DeletedCount == 0 {
			statusCode = 400 // All failed
		} else {
			statusCode = 207 // Partial success (Multi-Status)
		}
	}

	return c.Status(statusCode).JSON(response)
}
// GetTransactionsAggregateTable handles GET /transactions/aggregate-table
func GetTransactionsAggregateTable(c *fiber.Ctx) error {
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

	// Query transactions within date range
	var transactions []models.Transaction
	if err := database.DB.Preload("Category").Where("date BETWEEN ? AND ?", startDate, endDate).Find(&transactions).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch transactions",
		})
	}

	// Initialize response
	response := models.AggregateTableResponse{
		DateRange: models.DateRange{
			StartDate: startDateStr,
			EndDate:   endDateStr,
		},
	}

	// Maps to aggregate data by category
	incomeCategories := make(map[uint]*models.CategoryAggregate)
	expenseCategories := make(map[uint]*models.CategoryAggregate)

	var totalIncome, totalExpenses float64
	var incomeTransactionCount, expenseTransactionCount int

	// Process each transaction
	for _, t := range transactions {
		// Skip transfers as they don't have categories
		if t.Type == "transfer" || t.CategoryID == nil {
			continue
		}

		categoryID := *t.CategoryID
		
		if t.Type == "income" {
			totalIncome += t.Amount
			incomeTransactionCount++

			if agg, exists := incomeCategories[categoryID]; exists {
				agg.TotalAmount += t.Amount
				agg.TransactionCount++
			} else {
				incomeCategories[categoryID] = &models.CategoryAggregate{
					CategoryID:       categoryID,
					CategoryName:     t.Category.Name,
					TotalAmount:      t.Amount,
					TransactionCount: 1,
				}
			}
		} else if t.Type == "expense" {
			totalExpenses += t.Amount
			expenseTransactionCount++

			if agg, exists := expenseCategories[categoryID]; exists {
				agg.TotalAmount += t.Amount
				agg.TransactionCount++
			} else {
				expenseCategories[categoryID] = &models.CategoryAggregate{
					CategoryID:       categoryID,
					CategoryName:     t.Category.Name,
					TotalAmount:      t.Amount,
					TransactionCount: 1,
				}
			}
		}
	}

	// Convert maps to slices for JSON response
	for _, agg := range incomeCategories {
		response.Income.Categories = append(response.Income.Categories, *agg)
	}
	for _, agg := range expenseCategories {
		response.Expenses.Categories = append(response.Expenses.Categories, *agg)
	}

	// Set totals
	response.Income.TotalAmount = totalIncome
	response.Income.TotalTransactions = incomeTransactionCount
	response.Expenses.TotalAmount = totalExpenses
	response.Expenses.TotalTransactions = expenseTransactionCount

	// Set summary
	response.Summary.TotalIncome = totalIncome
	response.Summary.TotalExpenses = totalExpenses
	response.Summary.NetAmount = totalIncome - totalExpenses

	return c.JSON(response)
}

// CreateTransfer handles POST /transactions/transfer
func CreateTransfer(c *fiber.Ctx) error {
	var transferRequest models.TransferRequest

	if err := c.BodyParser(&transferRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate amount
	if transferRequest.Amount <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Amount must be greater than 0",
		})
	}

	// Validate source and destination accounts are different
	if transferRequest.BankAccountID == transferRequest.DestinationBankAccountID {
		return c.Status(400).JSON(fiber.Map{
			"error": "Cannot transfer to the same bank account",
		})
	}

	// Validate source bank account exists
	var sourceBankAccount models.BankAccount
	if err := database.DB.First(&sourceBankAccount, transferRequest.BankAccountID).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Source bank account not found",
		})
	}

	// Validate destination bank account exists
	var destBankAccount models.BankAccount
	if err := database.DB.First(&destBankAccount, transferRequest.DestinationBankAccountID).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Destination bank account not found",
		})
	}

	// Set default date if not provided
	if transferRequest.Date.IsZero() {
		transferRequest.Date = models.FlexibleDate{Time: time.Now()}
	}

	// Create transfer transaction
	transaction := models.Transaction{
		TransactionID:           transferRequest.TransactionID,
		Amount:                  transferRequest.Amount,
		Type:                    "transfer",
		CategoryID:              nil, // Transfers don't have categories
		BankAccountID:           transferRequest.BankAccountID,
		DestinationBankAccountID: &transferRequest.DestinationBankAccountID,
		Description:             transferRequest.Description,
		Date:                    transferRequest.Date,
	}

	if err := database.DB.Create(&transaction).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create transfer",
		})
	}

	// Load related data for response
	database.DB.Preload("BankAccount").Preload("DestinationBankAccount").First(&transaction, transaction.ID)

	// Convert to transfer response format
	response := models.TransferResponse{
		ID:            transaction.ID,
		TransactionID: transaction.TransactionID,
		Amount:        transaction.Amount,
		BankAccount: models.BankAccountResponse{
			ID:            transaction.BankAccount.ID,
			Name:          transaction.BankAccount.Name,
			AccountNumber: transaction.BankAccount.AccountNumber,
			BankName:      transaction.BankAccount.BankName,
			AccountType:   transaction.BankAccount.AccountType,
			Balance:       transaction.BankAccount.Balance,
			IsActive:      transaction.BankAccount.IsActive,
		},
		DestinationBankAccount: models.BankAccountResponse{
			ID:            transaction.DestinationBankAccount.ID,
			Name:          transaction.DestinationBankAccount.Name,
			AccountNumber: transaction.DestinationBankAccount.AccountNumber,
			BankName:      transaction.DestinationBankAccount.BankName,
			AccountType:   transaction.DestinationBankAccount.AccountType,
			Balance:       transaction.DestinationBankAccount.Balance,
			IsActive:      transaction.DestinationBankAccount.IsActive,
		},
		Description: transaction.Description,
		Date:        transaction.Date.Time,
		CreatedAt:   transaction.CreatedAt,
	}

	return c.Status(201).JSON(response)
}

// GetTransfers handles GET /transactions/transfers
func GetTransfers(c *fiber.Ctx) error {
	var transactions []models.Transaction
	query := database.DB.Preload("BankAccount").Preload("DestinationBankAccount").Where("type = ?", "transfer")

	// Apply bank account filter if provided
	if bankAccountID := c.Query("bank_account_id"); bankAccountID != "" {
		query = query.Where("bank_account_id = ? OR destination_bank_account_id = ?", bankAccountID, bankAccountID)
	}

	if err := query.Order("date DESC").Find(&transactions).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch transfers",
		})
	}

	// Convert to transfer response format
	var response []models.TransferResponse
	for _, t := range transactions {
		response = append(response, models.TransferResponse{
			ID:            t.ID,
			TransactionID: t.TransactionID,
			Amount:        t.Amount,
			BankAccount: models.BankAccountResponse{
				ID:            t.BankAccount.ID,
				Name:          t.BankAccount.Name,
				AccountNumber: t.BankAccount.AccountNumber,
				BankName:      t.BankAccount.BankName,
				AccountType:   t.BankAccount.AccountType,
				Balance:       t.BankAccount.Balance,
				IsActive:      t.BankAccount.IsActive,
			},
			DestinationBankAccount: models.BankAccountResponse{
				ID:            t.DestinationBankAccount.ID,
				Name:          t.DestinationBankAccount.Name,
				AccountNumber: t.DestinationBankAccount.AccountNumber,
				BankName:      t.DestinationBankAccount.BankName,
				AccountType:   t.DestinationBankAccount.AccountType,
				Balance:       t.DestinationBankAccount.Balance,
				IsActive:      t.DestinationBankAccount.IsActive,
			},
			Description: t.Description,
			Date:        t.Date.Time,
			CreatedAt:   t.CreatedAt,
		})
	}

	return c.JSON(response)
}
