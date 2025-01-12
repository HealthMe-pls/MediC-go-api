package controller

import (
	"github.com/HealthMe-pls/medic-go-api/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetEntrepreneur(db *gorm.DB,c *fiber.Ctx) error {
	var entrepreneur []model.Entrepreneur
	db.Find(&entrepreneur)
	return c.JSON(entrepreneur)
}

func GetEntrepreneurByUsername(db *gorm.DB, c *fiber.Ctx) error {
    username := c.Params("username")
    var entrepreneur model.Entrepreneur
    
    // Query the database for the entrepreneur with the provided username
    if err := db.First(&entrepreneur, "username = ?", username).Error; err != nil {
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
    // Get the username parameter from the URL
    username := c.Params("username")
    var entrepreneur model.Entrepreneur

    // Find the entrepreneur by username
    if err := db.First(&entrepreneur, "username = ?", username).Error; err != nil {
        return c.Status(fiber.StatusNotFound).SendString("Entrepreneur not found")
    }

    // Parse the updated details from the request body
    if err := c.BodyParser(&entrepreneur); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString("Failed to parse request body")
    }

    // Save the updated entrepreneur details to the database
    if result := db.Save(&entrepreneur); result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to update entrepreneur")
    }

    // Return the updated entrepreneur as a JSON response
    return c.JSON(entrepreneur)
}

// func DeleteEntrepreneur(db *gorm.DB, c *fiber.Ctx) error {
//     // Get the username parameter from the URL
//     username := c.Params("username")

//     // Delete the entrepreneur from the database by their username
//     if result := db.Delete(&model.Entrepreneur{}, "username = ?", username); result.Error != nil {
//         return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete entrepreneur")
//     }

//     // Return success message
//     return c.SendString("Entrepreneur successfully deleted")
// }

func DeleteEntrepreneurAndShops(db *gorm.DB, c *fiber.Ctx) error {
    // Get the entrepreneur's username from the URL parameter
    username := c.Params("username")

    // Delete all shops associated with the entrepreneur
    if err := db.Where("entrepreneur_username = ?", username).Delete(&model.Shop{}).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete shops")
    }

    // Now delete the entrepreneur
    if err := db.Where("username = ?", username).Delete(&model.Entrepreneur{}).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete entrepreneur")
    }

    return c.SendString("Entrepreneur and associated shops successfully deleted")
}

