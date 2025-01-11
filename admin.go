package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Get all Admins
func GetAdmins(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var admins []Admin
		db.Find(&admins)
		return c.JSON(admins)
	}
}

// Get Admin by Username
func GetAdminByUsername(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Params("username")
		var admin Admin
		if err := db.First(&admin, "username = ?", username).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Admin not found")
		}
		return c.JSON(admin)
	}
}

// Create Admin
func CreateAdmin(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		admin := new(Admin)
		if err := c.BodyParser(admin); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Failed to parse request body")
		}
		if result := db.Create(&admin); result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to create admin")
		}
		return c.JSON(admin)
	}
}

// Update Admin
func UpdateAdmin(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Params("username")
		var admin Admin
		if err := db.First(&admin, "username = ?", username).Error; err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Admin not found")
		}

		// Parse updated data from request body
		if err := c.BodyParser(&admin); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Failed to parse request body")
		}
		if result := db.Save(&admin); result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to update admin")
		}
		return c.JSON(admin)
	}
}

// Delete Admin
func DeleteAdmin(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		username := c.Params("username")
		if result := db.Delete(&Admin{}, "username = ?", username); result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete admin")
		}
		return c.SendString("Admin successfully deleted")
	}
}
