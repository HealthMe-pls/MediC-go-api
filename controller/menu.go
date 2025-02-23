package controller

import (

	"github.com/HealthMe-pls/medic-go-api/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)


// shop menu
// CreateShopMenu creates a new ShopMenu entry
func CreateShopMenuByAdmin(db *gorm.DB, c *fiber.Ctx) error {
	var shopMenu model.ShopMenu
	if err := c.BodyParser(&shopMenu); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	if err := db.Create(&shopMenu).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create shop menu",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(shopMenu)
}
func CreateMenuByEntrepreneur(db *gorm.DB, c *fiber.Ctx) error {
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
	var menu model.ShopMenu
	if err := c.BodyParser(&menu); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to parse request body",
			"details": err.Error(),
		})
	}

	// Ensure IsPublic is set to false and assign Shop ID
	menu.IsPublic = false
	menu.ShopID = shop.ID

	// Save the shop menu
	if err := db.Create(&menu).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create shop menu",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(menu)
}

// GetShopMenu retrieves a ShopMenu entry by ID
func GetShopMenu(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var shopMenu model.ShopMenu
	// if err := db.First(&shopMenu, id).Error; err != nil {
	// 	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
	// 		"error": "Shop menu not found",
	// 	})
	// }
	db.First(&shopMenu, id)
	return c.JSON(shopMenu)
}

// GetShopMenuByShopID retrieves ShopMenu entries by Shop ID
func GetShopMenuByShopID(db *gorm.DB, c *fiber.Ctx) error {
	shopID := c.Params("shop_id")
	var shopMenus []map[string]interface{}

	// Query specific fields
	if err := db.Model(&model.ShopMenu{}).
		Select("id, product_description, price, product_name, photo").
		Where("shop_id = ?", shopID).
		Find(&shopMenus).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve shop menus",
		})
	}

	return c.JSON(shopMenus)
}

// UpdateShopMenu updates a ShopMenu entry by ID
func UpdateShopMenu(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var shopMenu model.ShopMenu
	if err := db.First(&shopMenu, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Shop menu not found",
		})
	}

	if err := c.BodyParser(&shopMenu); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if err := db.Save(&shopMenu).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update shop menu",
		})
	}
	return c.JSON(shopMenu)
}

// DeleteShopMenu deletes a ShopMenu entry by ID
func DeleteShopMenu(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	if err := db.Delete(&model.ShopMenu{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete shop menu",
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}