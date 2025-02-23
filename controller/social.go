package controller

import (

	"github.com/HealthMe-pls/medic-go-api/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CreateSocialMedia creates a new SocialMedia entry
func CreateSocialMediaByAdmin(db *gorm.DB, c *fiber.Ctx) error {
	var socialMedia model.SocialMedia
	if err := c.BodyParser(&socialMedia); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	if err := db.Create(&socialMedia).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create social media",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(socialMedia)
}


func CreateSocialByEntrepreneur(db *gorm.DB, c *fiber.Ctx) error {
	entrepreneurID := c.Params("entrepreneur_id")
	shopID := c.Params("shop_id")

	// Check if the entrepreneur exists
	var entrepreneur model.Entrepreneur
	if err := db.First(&entrepreneur, "id = ?", entrepreneurID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Entrepreneur not found",
			"details": err.Error(),
		})
	}

	// Check if the shop exists and belongs to the entrepreneur
	var shop model.Shop
	if err := db.First(&shop, "id = ? AND entrepreneur_id = ?", shopID, entrepreneurID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Shop not found or does not belong to the entrepreneur",
			"details": err.Error(),
		})
	}

	// Parse request body
	var social model.SocialMedia
	if err := c.BodyParser(&social); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to parse request body",
			"details": err.Error(),
		})
	}

	// Ensure IsPublic is set to false and assign Shop ID
	social.IsPublic = false
	social.ShopID = shop.ID

	// Save the social media
	if err := db.Create(&social).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create social media",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(social)
}

// -----social media

// GetSocialMedia retrieves a SocialMedia entry by ID
func GetSocialMedia(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var socialMedia model.SocialMedia
	// if err := db.First(&socialMedia, id).Error; err != nil {
	// 	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
	// 		"error": "Social media not found",
	// 	})
	// }
	db.First(&socialMedia, id)
	return c.JSON(socialMedia)
}

// GetSocialMediaByShopID retrieves SocialMedia entries by Shop ID
func GetSocialMediaByShopID(db *gorm.DB, c *fiber.Ctx) error {
	shopID := c.Params("shop_id")
	var socialMedias []map[string]interface{}

	// Query specific fields
	if err := db.Model(&model.SocialMedia{}).
		Select("id, platform, link").
		Where("shop_id = ?", shopID).
		Find(&socialMedias).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve social media entries",
		})
	}

	return c.JSON(socialMedias)
}

// UpdateSocialMedia updates a SocialMedia entry by ID
func UpdateSocialMedia(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var socialMedia model.SocialMedia
	if err := db.First(&socialMedia, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Social media not found",
		})
	}

	if err := c.BodyParser(&socialMedia); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if err := db.Save(&socialMedia).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update social media",
		})
	}
	return c.JSON(socialMedia)
}

// DeleteSocialMedia deletes a SocialMedia entry by ID
func DeleteSocialMedia(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	if err := db.Delete(&model.SocialMedia{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete social media",
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}