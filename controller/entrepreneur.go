package controller

import (
	"fmt"
	"strconv"
	"github.com/HealthMe-pls/medic-go-api/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetEntrepreneur(db *gorm.DB,c *fiber.Ctx) error {
	var entrepreneur []model.Entrepreneur
	db.Find(&entrepreneur)
	return c.JSON(entrepreneur)
}

func GetEntrepreneurByID(db *gorm.DB, c *fiber.Ctx) error {
    id := c.Params("id")
    var entrepreneur model.Entrepreneur
    
    // Query the database for the entrepreneur with the provided id
    if err := db.First(&entrepreneur, "id = ?", id).Error; err != nil {
        // If an error occurs (e.g., no entrepreneur found), return a 404
        return c.Status(fiber.StatusNotFound).SendString("Entrepreneur not found")
    }

    // If successful, return the entrepreneur data as a JSON response
    return c.JSON(entrepreneur)
}

func CreateEntrepreneur(db *gorm.DB, c *fiber.Ctx) error {
	// Parse the request body into the Entrepreneur struct
	entrepreneur := new(model.Entrepreneur)
	if err := c.BodyParser(entrepreneur); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	// Check if the username already exists
    var existingEntrepreneur model.Entrepreneur
	err := db.Where("username = ?", entrepreneur.Username).First(&existingEntrepreneur).Error
	if err == nil {
		// Username already exists
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Username already exists",
		})
	} else if err != nil && err != gorm.ErrRecordNotFound {
		// Database error occurred while checking
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check username",
		})
	}

	// Save the Entrepreneur to the database
	if result := db.Create(&entrepreneur); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create entrepreneur",
		})
	}

	// Return the created Entrepreneur as a JSON response
	return c.Status(fiber.StatusCreated).JSON(entrepreneur)
}


func UpdateEntrepreneur(db *gorm.DB, c *fiber.Ctx) error {
    // Get the id parameter from the URL
    id := c.Params("id")
    var entrepreneur model.Entrepreneur

    // Find the entrepreneur by id
    if err := db.First(&entrepreneur, "id = ?", id).Error; err != nil {
        return c.Status(fiber.StatusNotFound).SendString("Entrepreneur not found")
    }

    // Parse the updated details from the request body
    if err := c.BodyParser(&entrepreneur); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString("Failed to parse request body")
    }
	encryptedPassword, err := EncryptPassword(entrepreneur.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to encrypt password",
		})
	}
	entrepreneur.Password = encryptedPassword
    // Save the updated entrepreneur details to the database
    if result := db.Save(&entrepreneur); result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to update entrepreneur")
    }

    // Return the updated entrepreneur as a JSON response
    return c.JSON(entrepreneur)
}
func DeleteEntrepreneurByID(db *gorm.DB, c *fiber.Ctx) error {
	// Get the entrepreneur ID parameter from the URL
	idParam := c.Params("id")

	// Convert entrepreneur ID to uint
	entrepreneurID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid entrepreneur ID",
		})
	}

	// Begin a transaction to ensure atomicity
	tx := db.Begin()

	// Step 1: Check if the Entrepreneur exists
	var entrepreneur model.Entrepreneur
	if err := tx.First(&entrepreneur, "id = ?", entrepreneurID).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Entrepreneur not found",
			"details": err.Error(),
		})
	}

	// Step 2: Retrieve all shops associated with the entrepreneur
	var shops []model.Shop
	if err := tx.Where("entrepreneur_id = ?", entrepreneurID).Find(&shops).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve shops",
		})
	}

	// Step 3: Delete each shop using DeleteShopByID within the same transaction
	for _, shop := range shops {
		if err := DeleteShopByID(tx, shop.ID); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   fmt.Sprintf("Failed to delete shop ID %d", shop.ID),
				"details": err.Error(),
			})
		}
	}

	// Step 4: Delete the Entrepreneur after all shops are deleted (within the same transaction)
	if err := tx.Where("id = ?", entrepreneurID).Delete(&model.Entrepreneur{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to delete entrepreneur",
			"details": err.Error(),
		})
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit transaction",
		})
	}

	// Return success message
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Entrepreneur deleted successfully",
	})
}

