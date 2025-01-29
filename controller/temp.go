package controller

import (
	"github.com/HealthMe-pls/medic-go-api/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Create a new TempShop
func CreateTempShop(db *gorm.DB, c *fiber.Ctx) error {
	var tempShop model.TempShop

	// Parse request body into tempShop
	if err := c.BodyParser(&tempShop); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Save to database
	if err := db.Create(&tempShop).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create TempShop",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(tempShop)
}

// Get all TempShops
func GetAllTempShops(db *gorm.DB, c *fiber.Ctx) error {
	var tempShops []model.TempShop

	// Fetch all records
	if err := db.Preload("DeletePhoto").Preload("DeleteSocial").Preload("DeleteMenu").
		Preload("ShopMenus").Preload("SocialMedia").Preload("Photos").
		Preload("ShopCategory").Find(&tempShops).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve TempShops",
		})
	}

	return c.JSON(tempShops)
}

// Get TempShop by TempID
func GetTempShopByID(db *gorm.DB, c *fiber.Ctx) error {
	tempID := c.Params("id")
	var tempShop model.TempShop

	// Fetch TempShop by TempID
	if err := db.Preload("DeletePhoto").Preload("DeleteSocial").Preload("DeleteMenu").
		Preload("ShopMenus").Preload("SocialMedia").Preload("Photos").
		Preload("ShopCategory").First(&tempShop, tempID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "TempShop not found",
		})
	}

	return c.JSON(tempShop)
}

// Get TempShops by ShopID
func GetTempShopsByShopID(db *gorm.DB, c *fiber.Ctx) error {
	shopID := c.Params("shop_id")
	var tempShops []model.TempShop

	// Fetch TempShops by ShopID
	if err := db.Where("shop_id = ?", shopID).Preload("DeletePhoto").Preload("DeleteSocial").Preload("DeleteMenu").
		Preload("ShopMenus").Preload("SocialMedia").Preload("Photos").
		Preload("ShopCategory").Find(&tempShops).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve TempShops for the given ShopID",
		})
	}

	return c.JSON(tempShops)
}

// Update a TempShop by TempID
func UpdateTempShop(db *gorm.DB, c *fiber.Ctx) error {
	tempID := c.Params("id")
	var tempShop model.TempShop

	// Check if TempShop exists
	if err := db.First(&tempShop, tempID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "TempShop not found",
		})
	}

	// Parse request body
	if err := c.BodyParser(&tempShop); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update TempShop
	if err := db.Save(&tempShop).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update TempShop",
		})
	}

	return c.JSON(tempShop)
}

func UpdateTempShopByShopID(db *gorm.DB, c *fiber.Ctx) error {
	shopID := c.Params("shop_id")
	var tempShop model.TempShop

	// Check if a TempShop exists with the given ShopID
	if err := db.Where("shop_id = ?", shopID).First(&tempShop).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No TempShop found for the given ShopID",
		})
	}

	// Parse the request body into tempShop
	if err := c.BodyParser(&tempShop); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update the TempShop
	if err := db.Save(&tempShop).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update TempShop",
		})
	}

	return c.JSON(tempShop)
}


// Delete a TempShop by TempID
func DeleteTempShop(db *gorm.DB, c *fiber.Ctx) error {
	tempID := c.Params("id")

	// Delete TempShop
	if err := db.Delete(&model.TempShop{}, tempID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete TempShop",
		})
	}

	return c.SendString("TempShop deleted successfully")
}
