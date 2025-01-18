package controller

import (
	"strings"
	"github.com/HealthMe-pls/medic-go-api/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func FilterShopsByKeyword(db *gorm.DB, c *fiber.Ctx) error {
	// Extract keyword from query parameters
	keyword := c.Query("keyword", "")
	if keyword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Keyword is required",
		})
	}

	// Use lowercased keyword for case-insensitive filtering
	lowerKeyword := "%" + strings.ToLower(keyword) + "%"

	// Define a slice to store the filtered shops
	var shops []model.Shop

	// Query the database
	err := db.Preload("ShopMenus").
		Where("LOWER(name) LIKE ? OR LOWER(full_description) LIKE ? OR LOWER(brief_description) LIKE ?", 
			lowerKeyword, lowerKeyword, lowerKeyword).
		Or("EXISTS (SELECT 1 FROM shop_menus WHERE shop_menus.shop_id = shops.id AND LOWER(shop_menus.product_name) LIKE ?)", lowerKeyword).
		Find(&shops).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to filter shops",
		})
	}

	// Return the filtered shops as JSON
	return c.JSON(shops)
}