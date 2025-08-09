package database

import (
	"log"
	"os"
	"strings"
	"time"

	"expense-api/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect establishes database connection with retry logic
func Connect() {
	var err error
	dbURL := os.Getenv("DB_URL")

	// Fallback to DATABASE_URL if DB_URL is not set
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_URL")
	}

	// If still no database URL, use SQLite as fallback
	if dbURL == "" {
		log.Println("No database URL found, using SQLite as fallback")
		dbURL = "sqlite://expense.db"
	}

	log.Printf("Attempting to connect to database with URL: %s", maskPassword(dbURL))

	// Retry connection logic
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			log.Printf("Retry attempt %d/%d", i+1, maxRetries)
			time.Sleep(time.Duration(i*2) * time.Second)
		}

		if strings.HasPrefix(dbURL, "sqlite://") {
			// SQLite connection
			dbPath := strings.TrimPrefix(dbURL, "sqlite://")
			DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			})
		} else {
			// PostgreSQL connection
			DB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			})
		}

		if err == nil {
			break
		}

		log.Printf("Database connection attempt %d failed: %v", i+1, err)
	}

	if err != nil {
		log.Fatal("Failed to connect to database after all retries:", err)
	}

	log.Println("Database connected successfully")
}

// maskPassword masks the password in database URL for logging
func maskPassword(dbURL string) string {
	if strings.Contains(dbURL, "@") {
		parts := strings.Split(dbURL, "@")
		if len(parts) == 2 {
			userPass := strings.Split(parts[0], "://")
			if len(userPass) == 2 {
				userParts := strings.Split(userPass[1], ":")
				if len(userParts) == 2 {
					return userPass[0] + "://" + userParts[0] + ":***@" + parts[1]
				}
			}
		}
	}
	return dbURL
}

// Migrate runs database migrations
func Migrate() {
	err := DB.AutoMigrate(&models.Category{}, &models.Transaction{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migrated successfully")
}

// SeedDefaultCategories populates the database with default categories
func SeedDefaultCategories() {
	var count int64
	DB.Model(&models.Category{}).Count(&count)

	if count > 0 {
		log.Println("Categories already seeded, skipping...")
		return
	}

	defaultCategories := []models.Category{
		{Name: "Food", Type: "expense"},
		{Name: "Transport", Type: "expense"},
		{Name: "Bills", Type: "expense"},
		{Name: "Shopping", Type: "expense"},
		{Name: "Salary", Type: "income"},
		{Name: "Freelance", Type: "income"},
		{Name: "Investments", Type: "income"},
	}

	for _, category := range defaultCategories {
		if err := DB.Create(&category).Error; err != nil {
			log.Printf("Failed to create category %s: %v", category.Name, err)
		}
	}

	log.Printf("Seeded %d default categories", len(defaultCategories))
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
