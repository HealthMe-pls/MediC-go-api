package controller

import (
	"strconv"
	"fmt"
	"github.com/HealthMe-pls/medic-go-api/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)


func GetAvailableMenus(db *gorm.DB, c *fiber.Ctx) error {
	var availableMenus []model.ShopMenu

	// Fetch only menus where is_public = true
	if err := db.Where("is_public = ?", true).Find(&availableMenus).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve available menus",
		})
	}

	return c.JSON(availableMenus)
}

func GetAvailableMenusByShopID(db *gorm.DB, c *fiber.Ctx) error {
	shopID := c.Params("shop_id")
	var availableMenus []model.ShopMenu

	// Fetch menus where shop_id matches and is_public is true
	if err := db.Where("shop_id = ? AND is_public = ?", shopID, true).Find(&availableMenus).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve available menus for the shop",
		})
	}

	return c.JSON(availableMenus)
}
func GetAvailableMenusDetailByshopID(db *gorm.DB, c *fiber.Ctx) error {
	shopID, err := strconv.ParseUint(c.Params("shop_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid shop ID",
		})
	}

	menus, err := GetAvailableMenusHelper(db, uint(shopID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve available menus",
		})
	}

	return c.JSON(menus)
}

func GetAvailableMenusHelper(db *gorm.DB, shopID uint) ([]fiber.Map, error) {
	var shopMenus []model.ShopMenu
	if err := db.Where("shop_id = ? AND is_public = ?", shopID, true).Find(&shopMenus).Error; err != nil {
		return nil, err
	}

	var result []fiber.Map
	for _, menu := range shopMenus {
		// Fetch all available photos related to the menu
		menuPhotos, err := getAvailablePhotosByMenuID(db, menu.ID)
		if err != nil {
			return nil, err
		}

		result = append(result, fiber.Map{
			"id":                  menu.ID,
			"product_name":        menu.ProductName,
			"product_description": menu.ProductDescription,
			"price":               menu.Price,
			"is_public":           menu.IsPublic,
			"photos":              menuPhotos, // Include only available photos
		})
	}
	return result, nil
}

func getAvailablePhotosByMenuID(db *gorm.DB, menuID uint) ([]model.Photo, error) {
	var photos []model.Photo
	if err := db.Where("is_public = ? AND menu_id = ?", true, menuID).Find(&photos).Error; err != nil {
		return nil, err
	}
	return photos, nil
}

func GetAvailablePhotosByMenuID(db *gorm.DB, c *fiber.Ctx) error {
	menuID := c.Params("menu_id")

	var photos []model.Photo
	if err := db.Where("is_public = ? AND menu_id = ?", true, menuID).Find(&photos).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve available photos for menu",
		})
	}

	return c.JSON(photos)
}

func GetAvailablePhotosByShopID(db *gorm.DB, c *fiber.Ctx) error {
	shopID := c.Params("shop_id")

	var photos []model.Photo
	if err := db.Where("is_public = ? AND shop_id = ?", true, shopID).Find(&photos).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve available photos for shop",
		})
	}

	return c.JSON(photos)
}

func GetAvailableSocialByShopID(db *gorm.DB, c *fiber.Ctx) error {
	shopID := c.Params("shop_id")

	// Define a slice to store the filtered social media records
	var socialMedias []model.SocialMedia

	// Query the database for social media records where shop_id matches and is_public is true
	err := db.Where("shop_id = ? AND is_public = ?", shopID, true).
		Find(&socialMedias).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch social media",
			"details": err.Error(),
		})
	}

	// Return the filtered social media records as JSON
	return c.JSON(socialMedias)
}


//shop
func GetAvailableShopDetailByID(db *gorm.DB, c *fiber.Ctx) error {
	shopID := c.Params("shop_id") // shop_id is still a string
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
	if err := db.Preload("Entrepreneur").Preload("ShopCategory").First(&shop, "id = ?", uint(shopIDUint)).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve shop",
			"details": err.Error(),
		})
	}

	// Fetch only available data using helper functions
	shopOpenDates, err := getShopOpenDates(db, uint(shopIDUint)) // Fetch shop open dates
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve shop open dates",
			"details": err.Error(),
		})
	}

	availableMenus, err := GetAvailableMenusHelper(db, uint(shopIDUint)) // Fetch only public shop menus
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve available shop menus",
			"details": err.Error(),
		})
	}

	availableSocialMedia, err := GetAvailableSocialHelper(db, uint(shopIDUint)) // Fetch only public social media links
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve available social media",
			"details": err.Error(),
		})
	}

	availablePhotos, err := getAvailablePhotosByShopID(db, uint(shopIDUint)) // Fetch only public photos
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve available shop photos",
			"details": err.Error(),
		})
	}

	// Construct the response with only available data
	shopResponse := fiber.Map{
		"shop_id":         shop.ID,
		"name":            shop.Name,
		"entrepreneur_id": shop.Entrepreneur.ID,
		"entrepreneur":    shop.Entrepreneur.Title + " " + shop.Entrepreneur.FirstName + " " + shop.Entrepreneur.MiddleName + " " + shop.Entrepreneur.LastName,
		"category_id":     shop.ShopCategory.ID,
		"category":        shop.ShopCategory.Name,
		"open_status":     shop.OpenStatus,
		"description":     shop.Description,
		"photos":          availablePhotos,  // Include only available photos
		"shop_open_dates": shopOpenDates,
		"menus":           availableMenus,   // Include only public menus
		"social_media":    availableSocialMedia, // Include only public social media
	}

	return c.JSON(shopResponse)
}

func GetAvailableSocialHelper(db *gorm.DB, shopID uint) ([]model.SocialMedia, error) {
	var socialMedias []model.SocialMedia
	err := db.Where("shop_id = ? AND is_public = ?", shopID, true).
		Find(&socialMedias).Error
	if err != nil {
		return nil, err
	}
	return socialMedias, nil
}

func getAvailablePhotosByShopID(db *gorm.DB, shopID uint) ([]model.Photo, error) {
	var photos []model.Photo
	if err := db.Where("is_public = ? AND shop_id = ?", true, shopID).Find(&photos).Error; err != nil {
		return nil, err
	}
	return photos, nil
}
