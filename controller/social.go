package controller

import (
	"fmt"
	"strconv"
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

func CreateSocialWithTemp(db *gorm.DB, c *fiber.Ctx) error {
	social := new(model.SocialMedia)
	if err := c.BodyParser(social); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to parse request body",
			"details": err.Error(),
		})
	}

	// Set is_public to false before saving
	social.IsPublic = false

	// Save the social media entry in the SocialMedia table
	if result := db.Create(&social); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create social media",
			"details": result.Error.Error(),
		})
	}

	// Ensure TempID is not nil before dereferencing
	var tempID uint
	if social.TempID != nil {
		tempID = *social.TempID
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "TempID is required",
		})
	}

	// Create a corresponding TempSocial entry
	tempSocial := model.TempSocial{
		TempID:   tempID,
		SocialID: social.ID,
		Name:     social.Name,
		Platform: social.Platform,
		Link:     social.Link,
	}

	if result := db.Create(&tempSocial); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create temp social media",
			"details": result.Error.Error(),
		})
	}

	// Return created social media and temp social media
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"social":     social,
		"tempSocial": tempSocial,
	})
}

func GetShopIDBySocialID(db *gorm.DB, c *fiber.Ctx) error {
	socialID := c.Params("social_id")

	// Check if the social media entry exists
	var socialMedia model.SocialMedia
	if err := db.First(&socialMedia, "id = ?", socialID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Social media entry not found",
			"details": err.Error(),
		})
	}

	// Return the ShopID associated with the Social Media
	return c.JSON(fiber.Map{
		"shop_id": socialMedia.ShopID,
	})
}

func UpdateSocialBySocialID(db *gorm.DB, c *fiber.Ctx) error {
	socialID, err := strconv.Atoi(c.Params("social_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid social ID",
		})
	}

	// Fetch the SocialMedia record
	var social model.SocialMedia
	if err := db.First(&social, socialID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "SocialMedia not found",
		})
	}

	// Check if TempSocial exists for this SocialID
	var tempSocial model.TempSocial
	if err := db.Where("social_id = ?", social.ID).First(&tempSocial).Error; err != nil {
		// If not found, create a new TempSocial entry
		tempSocial = model.TempSocial{
			SocialID: social.ID,
			TempID:   *social.TempID, // Ensure TempID is not nil
			Name:     social.Name,
			Platform: social.Platform,
			Link:     social.Link,
		}
		if err := db.Create(&tempSocial).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create temp social media",
			})
		}
	} else {
		// If found, update existing TempSocial
		tempSocial.Name = social.Name
		tempSocial.Platform = social.Platform
		tempSocial.Link = social.Link

		if err := db.Save(&tempSocial).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update temp social media",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Temp social media updated successfully",
		"tempSocial": tempSocial,
	})
}

func UpdateSocialFromTemp(db *gorm.DB, tempID uint) error {
	// Fetch the TempSocial record by TempID
	var tempSocial model.TempSocial
	if err := db.First(&tempSocial, "temp_id = ?", tempID).Error; err != nil {
		return fmt.Errorf("TempSocial not found: %w", err)
	}

	// Fetch the corresponding SocialMedia record by SocialID
	var social model.SocialMedia
	if err := db.First(&social, "id = ?", tempSocial.SocialID).Error; err != nil {
		return fmt.Errorf("SocialMedia not found for SocialID %d: %w", tempSocial.SocialID, err)
	}

	// Update SocialMedia with values from TempSocial
	social.Name = tempSocial.Name
	social.Platform = tempSocial.Platform
	social.Link = tempSocial.Link

	// Save the updated SocialMedia record
	if err := db.Save(&social).Error; err != nil {
		return fmt.Errorf("failed to update SocialMedia: %w", err)
	}

	return nil
}
