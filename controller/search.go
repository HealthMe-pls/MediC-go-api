package controller

import (
	"strings"

	"github.com/HealthMe-pls/medic-go-api/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)
func SearchShopsByKeyword(db *gorm.DB, c *fiber.Ctx) error {
	// Extract keyword from query parameters
	keyword := c.Query("keyword", "")
	if keyword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Keyword is required",
		})
	}

	// Use lowercased keyword for case-insensitive filtering
	lowerKeyword := "%" + strings.ToLower(keyword) + "%"

	// Define a map to store the results
	var results []fiber.Map

	// Step 1: Search for shops matching the keyword in their name
	var shops []model.Shop
	if err := db.Where("LOWER(name) LIKE ?", lowerKeyword).Find(&shops).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to search shops",
			"details": err.Error(),
		})
	}

	// Add matched shops to the results
	for _, shop := range shops {
		results = append(results, fiber.Map{
			"shop_id":   shop.ID,
			"matchWord": shop.Name, // Shop name is the matchWord
		})
	}

	// Step 2: Search for shops matching the keyword in their menus
	var shopMenus []model.ShopMenu
	if err := db.Preload("Shop"). // Preload the associated Shop for ShopMenu
		Where("LOWER(product_name) LIKE ? OR LOWER(product_description) LIKE ?", lowerKeyword, lowerKeyword).
		Find(&shopMenus).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to search shop menus",
			"details": err.Error(),
		})
	}

	// Add matched menus to the results, avoiding duplicates
	seenShops := make(map[uint]bool) // Track shop IDs already added
	for _, result := range results {
		seenShops[result["shop_id"].(uint)] = true
	}

	for _, menu := range shopMenus {
		if !seenShops[menu.ShopID] { // Only add if the shop isn't already in results
			results = append(results, fiber.Map{
				"shop_id":   menu.ShopID,
				"matchWord": menu.ProductName, // Menu name is the matchWord
			})
			seenShops[menu.ShopID] = true
		}
	}

	// Return the combined results as JSON
	return c.JSON(results)
}





