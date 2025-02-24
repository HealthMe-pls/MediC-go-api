package controller

import (
	"strconv"
	"fmt"
	"github.com/HealthMe-pls/medic-go-api/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)
func UpdateShopFromTemp(db *gorm.DB, tempID uint) error {
	// Fetch the TempShop record by TempID
	var tempShop model.TempShop
	if err := db.First(&tempShop, "temp_id = ?", tempID).Error; err != nil {
		return fmt.Errorf("TempShop not found: %w", err)
	}

	// Ensure the TempShop has a valid ShopID
	if tempShop.ShopID == nil {
		return fmt.Errorf("TempShop with TempID %d has no associated ShopID", tempID)
	}

	// Fetch the corresponding Shop record
	var shop model.Shop
	if err := db.First(&shop, "id = ?", *tempShop.ShopID).Error; err != nil {
		return fmt.Errorf("Shop not found for ShopID %d: %w", *tempShop.ShopID, err)
	}

	// Update the Shop with values from TempShop
	shop.Name = tempShop.Name
	shop.Description = tempShop.Description
	if tempShop.ShopCategoryID != nil {
		shop.ShopCategoryID = *tempShop.ShopCategoryID
	}

	// Save the updated Shop record
	if err := db.Save(&shop).Error; err != nil {
		return fmt.Errorf("failed to update Shop: %w", err)
	}

	return nil
}

func HandleUpdateShopFromTemp(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")

	// Convert id to uint
	tempID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid TempShop ID",
		})
	}

	// Call the function to update the shop
	if err := UpdateShopFromTemp(db, uint(tempID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update shop",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"message": "Shop updated successfully"})
}

func UpdateStatusToPublicByTempID(db *gorm.DB, tempID uint) error {
	// Update IsPublic for SocialMedia by TempID
	if err := db.Model(&model.SocialMedia{}).
		Where("temp_id = ?", tempID).
		Update("is_public", true).Error; err != nil {
		return err
	}

	// Update IsPublic for ShopMenu by TempID
	if err := db.Model(&model.ShopMenu{}).
		Where("temp_id = ?", tempID).
		Update("is_public", true).Error; err != nil {
		return err
	}

	// Update IsPublic for Photo by TempID
	if err := db.Model(&model.Photo{}).
		Where("temp_id = ?", tempID).
		Update("is_public", true).Error; err != nil {
		return err
	}

	return nil
}

func GetTempIDByShopID(db *gorm.DB, c *fiber.Ctx) error {
	shopID := c.Params("shop_id")

	// Check if the TempShop entry exists with the given ShopID
	var tempShop model.TempShop
	if err := db.Where("shop_id = ?", shopID).First(&tempShop).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "TempShop not found for given ShopID",
			"details": err.Error(),
		})
	}

	// Return the TempID associated with the Shop
	return c.JSON(fiber.Map{
		"temp_id": tempShop.TempID,
	})
}
