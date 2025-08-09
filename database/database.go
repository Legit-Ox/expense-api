package database

import (
	"log"
	"os"
	"strings"

	"expense-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect establishes database connection
func Connect() {
	var err error
	dbURL := os.Getenv("DB_URL")

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

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")
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