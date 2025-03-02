package controller

import (
	"strconv"
	"fmt"
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
