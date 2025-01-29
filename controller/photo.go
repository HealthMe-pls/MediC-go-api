package controller

import (
	"github.com/HealthMe-pls/medic-go-api/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"fmt"
	"strconv"
)


// CreatePhoto creates a new Photo entry
func CreatePhoto(db *gorm.DB, c *fiber.Ctx) error {
	var photo model.Photo
	if err := c.BodyParser(&photo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	if err := db.Create(&photo).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create photo",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(photo)
}

// GetPhoto retrieves a Photo entry by ID
func GetPhoto(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var photo model.Photo
	if err := db.First(&photo, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Photo not found",
		})
	}
	return c.JSON(photo)
}

// GetPhotoByMenuID retrieves Photo entries by Menu ID
func GetPhotoByMenuID(db *gorm.DB, c *fiber.Ctx) error {
	menuID := c.Params("menu_id")
	var photos []model.Photo
	if err := db.Where("menu_id = ?", menuID).Find(&photos).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve photos by menu ID",
		})
	}
	return c.JSON(photos)
}

// GetPhotoByShopID retrieves Photo entries by Shop ID
func GetPhotoByShopID(db *gorm.DB, c *fiber.Ctx) error {
	shopID := c.Params("shop_id")
	var photos []model.Photo
	if err := db.Where("shop_id = ?", shopID).Find(&photos).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve photos by shop ID",
		})
	}
	return c.JSON(photos)
}

// UpdatePhoto updates a Photo entry by ID
func UpdatePhoto(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var photo model.Photo
	if err := db.First(&photo, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Photo not found",
		})
	}

	if err := c.BodyParser(&photo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if err := db.Save(&photo).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update photo",
		})
	}
	return c.JSON(photo)
}

// DeletePhoto deletes a Photo entry by ID
func DeletePhoto(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	if err := db.Delete(&model.Photo{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete photo",
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}















func CreatePhotoByShopID(db *gorm.DB, c *fiber.Ctx) error {
	// Extract Shop ID from request params
	shopID, err := strconv.ParseUint(c.Params("shop_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid shop ID",
		})
	}

	// Read file from request
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read image",
		})
	}

	// Save the file to the server
	filePath := fmt.Sprintf("./uploads/%s", file.Filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save image",
		})
	}

	// Fetch TempID from TempShop where ShopID matches
	var tempShop model.TempShop
	if err := db.Where("shop_id = ?", shopID).First(&tempShop).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "TempShop not found for given ShopID",
		})
	}

	// Create new photo record
	photo := model.Photo{
		PathFile: file.Filename,
		ShopID:   uintPtr(uint(shopID)),
		TempID:   &tempShop.TempID, // Get TempID from TempShop
		IsPublic: false,            // Always false
	}

	// Save to database
	if err := db.Create(&photo).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create photo",
		})
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(photo)
}

// Helper function to convert uint to *uint
func uintPtr(u uint) *uint {
	return &u
}
