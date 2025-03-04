package controller

import (
	"fmt"
	"os"
	"strconv"

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
// UpdateShopMenu updates a ShopMenu entry and its corresponding TempMenu entry
func UpdateShopMenu(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")

	// Find the existing ShopMenu
	var shopMenu model.ShopMenu
	if err := db.First(&shopMenu, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Shop menu not found",
		})
	}

	// Parse the request body
	if err := c.BodyParser(&shopMenu); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Begin transaction to update both tables
	tx := db.Begin()

	// Update the ShopMenu entry
	if err := tx.Save(&shopMenu).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update shop menu",
		})
	}

	// Check if a TempMenu exists for this menu ID
	var tempMenu model.TempMenu
	if err := tx.Where("menu_id = ?", shopMenu.ID).First(&tempMenu).Error; err == nil {
		// Update TempMenu if it exists
		tempMenu.ProductDescription = shopMenu.ProductDescription
		tempMenu.Price = shopMenu.Price
		tempMenu.ProductName = shopMenu.ProductName

		if err := tx.Save(&tempMenu).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update temp menu",
			})
		}
	}

	// Commit the transaction
	tx.Commit()

	return c.JSON(shopMenu)
}
// func DeleteShopMenu(tx *gorm.DB, c *fiber.Ctx) error {
// 	id := c.Params("id")

// 	// Step 1: Check if there are photos linked to the menu
// 	var photos []model.Photo
// 	if err := tx.Where("menu_id = ?", id).Find(&photos).Error; err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Failed to check associated photos",
// 		})
// 	}

// 	// Step 2: Delete associated photos
// 	for _, photo := range photos {
// 		filePath := fmt.Sprintf("./uploads/%s", photo.PathFile)

// 		// Try deleting the file, log an error if it fails but continue
// 		if err := os.Remove(filePath); err != nil {
// 			fmt.Println("Error deleting file:", err)
// 		}

// 		// Delete photo from DB
// 		if err := tx.Delete(&photo).Error; err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 				"error": "Failed to delete associated photo from database",
// 			})
// 		}
// 	}

// 	// Step 3: Delete the corresponding TempMenu entry
// 	if err := tx.Where("menu_id = ?", id).Delete(&model.TempMenu{}).Error; err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Failed to delete associated TempMenu",
// 		})
// 	}

// 	// Step 4: Delete the ShopMenu entry
// 	if err := tx.Where("id = ?", id).Delete(&model.ShopMenu{}).Error; err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Failed to delete shop menu",
// 		})
// 	}

// 	// No need to commit or rollback here (handled in the main function)
// 	return nil
// }

// DeleteShopMenu deletes a menu and related data within the same transaction
func DeleteShopMenu(tx *gorm.DB, menuID uint) error {
	// Step 1: Retrieve and delete associated photos
	var photos []model.Photo
	if err := tx.Where("menu_id = ?", menuID).Find(&photos).Error; err != nil {
		return fmt.Errorf("failed to check associated photos")
	}

	// Delete photos from filesystem
	for _, photo := range photos {
		filePath := fmt.Sprintf("./uploads/%s", photo.PathFile)
		if err := os.Remove(filePath); err != nil {
			fmt.Println("Error deleting file:", err)
		}
	}

	// Delete all photos in DB
	if err := tx.Where("menu_id = ?", menuID).Delete(&model.Photo{}).Error; err != nil {
		return fmt.Errorf("failed to delete photos from database")
	}

	// Step 2: Delete the corresponding TempMenu entry
	if err := tx.Where("menu_id = ?", menuID).Delete(&model.TempMenu{}).Error; err != nil {
		return fmt.Errorf("failed to delete associated TempMenu")
	}

	// Step 3: Delete the ShopMenu entry
	if err := tx.Where("id = ?", menuID).Delete(&model.ShopMenu{}).Error; err != nil {
		return fmt.Errorf("failed to delete shop menu")
	}

	return nil
}

// DeleteShopMenuByID extracts menu_id from request, converts to uint, and calls DeleteShopMenu
func DeleteShopMenuByID(db *gorm.DB, c *fiber.Ctx) error {
	// Get menu ID from request parameters
	menuID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid menu ID",
		})
	}

	// Convert to uint
	uintMenuID := uint(menuID)

	// Begin transaction
	tx := db.Begin()

	// Call DeleteShopMenu within the transaction
	if err := DeleteShopMenu(tx, uintMenuID); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Commit transaction if all operations succeed
	tx.Commit()

	// Return success response
	return c.SendString("Shop menu successfully deleted")
}


func CreateMenuWithTemp(db *gorm.DB, c *fiber.Ctx, isPublic bool) error {
	menu := new(model.ShopMenu)
	if err := c.BodyParser(menu); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to parse request body",
			"details": err.Error(),
		})
	}

	// Assign isPublic from function parameter
	menu.IsPublic = isPublic
    // Find TempShop that has the same ShopID as the social media
    var tempShop model.TempShop
    if err := db.Where("shop_id = ?", menu.ShopID).First(&tempShop).Error; err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error":   "TempShop not found for this ShopID",
            "details": err.Error(),
        })
    }
	menu.TempID = &tempShop.TempID
	// Save the menu in ShopMenu table
	if result := db.Create(&menu); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create menu",
			"details": result.Error.Error(),
		})
	}


	// Create a corresponding TempMenu entry
	tempMenu := model.TempMenu{
		TempID:             tempShop.TempID,
		MenuID:             menu.ID,
		ProductName:        menu.ProductName,
		ProductDescription: menu.ProductDescription,
		Price:              menu.Price,
	}

	if result := db.Create(&tempMenu); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create temp menu",
			"details": result.Error.Error(),
		})
	}

	// Return created menu and temp menu
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"menu":     menu,
		"tempMenu": tempMenu,
	})
}



func CreateTempMenu(db *gorm.DB, c *fiber.Ctx) error {
	tempMenu := new(model.TempMenu)
	if err := c.BodyParser(tempMenu); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to parse request body",
			"details": err.Error(),
		})
	}

	if result := db.Create(&tempMenu); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create temp menu",
			"details": result.Error.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(tempMenu)
}


func GetShopIDByMenuID(db *gorm.DB, c *fiber.Ctx) error {
	menuID := c.Params("menu_id")

	// Check if the menu exists
	var shopMenu model.ShopMenu
	if err := db.First(&shopMenu, "id = ?", menuID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Menu not found",
			"details": err.Error(),
		})
	}

	// Return the ShopID associated with the Menu
	return c.JSON(fiber.Map{
		"shop_id": shopMenu.ShopID,
	})
}

func UpdateTempMenuByMenuID(db *gorm.DB, c *fiber.Ctx) error {
	menuID, err := strconv.Atoi(c.Params("menu_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid menu ID",
		})
	}

	// Fetch the ShopMenu
	var menu model.ShopMenu
	if err := db.First(&menu, menuID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "ShopMenu not found",
		})
	}

	// Check if TempMenu exists for this MenuID
	var tempMenu model.TempMenu
	if err := db.Where("menu_id = ?", menu.ID).First(&tempMenu).Error; err != nil {
		// If not found, create a new TempMenu entry
		tempMenu = model.TempMenu{
			MenuID:             menu.ID,
			TempID:             *menu.TempID, // Ensure TempID is not nil
			ProductDescription: menu.ProductDescription,
			Price:              menu.Price,
			ProductName:        menu.ProductName,
		}
		if err := db.Create(&tempMenu).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create temp menu",
			})
		}
	} else {
		// If found, update existing TempMenu
		tempMenu.ProductDescription = menu.ProductDescription
		tempMenu.Price = menu.Price
		tempMenu.ProductName = menu.ProductName

		if err := db.Save(&tempMenu).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update temp menu",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Temp menu updated successfully",
		"temp_menu": tempMenu,
	})
}


func UpdateMenuFromTemp(db *gorm.DB, tempID uint) error {
	// Fetch the TempMenu record by TempID
	var tempMenu model.TempMenu
	if err := db.First(&tempMenu, "temp_id = ?", tempID).Error; err != nil {
		return fmt.Errorf("TempMenu not found: %w", err)
	}

	// Fetch the corresponding ShopMenu record by MenuID
	var shopMenu model.ShopMenu
	if err := db.First(&shopMenu, "id = ?", tempMenu.MenuID).Error; err != nil {
		return fmt.Errorf("ShopMenu not found for MenuID %d: %w", tempMenu.MenuID, err)
	}

	// Update ShopMenu with values from TempMenu
	shopMenu.ProductDescription = tempMenu.ProductDescription
	shopMenu.Price = tempMenu.Price
	shopMenu.ProductName = tempMenu.ProductName

	// Save the updated ShopMenu record
	if err := db.Save(&shopMenu).Error; err != nil {
		return fmt.Errorf("failed to update ShopMenu: %w", err)
	}

	return nil
}
