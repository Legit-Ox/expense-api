package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"expense-api/models"
)

// CreateBankAccount creates a new bank account
func CreateBankAccount(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var bankAccount models.BankAccount
		
		if err := c.BodyParser(&bankAccount); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}

		// Validate required fields
		if bankAccount.Name == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Bank account name is required",
			})
		}

		if bankAccount.BankName == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Bank name is required",
			})
		}

		if bankAccount.AccountType == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Account type is required",
			})
		}

		// Validate account type
		validTypes := []string{"checking", "savings", "credit", "investment", "other"}
		isValidType := false
		for _, validType := range validTypes {
			if bankAccount.AccountType == validType {
				isValidType = true
				break
			}
		}
		if !isValidType {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid account type. Must be one of: checking, savings, credit, investment, other",
			})
		}

		// Create bank account
		if err := db.Create(&bankAccount).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create bank account",
			})
		}

		// Return response
		response := models.BankAccountResponse{
			ID:            bankAccount.ID,
			Name:          bankAccount.Name,
			AccountNumber: bankAccount.AccountNumber,
			BankName:      bankAccount.BankName,
			AccountType:   bankAccount.AccountType,
			Balance:       bankAccount.Balance,
			IsActive:      bankAccount.IsActive,
		}

		return c.Status(fiber.StatusCreated).JSON(response)
	}
}

// GetBankAccounts retrieves all bank accounts
func GetBankAccounts(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var bankAccounts []models.BankAccount
		
		// Get active accounts by default, unless include_inactive=true
		query := db.Where("is_active = ?", true)
		if c.Query("include_inactive") == "true" {
			query = db
		}

		if err := query.Find(&bankAccounts).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve bank accounts",
			})
		}

		// Convert to response format
		var responses []models.BankAccountResponse
		for _, account := range bankAccounts {
			responses = append(responses, models.BankAccountResponse{
				ID:            account.ID,
				Name:          account.Name,
				AccountNumber: account.AccountNumber,
				BankName:      account.BankName,
				AccountType:   account.AccountType,
				Balance:       account.Balance,
				IsActive:      account.IsActive,
			})
		}

		return c.JSON(responses)
	}
}

// GetBankAccount retrieves a specific bank account by ID
func GetBankAccount(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid bank account ID",
			})
		}

		var bankAccount models.BankAccount
		if err := db.First(&bankAccount, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Bank account not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve bank account",
			})
		}

		response := models.BankAccountResponse{
			ID:            bankAccount.ID,
			Name:          bankAccount.Name,
			AccountNumber: bankAccount.AccountNumber,
			BankName:      bankAccount.BankName,
			AccountType:   bankAccount.AccountType,
			Balance:       bankAccount.Balance,
			IsActive:      bankAccount.IsActive,
		}

		return c.JSON(response)
	}
}

// UpdateBankAccount updates a bank account
func UpdateBankAccount(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid bank account ID",
			})
		}

		var existingAccount models.BankAccount
		if err := db.First(&existingAccount, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Bank account not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve bank account",
			})
		}

		var updateData models.BankAccount
		if err := c.BodyParser(&updateData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}

		// Validate account type if provided
		if updateData.AccountType != "" {
			validTypes := []string{"checking", "savings", "credit", "investment", "other"}
			isValidType := false
			for _, validType := range validTypes {
				if updateData.AccountType == validType {
					isValidType = true
					break
				}
			}
			if !isValidType {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid account type. Must be one of: checking, savings, credit, investment, other",
				})
			}
		}

		// Update fields
		if updateData.Name != "" {
			existingAccount.Name = updateData.Name
		}
		if updateData.AccountNumber != "" {
			existingAccount.AccountNumber = updateData.AccountNumber
		}
		if updateData.BankName != "" {
			existingAccount.BankName = updateData.BankName
		}
		if updateData.AccountType != "" {
			existingAccount.AccountType = updateData.AccountType
		}
		if updateData.Balance != 0 {
			existingAccount.Balance = updateData.Balance
		}
		// Handle IsActive explicitly since it's a boolean
		existingAccount.IsActive = updateData.IsActive

		if err := db.Save(&existingAccount).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update bank account",
			})
		}

		response := models.BankAccountResponse{
			ID:            existingAccount.ID,
			Name:          existingAccount.Name,
			AccountNumber: existingAccount.AccountNumber,
			BankName:      existingAccount.BankName,
			AccountType:   existingAccount.AccountType,
			Balance:       existingAccount.Balance,
			IsActive:      existingAccount.IsActive,
		}

		return c.JSON(response)
	}
}

// DeleteBankAccount soft deletes a bank account
func DeleteBankAccount(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid bank account ID",
			})
		}

		// Check if bank account exists
		var bankAccount models.BankAccount
		if err := db.First(&bankAccount, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Bank account not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve bank account",
			})
		}

		// Check if there are any transactions associated with this account
		var transactionCount int64
		if err := db.Model(&models.Transaction{}).Where("bank_account_id = ? OR destination_bank_account_id = ?", id, id).Count(&transactionCount).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to check for associated transactions",
			})
		}

		if transactionCount > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot delete bank account with associated transactions. Consider deactivating instead.",
			})
		}

		// Soft delete the bank account
		if err := db.Delete(&bankAccount).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete bank account",
			})
		}

		return c.Status(fiber.StatusNoContent).Send(nil)
	}
}
