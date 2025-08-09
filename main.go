package main

import (
	"log"
	"os"

	"expense-api/database"
	"expense-api/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	database.Connect()
	database.Migrate()
	database.SeedDefaultCategories()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"message": "Expense API is running",
		})
	})

	// API routes
	api := app.Group("/api")

	// Transaction routes
	transactions := api.Group("/transactions")
	transactions.Post("/", handlers.CreateTransaction)
	transactions.Get("/", handlers.GetTransactions)
	transactions.Get("/aggregate", handlers.GetTransactionsAggregate)
	transactions.Get("/date-range", handlers.GetTransactionsByDateRange)
	transactions.Get("/:id", handlers.GetTransaction)
	transactions.Put("/:id", handlers.UpdateTransaction)
	transactions.Delete("/:id", handlers.DeleteTransaction)

	// Category routes
	categories := api.Group("/categories")
	categories.Post("/", handlers.CreateCategory)
	categories.Get("/", handlers.GetCategories)
	categories.Get("/:id", handlers.GetCategory)
	categories.Put("/:id", handlers.UpdateCategory)
	categories.Delete("/:id", handlers.DeleteCategory)

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
