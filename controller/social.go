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
// UpdateSocialMedia updates a SocialMedia entry by ID and also updates the TempSocial entry
func UpdateSocialMedia(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var socialMedia model.SocialMedia

	// Retrieve the existing SocialMedia entry by ID
	if err := db.First(&socialMedia, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Social media not found",
		})
	}

	// Parse the request body into the SocialMedia struct
	if err := c.BodyParser(&socialMedia); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Update the SocialMedia entry
	if err := db.Save(&socialMedia).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update social media",
		})
	}

	// Now update or create the corresponding TempSocial entry
	var tempSocial model.TempSocial
	// Check if a TempSocial exists for this SocialMedia
	if err := db.Where("social_id = ?", socialMedia.ID).First(&tempSocial).Error; err != nil {
		// If TempSocial doesn't exist, create a new entry
		tempSocial = model.TempSocial{
			SocialID: socialMedia.ID,
			Name:     socialMedia.Name,
			Platform: socialMedia.Platform,
			Link:     socialMedia.Link,
		}
		if err := db.Create(&tempSocial).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create TempSocial",
			})
		}
	} else {
		// If TempSocial exists, update it
		tempSocial.Name = socialMedia.Name
		tempSocial.Platform = socialMedia.Platform
		tempSocial.Link = socialMedia.Link
		if err := db.Save(&tempSocial).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update TempSocial",
			})
		}
	}

	// Return the updated SocialMedia and TempSocial as JSON
	return c.JSON(fiber.Map{
		"social_media": socialMedia,
		"temp_social":  tempSocial,
	})
}


// DeleteSocialMedia deletes a SocialMedia entry and its related TempSocial entry
func DeleteSocialMedia(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")

	// Begin a transaction to ensure atomicity
	tx := db.Begin()

	// Delete TempSocial entries related to the SocialMedia entry
	if err := tx.Where("social_id = ?", id).Delete(&model.TempSocial{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete associated TempSocial",
		})
	}

	// Delete the SocialMedia entry
	if err := tx.Delete(&model.SocialMedia{}, id).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete social media",
		})
	}

	// Commit the transaction
	tx.Commit()

	return c.SendStatus(fiber.StatusNoContent)
}

func CreateSocialWithTemp(db *gorm.DB, c *fiber.Ctx, isPublic bool) error {
    social := new(model.SocialMedia)
    
    // Parse request body
    if err := c.BodyParser(social); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error":   "Failed to parse request body",
            "details": err.Error(),
        })
    }

    // Assign isPublic from function parameter
    social.IsPublic = isPublic

    // Find TempShop that has the same ShopID as the social media
    var tempShop model.TempShop
    if err := db.Where("shop_id = ?", social.ShopID).First(&tempShop).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error":   "TempShop not found for this ShopID",
            "details": err.Error(),
        })
    }

    // Assign the found TempID to the new social media entry
    social.TempID = &tempShop.TempID

    // Save the SocialMedia entry in the database
    if result := db.Create(&social); result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error":   "Failed to create social media",
            "details": result.Error.Error(),
        })
    }

    // Create a corresponding TempSocial entry
    tempSocial := model.TempSocial{
        TempID:   tempShop.TempID, // Use the found TempID
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
