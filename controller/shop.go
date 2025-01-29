package controller

import (
	"fmt"
	"strconv"

	"github.com/HealthMe-pls/medic-go-api/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// func Shop(db *gorm.DB, c *fiber.Ctx) error {
// 	var entrepreneur []model.Entrepreneur
// 	db.Find(&entrepreneur)
// 	return c.JSON(entrepreneur)
// }
func GetShopByID(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var shop model.Shop
	if err := db.First(&shop, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Shop not found",
			"details": err.Error(),
		})
	}
	// db.First(&shop, id)
	return c.JSON(shop)
}
func GetShops(db *gorm.DB, c *fiber.Ctx) error {
	var shops []model.Shop
	if err := db.Find(&shops).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve shop categories",
			"details": err.Error(),
		})
	}
	return c.JSON(shops)
}

func GetShopDetail(db *gorm.DB, c *fiber.Ctx) error {
	var shops []model.Shop

	// Fetch basic shop details with Entrepreneur and ShopCategory preloaded
	if err := db.Preload("Entrepreneur").Preload("ShopCategory").Find(&shops).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve shops",
			"details": err.Error(),
		})
	}

	// Construct the detailed response
	var shopResponses []fiber.Map
	for _, shop := range shops {
		shopID := shop.ID

		// Fetch shop open dates
		shopOpenDates, err := getShopOpenDates(db, shopID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to retrieve shop open dates",
				"details": err.Error(),
			})
		}

		// Fetch shop menus
		shopMenus, err := getShopMenus(db, shopID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to retrieve shop menus",
				"details": err.Error(),
			})
		}

		// Fetch social media
		socialMedias, err := getSocialMedia(db, shopID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to retrieve social media entries",
				"details": err.Error(),
			})
		}

		// Fetch all photos related to the shop by shopID
		shopPhotos, err := getPhotosByShopID(db, shopID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to retrieve photos",
				"details": err.Error(),
			})
		}

		// Construct the shop response
		shopResponses = append(shopResponses, fiber.Map{
			"shop_id":         shop.ID,
			"name":            shop.Name,
			"entrepreneur_id": shop.Entrepreneur.ID,
			"entrepreneur":    shop.Entrepreneur.Title + " " + shop.Entrepreneur.FirstName + " " + shop.Entrepreneur.MiddleName + " " + shop.Entrepreneur.LastName,
			"category_id":     shop.ShopCategory.ID,
			"category":        shop.ShopCategory.Name,
			"open_status":     shop.OpenStatus,
			"description":     shop.Description,
			"photos":          shopPhotos, // Updated to include all photos related to the shop
			"shop_open_dates": shopOpenDates,
			"menus":           shopMenus,
			"social_media":    socialMedias,
		})
	}

	return c.JSON(shopResponses)
}

func stringToUint(shopID string) (uint, error) {
	// Log shopID to verify it
	fmt.Println("Received shopID:", shopID)

	// Try to convert the string to uint
	id, err := strconv.ParseUint(shopID, 10, 32)
	if err != nil {
		// Log error for debugging
		fmt.Println("Error parsing shopID:", err)
		return 0, err
	}

	return uint(id), nil
}

func GetShopDetailByID(db *gorm.DB, c *fiber.Ctx) error {
	shopID := c.Params("id") // shop_id is still a string
	fmt.Println("shopID from URL parameter:", shopID)

	// Convert the string shopID to uint
	shopIDUint, err := stringToUint(shopID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid shop ID",
		})
	}

	// Fetch a single shop by ID with Entrepreneur and ShopCategory preloaded
	var shop model.Shop
	if err := db.Preload("Entrepreneur").Preload("ShopCategory").First(&shop, "id = ?", shopIDUint).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve shop",
			"details": err.Error(),
		})
	}

	// Fetch related data using helper functions
	shopOpenDates, err := getShopOpenDates(db, shopIDUint) // Fetch shop open dates
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve shop open dates",
			"details": err.Error(),
		})
	}

	shopMenus, err := getShopMenus(db, shopIDUint) // Fetch shop menus
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve shop menus",
			"details": err.Error(),
		})
	}

	socialMedias, err := getSocialMedia(db, shopIDUint) // Fetch social media links
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve social media entries",
			"details": err.Error(),
		})
	}

	// Fetch all photos related to the shop using shopID
	shopPhotos, err := getPhotosByShopID(db, shopIDUint)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve shop photos",
			"details": err.Error(),
		})
	}

	// Construct the shop response
	shopResponse := fiber.Map{
		"shop_id":         shop.ID,
		"name":            shop.Name,
		"entrepreneur_id": shop.Entrepreneur.ID,
		"entrepreneur":    shop.Entrepreneur.Title + " " + shop.Entrepreneur.FirstName + " " + shop.Entrepreneur.MiddleName + " " + shop.Entrepreneur.LastName,
		"category_id":     shop.ShopCategory.ID,
		"category":        shop.ShopCategory.Name,
		"open_status":     shop.OpenStatus,
		"description":     shop.Description,
		"photos":          shopPhotos, // Include all photos related to the shop
		"shop_open_dates": shopOpenDates,
		"menus":           shopMenus,
		"social_media":    socialMedias,
	}

	return c.JSON(shopResponse)
}


func CreateShop(db *gorm.DB, c *fiber.Ctx) error {
	// Parse the request body into the Shop struct
	shop := new(model.Shop)
	if err := c.BodyParser(shop); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to parse request body",
			"details": err.Error(),
		})
	}

	// Check if the ShopCategory exists by its ID
	var shopCategory model.ShopCategory
	if err := db.First(&shopCategory, shop.ShopCategoryID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "ShopCategory not found",
			"details": err.Error(),
		})
	}

	// Save the Shop to the database
	if result := db.Create(&shop); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create shop",
		})
	}

	// Return the created Shop as a JSON response
	return c.Status(fiber.StatusCreated).JSON(shop)
}
func UpdateShop(db *gorm.DB, c *fiber.Ctx) error {
	// Get the shop ID parameter from the URL
	id := c.Params("id")
	var shop model.Shop

	// Find the shop by ID
	if err := db.First(&shop, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Shop not found",
			"details": err.Error(),
		})
		// return c.Status(fiber.StatusNotFound).SendString("Shop not found")
	}

	// Parse the updated details from the request body
	if err := c.BodyParser(&shop); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to parse request body")
	}

	// Save the updated shop details to the database
	if result := db.Save(&shop); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update shop")
	}

	// Return the updated shop as a JSON response
	return c.JSON(shop)
}
func DeleteShop(db *gorm.DB, c *fiber.Ctx) error {
	// Get the shop ID parameter from the URL
	id := c.Params("id")

	// Delete the shop from the database by its ID
	if result := db.Delete(&model.Shop{}, id); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete shop")
	}

	// Return success message
	return c.SendString("Shop successfully deleted")
}

