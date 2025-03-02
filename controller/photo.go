package controller

import (
	"fmt"
	"strconv"
	"github.com/HealthMe-pls/medic-go-api/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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
// CreatePhotoByMenuID creates a Photo by MenuID with IsPublic set to false
func CreatePhotoByMenuID(db *gorm.DB, c *fiber.Ctx, isPublic bool) error {
	// Parse menu_id from the URL params
	menuID, err := strconv.Atoi(c.Params("menu_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid menu ID",
		})
	}
	unitMenuID := uint(menuID)

	// Query the Menu to get the tempid
	var menu model.ShopMenu
	if err := db.First(&menu, unitMenuID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Menu not found",
		})
	}

	// Parse request body into the Photo struct
	// Read file from request
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to read image")
	}
	
	// Save the file to the server
	filePath := fmt.Sprintf("./uploads/%s", file.Filename)
	fileName := file.Filename
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to save image")
	}

	photo := new(model.Photo)
	if err := c.BodyParser(photo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to parse request body",
			"details": err.Error(),
		})
	}
	photo.PathFile = fileName
	// Set IsPublic to false and assign MenuID
	photo.IsPublic = isPublic
	photo.MenuID = &unitMenuID // Convert menuID to uint
	// Set tempid from Menu to Photo
	photo.TempID = menu.TempID // Assuming TempID is the field in both models

	// Create the Photo entry in the database
	if err := db.Create(&photo).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create photo",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"photo": photo,
	})
}

// // CreatePhotoByMenuID creates a Photo by MenuID with IsPublic set to false
// func CreatePhotoByMenuID(db *gorm.DB, c *fiber.Ctx, isPublic bool) error {
// 	// Parse menu_id from the URL params
// 	menuID, err := strconv.Atoi(c.Params("menu_id"))
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Invalid menu ID",
// 		})
// 	}
// 	unitMenuID := uint(menuID)
// 	// Parse request body into the Photo struct
// 	// Read file from request
// 	file, err := c.FormFile("image")
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).SendString("Failed to read image")
// 	}
	
// 	// Save the file to the server
// 	filePath := fmt.Sprintf("./uploads/%s", file.Filename)
// 	fileName := file.Filename
// 	if err := c.SaveFile(file, filePath); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).SendString("Failed to save image")
// 	}
// 	photo := new(model.Photo)
// 	if err := c.BodyParser(photo); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error":   "Failed to parse request body",
// 			"details": err.Error(),
// 		})
// 	}
// 	photo.PathFile = fileName
// 	// Set IsPublic to false and assign MenuID
// 	photo.IsPublic = isPublic
// 	photo.MenuID = &unitMenuID // Convert menuID to uint

// 	// Create the Photo entry in the database
// 	if err := db.Create(&photo).Error; err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error":   "Failed to create photo",
// 			"details": err.Error(),
// 		})
// 	}

// 	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
// 		"photo": photo,
// 	})
// }
// CreatePhotoByShopID creates a Photo by ShopID with IsPublic set to false
func CreatePhotoByShopID(db *gorm.DB, c *fiber.Ctx, isPublic bool) error {
	// Parse shop_id from the URL params
	shopID, err := strconv.Atoi(c.Params("shop_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid shop ID",
		})
	}
	unitShopId := uint(shopID)

	// Query the TempShop entry by ShopID to get the TempID
	var tempShop model.TempShop
	if err := db.Where("shop_id = ?", unitShopId).First(&tempShop).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "TempShop not found for the given ShopID",
		})
	}

	// Read file from request
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to read image")
	}
	
	// Save the file to the server
	filePath := fmt.Sprintf("./uploads/%s", file.Filename)
	fileName := file.Filename
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to save image")
	}

	// Parse request body into the Photo struct
	photo := new(model.Photo)
	if err := c.BodyParser(photo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to parse request body",
			"details": err.Error(),
		})
	}

	// Set IsPublic to false and assign ShopID
	photo.PathFile = fileName
	photo.IsPublic = isPublic
	photo.ShopID = &unitShopId // Convert shopID to uint
	// Set TempID from TempShop to Photo
	photo.TempID = &tempShop.TempID // Assuming TempID is the field in both models

	// Create the Photo entry in the database
	if err := db.Create(&photo).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create photo",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"photo": photo,
	})
}


// CreatePhotoByWorkshopID creates a Photo by WorkshopID with IsPublic set to false
func CreatePhotoByWorkshopID(db *gorm.DB, c *fiber.Ctx) error {
	workshopID, err := strconv.Atoi(c.Params("workshop_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid workshop ID",
		})
	}

	// Convert workshopID from int to uint
	uintWorkshopID := uint(workshopID)
	// Read file from request
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to read image")
	}
	
	// Save the file to the server
	filePath := fmt.Sprintf("./uploads/%s", file.Filename)
	fileName := file.Filename
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to save image")
	}

	// Parse request body into the Photo struct
	photo := new(model.Photo)
	if err := c.BodyParser(photo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to parse request body",
			"details": err.Error(),
		})
	}

	// Set IsPublic to true
	photo.PathFile = fileName
	photo.IsPublic = true
	photo.WorkshopID = &uintWorkshopID // Assign pointer to uint

	// Create the Photo entry in the database
	if err := db.Create(&photo).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create photo",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"photo": photo,
	})
}

