package handlers

import (
	"expense-api/database"
	"expense-api/models"

	"github.com/gofiber/fiber/v2"
)

// CreateCategory handles POST /categories
func CreateCategory(c *fiber.Ctx) error {
	var category models.Category

	if err := c.BodyParser(&category); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate category type
	if category.Type != "expense" && category.Type != "income" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Type must be either 'expense' or 'income'",
		})
	}

	// Check if category name already exists
	var existingCategory models.Category
	if err := database.DB.Where("name = ?", category.Name).First(&existingCategory).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Category with this name already exists",
		})
	}

	if err := database.DB.Create(&category).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create category",
		})
	}

	return c.Status(201).JSON(category)
}

// GetCategories handles GET /categories
func GetCategories(c *fiber.Ctx) error {
	var categories []models.Category

	if err := database.DB.Find(&categories).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch categories",
		})
	}

	// Convert to response format
	var response []models.CategoryResponse
	for _, c := range categories {
		response = append(response, models.CategoryResponse{
			ID:   c.ID,
			Name: c.Name,
			Type: c.Type,
		})
	}

	return c.JSON(response)
}

// GetCategory handles GET /categories/:id
func GetCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	var category models.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	response := models.CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
		Type: category.Type,
	}

	return c.JSON(response)
}

// DeleteCategory handles DELETE /categories/:id
func DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	var category models.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	// Check if category is being used by any transactions
	var count int64
	database.DB.Model(&models.Transaction{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Cannot delete category that has associated transactions",
		})
	}

	if err := database.DB.Delete(&category).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete category",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Category deleted successfully",
	})
}

// UpdateCategory handles PUT /categories/:id
func UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	var category models.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Category not found",
		})
	}

	var updateData map[string]interface{}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate category type if provided
	if categoryType, exists := updateData["type"]; exists {
		if categoryType != "expense" && categoryType != "income" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Type must be either 'expense' or 'income'",
			})
		}
	}

	// Check if name already exists (excluding current category)
	if name, exists := updateData["name"]; exists {
		var existingCategory models.Category
		if err := database.DB.Where("name = ? AND id != ?", name, id).First(&existingCategory).Error; err == nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Category with this name already exists",
			})
		}
	}

	if err := database.DB.Model(&category).Updates(updateData).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update category",
		})
	}

	// Load updated category
	database.DB.First(&category, id)

	return c.JSON(category)
}
