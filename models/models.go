package models

import (
	"time"

	"gorm.io/gorm"
)

// Transaction represents an expense or income transaction
type Transaction struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Amount      float64   `json:"amount" gorm:"not null"`
	Type        string    `json:"type" gorm:"not null;check:type IN ('expense', 'income')"`
	CategoryID  uint      `json:"category_id" gorm:"not null"`
	Category    Category  `json:"category" gorm:"foreignKey:CategoryID"`
	Description string    `json:"description" gorm:"not null"`
	Date        time.Time `json:"date" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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

// TransactionResponse represents the response structure for transactions
type TransactionResponse struct {
	ID          uint      `json:"id"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"`
	CategoryID  uint      `json:"category_id"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
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