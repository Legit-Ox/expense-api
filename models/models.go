package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// FlexibleDate handles multiple date formats
type FlexibleDate struct {
	time.Time
}

// UnmarshalJSON handles multiple date formats
func (fd *FlexibleDate) UnmarshalJSON(data []byte) error {
	dateStr := strings.Trim(string(data), `"`)
	
	// Try different date formats
	formats := []string{
		"02-01-2006", // dd-mm-yyyy
		"2006-01-02", // yyyy-mm-dd
		"2006-01-02T15:04:05Z", // ISO format
		"2006-01-02T15:04:05Z07:00", // ISO with timezone
		"01/02/2006", // mm/dd/yyyy
		"02/01/2006", // dd/mm/yyyy
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			fd.Time = t
			return nil
		}
	}
	
	return json.Unmarshal(data, &fd.Time)
}

// MarshalJSON outputs in ISO format
func (fd FlexibleDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(fd.Time.Format("2006-01-02T15:04:05Z"))
}

// Value implements the driver.Valuer interface for database storage
func (fd FlexibleDate) Value() (driver.Value, error) {
	return fd.Time, nil
}

// Scan implements the sql.Scanner interface for database retrieval
func (fd *FlexibleDate) Scan(value interface{}) error {
	if value == nil {
		fd.Time = time.Time{}
		return nil
	}
	
	switch v := value.(type) {
	case time.Time:
		fd.Time = v
		return nil
	case string:
		t, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return err
		}
		fd.Time = t
		return nil
	default:
		return fmt.Errorf("cannot scan %T into FlexibleDate", value)
	}
}

// BankAccount represents a bank account
type BankAccount struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Name          string         `json:"name" gorm:"not null"`
	AccountNumber string         `json:"account_number"`
	BankName      string         `json:"bank_name" gorm:"not null"`
	AccountType   string         `json:"account_type" gorm:"not null;check:account_type IN ('checking', 'savings', 'credit', 'investment', 'other')"`
	Balance       float64        `json:"balance" gorm:"default:0"`
	IsActive      bool           `json:"is_active" gorm:"default:true"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// Transaction represents an expense, income, or transfer transaction
type Transaction struct {
	ID                      uint         `json:"id" gorm:"primaryKey"`
	TransactionID           string       `json:"transaction_id" gorm:"index"`
	Amount                  float64      `json:"amount" gorm:"not null"`
	Type                    string       `json:"type" gorm:"not null;check:type IN ('expense', 'income', 'transfer')"`
	CategoryID              *uint        `json:"category_id"` // Nullable for transfers
	Category                Category     `json:"category" gorm:"foreignKey:CategoryID"`
	BankAccountID           uint         `json:"bank_account_id" gorm:"not null"`
	BankAccount             BankAccount  `json:"bank_account" gorm:"foreignKey:BankAccountID"`
	DestinationBankAccountID *uint       `json:"destination_bank_account_id"` // For transfers
	DestinationBankAccount  BankAccount  `json:"destination_bank_account" gorm:"foreignKey:DestinationBankAccountID"`
	Description             string       `json:"description" gorm:"not null"`
	Date                    FlexibleDate `json:"date" gorm:"not null"`
	CreatedAt               time.Time    `json:"created_at"`
	UpdatedAt               time.Time    `json:"updated_at"`
}

// Category represents a transaction category
type Category struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null;unique"`
	Type      string         `json:"type" gorm:"not null;check:type IN ('expense', 'income')"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// BankAccountResponse represents the response structure for bank accounts
type BankAccountResponse struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name"`
	AccountNumber string  `json:"account_number"`
	BankName      string  `json:"bank_name"`
	AccountType   string  `json:"account_type"`
	Balance       float64 `json:"balance"`
	IsActive      bool    `json:"is_active"`
}

// TransactionResponse represents the response structure for transactions
type TransactionResponse struct {
	ID                      uint                  `json:"id"`
	TransactionID           string                `json:"transaction_id"`
	Amount                  float64               `json:"amount"`
	Type                    string                `json:"type"`
	CategoryID              *uint                 `json:"category_id"`
	Category                string                `json:"category"`
	BankAccountID           uint                  `json:"bank_account_id"`
	BankAccount             BankAccountResponse   `json:"bank_account"`
	DestinationBankAccountID *uint                `json:"destination_bank_account_id"`
	DestinationBankAccount  *BankAccountResponse  `json:"destination_bank_account"`
	Description             string                `json:"description"`
	Date                    time.Time             `json:"date"`
	CreatedAt               time.Time             `json:"created_at"`
}

// CategoryResponse represents the response structure for categories
type CategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// AggregateResponse represents the aggregation response
type AggregateResponse struct {
	Categories     map[string]float64 `json:"categories"`
	TotalIncome   float64            `json:"total_income"`
	TotalExpenses float64            `json:"total_expenses"`
	NetAmount     float64            `json:"net_amount"`
}

// BulkTransactionRequest represents a request to create multiple transactions
type BulkTransactionRequest struct {
	Transactions []Transaction `json:"transactions" validate:"required,min=1,max=5000"`
}

// BulkTransactionResponse represents the response for bulk transaction creation
type BulkTransactionResponse struct {
	Success      []TransactionResponse `json:"success"`
	Failed       []BulkTransactionError `json:"failed"`
	TotalCount   int                   `json:"total_count"`
	SuccessCount int                   `json:"success_count"`
	FailedCount  int                   `json:"failed_count"`
}

// BulkTransactionError represents an error for a specific transaction in bulk operation
type BulkTransactionError struct {
	Index       int    `json:"index"`
	Transaction Transaction `json:"transaction"`
	Error       string `json:"error"`
}

// BulkDeleteRequest represents a request to delete multiple transactions
type BulkDeleteRequest struct {
	TransactionIDs []uint `json:"transaction_ids" validate:"required,min=1,max=1000"`
}

// BulkDeleteResponse represents the response for bulk transaction deletion
type BulkDeleteResponse struct {
	Deleted      []uint                `json:"deleted"`
	Failed       []BulkDeleteError     `json:"failed"`
	TotalCount   int                   `json:"total_count"`
	DeletedCount int                   `json:"deleted_count"`
	FailedCount  int                   `json:"failed_count"`
}

// BulkDeleteError represents an error for a specific transaction ID in bulk delete operation
type BulkDeleteError struct {
	TransactionID uint   `json:"transaction_id"`
	Error         string `json:"error"`
}

// CategoryAggregate represents category-wise aggregation data
type CategoryAggregate struct {
	CategoryID       uint    `json:"category_id"`
	CategoryName     string  `json:"category_name"`
	TotalAmount      float64 `json:"total_amount"`
	TransactionCount int     `json:"transaction_count"`
}

// TypeAggregate represents aggregation data for a transaction type (income/expense)
type TypeAggregate struct {
	Categories        []CategoryAggregate `json:"categories"`
	TotalAmount       float64             `json:"total_amount"`
	TotalTransactions int                 `json:"total_transactions"`
}

// DateRange represents a date range
type DateRange struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// AggregateTableResponse represents the response for aggregate table endpoint
type AggregateTableResponse struct {
	DateRange DateRange     `json:"date_range"`
	Income    TypeAggregate `json:"income"`
	Expenses  TypeAggregate `json:"expenses"`
	Summary   struct {
		NetAmount     float64 `json:"net_amount"`
		TotalIncome   float64 `json:"total_income"`
		TotalExpenses float64 `json:"total_expenses"`
	} `json:"summary"`
}

// TransferRequest represents a request to create a transfer between accounts
type TransferRequest struct {
	Amount                  float64      `json:"amount" validate:"required,gt=0"`
	BankAccountID           uint         `json:"bank_account_id" validate:"required"`
	DestinationBankAccountID uint        `json:"destination_bank_account_id" validate:"required"`
	Description             string       `json:"description" validate:"required"`
	Date                    FlexibleDate `json:"date" validate:"required"`
	TransactionID           string       `json:"transaction_id"`
}

// TransferResponse represents the response for a transfer transaction
type TransferResponse struct {
	ID                      uint                 `json:"id"`
	TransactionID           string               `json:"transaction_id"`
	Amount                  float64              `json:"amount"`
	BankAccount             BankAccountResponse  `json:"bank_account"`
	DestinationBankAccount  BankAccountResponse  `json:"destination_bank_account"`
	Description             string               `json:"description"`
	Date                    time.Time            `json:"date"`
	CreatedAt               time.Time            `json:"created_at"`
} 