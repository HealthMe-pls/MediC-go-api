package controller

import (
	"github.com/HealthMe-pls/medic-go-api/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetAllWorkshops retrieves all workshops
func GetWorkshops(db *gorm.DB, c *fiber.Ctx) error {
	var workshops []model.Workshop
	if err := db.Find(&workshops).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch workshops",
		})
	}
	return c.JSON(workshops)
}

// GetWorkshopByID retrieves a workshop by its ID
func GetWorkshopByID(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var workshop model.Workshop
	if err := db.First(&workshop, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Workshop not found",
		})
	}
	// db.First(&workshop, "id = ?", id)
	return c.JSON(workshop)
}

// CreateWorkshop creates a new workshop
func CreateWorkshop(db *gorm.DB, c *fiber.Ctx) error {
	var workshop model.Workshop
	if err := c.BodyParser(&workshop); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
			"details": err.Error(),
		})
	}
	if err := db.Create(&workshop).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create workshop",
			"details": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(workshop)
}

// UpdateWorkshop updates an existing workshop
func UpdateWorkshop(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var workshop model.Workshop
	if err := db.First(&workshop, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Workshop not found",
		})
	}
	if err := c.BodyParser(&workshop); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	if err := db.Save(&workshop).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update workshop",
		})
	}
	return c.JSON(workshop)
}

// DeleteWorkshop deletes a workshop by its ID
func DeleteWorkshop(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	if err := db.Delete(&model.Workshop{}, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete workshop",
		})
	}
	return c.SendString("Workshop successfully deleted")
}
