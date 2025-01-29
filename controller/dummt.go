package controller

import (
	"strconv"

	"github.com/HealthMe-pls/medic-go-api/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func getShopOpenDates(db *gorm.DB, shopID uint) ([]fiber.Map, error) {
	var shopOpenDates []model.ShopOpenDate
	if err := db.Where("shop_id = ?", shopID).Find(&shopOpenDates).Error; err != nil {
		return nil, err
	}

	var result []fiber.Map
	for _, date := range shopOpenDates {
		result = append(result, fiber.Map{
			"id":         date.ID,
			"start_time": date.StartTime,
			"end_time":   date.EndTime,
		})
	}
	return result, nil
}

// Helper function to fetch shop menus
func getShopMenus(db *gorm.DB, shopID uint) ([]fiber.Map, error) {
	var shopMenus []model.ShopMenu
	if err := db.Where("shop_id = ?", shopID).Find(&shopMenus).Error; err != nil {
		return nil, err
	}

	var result []fiber.Map
	for _, menu := range shopMenus {
		// Fetch all photos related to the menu by MenuID
		menuPhotos, err := getPhotosByMenuID(db, menu.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, fiber.Map{
			"id":                  menu.ID,
			"product_name":        menu.ProductName,
			"product_description": menu.ProductDescription,
			"price":               menu.Price,
			"photos":              menuPhotos, // Include all photos related to the menu
		})
	}
	return result, nil
}

// Helper function to fetch social media
func getSocialMedia(db *gorm.DB, shopID uint) ([]fiber.Map, error) {
	var socialMedias []model.SocialMedia
	if err := db.Where("shop_id = ?", shopID).Find(&socialMedias).Error; err != nil {
		return nil, err
	}

	var result []fiber.Map
	for _, social := range socialMedias {
		result = append(result, fiber.Map{
			"id":       social.ID,
			"platform": social.Platform,
			"link":     social.Link,
		})
	}
	return result, nil
}

func getPhotosByShopID(db *gorm.DB, shopID uint) ([]fiber.Map, error) {
	var photos []model.Photo
	if err := db.Where("shop_id = ?", shopID).Find(&photos).Error; err != nil {
		return nil, err
	}

	var result []fiber.Map
	for _, photo := range photos {
		result = append(result, fiber.Map{
			"photo_id": photo.ID,
			"pathfile": photo.PathFile,
		})
	}
	return result, nil
}

func getPhotosByMenuID(db *gorm.DB, menuID uint) ([]fiber.Map, error) {
	var photos []model.Photo
	if err := db.Where("menu_id = ?", menuID).Find(&photos).Error; err != nil {
		return nil, err
	}

	var result []fiber.Map
	for _, photo := range photos {
		result = append(result, fiber.Map{
			"photo_id": photo.ID,
			"pathfile": photo.PathFile,
		})
	}
	return result, nil
}

func GetShopsByCategory(db *gorm.DB, c *fiber.Ctx) error {
	// Get the ShopCategoryID parameter from the URL
	shopCategoryID := c.Params("shop_category_id")

	var shops []model.Shop

	// Query shops by the ShopCategoryID
	if err := db.Where("shop_category_id = ?", shopCategoryID).Find(&shops).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("No shops found for this category")
	}

	// Return the shops as a JSON response
	return c.JSON(shops)
}

// -----social media
// CreateSocialMedia creates a new SocialMedia entry
func CreateSocialMedia(db *gorm.DB, c *fiber.Ctx) error {
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

// shop menu
// CreateShopMenu creates a new ShopMenu entry
func CreateShopMenu(db *gorm.DB, c *fiber.Ctx) error {
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

//-----shop category

func CreateShopCategory(db *gorm.DB, c *fiber.Ctx) error {
	// Parse the request body into the ShopCategory struct
	shopCategory := new(model.ShopCategory)
	if err := c.BodyParser(shopCategory); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	// Create the ShopCategory record in the database
	if err := db.Create(shopCategory).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create ShopCategory",
		})
	}

	// Return the newly created ShopCategory as JSON response
	return c.Status(fiber.StatusCreated).JSON(shopCategory)
}

func GetShopCategories(db *gorm.DB, c *fiber.Ctx) error {
	var categories []model.ShopCategory
	if err := db.Find(&categories).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve shop categories",
		})
	}
	return c.JSON(categories)
}
func GetShopCategoryByID(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var category model.ShopCategory
	// if err := db.First(&category, id).Error; err != nil {
	// 	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
	// 		"error": "Shop category not found",
	// 	})
	// }
	db.First(&category, id)
	return c.JSON(category)
}

func UpdateShopCategory(db *gorm.DB, c *fiber.Ctx) error {
	// Get the shop category ID from the URL parameters
	shopCategoryID := c.Params("id")

	// Convert the ID to uint
	id, err := strconv.Atoi(shopCategoryID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ShopCategory ID format")
	}

	// Fetch the existing shop category
	var shopCategory model.ShopCategory
	if err := db.First(&shopCategory, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("ShopCategory not found")
	}

	// Parse the request body for updates
	if err := c.BodyParser(&shopCategory); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to parse request body")
	}

	// Save the updated shop category to the database
	if err := db.Save(&shopCategory).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update ShopCategory")
	}

	// Return the updated shop category as JSON
	return c.Status(fiber.StatusOK).JSON(shopCategory)
}

func DeleteShopCategory(db *gorm.DB, c *fiber.Ctx) error {
	// Get the shop category ID from the URL parameters
	shopCategoryID := c.Params("id")

	// Convert the ID to uint
	id, err := strconv.Atoi(shopCategoryID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ShopCategory ID format")
	}

	// Attempt to delete the shop category
	if err := db.Delete(&model.ShopCategory{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete ShopCategory")
	}

	// Return success message
	return c.SendString("ShopCategory deleted successfully")
}
