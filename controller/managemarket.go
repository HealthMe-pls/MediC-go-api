package controller

import (
	"github.com/HealthMe-pls/medic-go-api/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CreateMarketOpenDate creates a new MarketOpenDate entry
func CreateMarketOpenDate(db *gorm.DB, c *fiber.Ctx) error {
	var marketOpenDate model.MarketOpenDate
	if err := c.BodyParser(&marketOpenDate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	if err := db.Create(&marketOpenDate).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create market open date",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(marketOpenDate)
}
// GetAllMarketDates retrieves all MarketOpenDate entries
func GetAllMarketDates(db *gorm.DB, c *fiber.Ctx) error {
	var marketOpenDates []model.MarketOpenDate

	// Fetch all market open dates with related shop open dates
	if err := db.Preload("ShopOpenDates").Find(&marketOpenDates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve market open dates",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"market_open_dates": marketOpenDates,
	})
}

// GetMarketOpenDate retrieves a MarketOpenDate by ID
func GetMarketOpenDate(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var marketOpenDate model.MarketOpenDate
	if err := db.First(&marketOpenDate, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Market open date not found",
		})
	}
	return c.JSON(marketOpenDate)
}

// UpdateMarketOpenDate updates an existing MarketOpenDate by ID
func UpdateMarketOpenDate(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var marketOpenDate model.MarketOpenDate
	if err := db.First(&marketOpenDate, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Market open date not found",
		})
	}

	if err := c.BodyParser(&marketOpenDate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if err := db.Save(&marketOpenDate).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update market open date",
		})
	}
	return c.JSON(marketOpenDate)
}

// DeleteMarketOpenDate deletes a MarketOpenDate by ID
func DeleteMarketOpenDate(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	if err := db.Delete(&model.MarketOpenDate{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete market open date",
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

//shop open time
// CreateShopOpenDate creates a new ShopOpenDate entry
func CreateShopOpenDate(db *gorm.DB, c *fiber.Ctx) error {
	var shopOpenDate model.ShopOpenDate
	if err := c.BodyParser(&shopOpenDate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	if err := db.Create(&shopOpenDate).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create shop open date",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(shopOpenDate)
}
// GetAllShopTimes retrieves all shop open dates
func GetAllShopTimes(db *gorm.DB, c *fiber.Ctx) error {
	var shopOpenDates []model.ShopOpenDate

	// Fetch all shop open times with related shop and market open date
	if err := db.Preload("Shop").Preload("MarketOpenDate").Find(&shopOpenDates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve shop open times",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"shop_open_dates": shopOpenDates,
	})
}
// GetShopOpenDate retrieves a ShopOpenDate entry by ID
func GetShopOpenDate(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var shopOpenDate model.ShopOpenDate
	if err := db.First(&shopOpenDate, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Shop open date not found",
		})
	}
	return c.JSON(shopOpenDate)
}

// GetShopOpenDateByShopID retrieves ShopOpenDate entries by Shop ID
func GetShopOpenDateByShopID(db *gorm.DB, c *fiber.Ctx) error {
	shopID := c.Params("shop_id")
	var shopOpenDates []map[string]interface{}

	// Query the required fields and map the results directly
	if err := db.Model(&model.ShopOpenDate{}).
		Select("id, start_time, end_time").
		Where("shop_id = ?", shopID).
		Find(&shopOpenDates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve shop open dates",
		})
	}

	return c.JSON(shopOpenDates)
}

// UpdateShopOpenDate updates a ShopOpenDate entry by ID
func UpdateShopOpenDate(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var shopOpenDate model.ShopOpenDate
	if err := db.First(&shopOpenDate, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Shop open date not found",
		})
	}

	if err := c.BodyParser(&shopOpenDate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if err := db.Save(&shopOpenDate).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update shop open date",
		})
	}
	return c.JSON(shopOpenDate)
}

// DeleteShopOpenDate deletes a ShopOpenDate entry by ID
func DeleteShopOpenDate(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	if err := db.Delete(&model.ShopOpenDate{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete shop open date",
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

//contrat to admin
// CreateContactToAdmin creates a new ContactToAdmin entry
func CreateContactToAdmin(db *gorm.DB, c *fiber.Ctx) error {
	var contactToAdmin model.ContactToAdmin
	if err := c.BodyParser(&contactToAdmin); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	if err := db.Create(&contactToAdmin).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create contact",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(contactToAdmin)
}


func GetAllContacts(db *gorm.DB, c *fiber.Ctx) error {
	var contacts []model.ContactToAdmin

	// Fetch all contact requests from the database
	if err := db.Find(&contacts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve contacts",
		})
	}

	return c.JSON(contacts)
}

// GetContactToAdmin retrieves a ContactToAdmin entry by ID
func GetContactToAdmin(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var contactToAdmin model.ContactToAdmin
	if err := db.First(&contactToAdmin, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Contact not found",
		})
	}
	return c.JSON(contactToAdmin)
}

// UpdateContactToAdmin updates a ContactToAdmin entry by ID
func UpdateContactToAdmin(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var contactToAdmin model.ContactToAdmin
	if err := db.First(&contactToAdmin, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Contact not found",
		})
	}

	if err := c.BodyParser(&contactToAdmin); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if err := db.Save(&contactToAdmin).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update contact",
		})
	}
	return c.JSON(contactToAdmin)
}

// DeleteContactToAdmin deletes a ContactToAdmin entry by ID
func DeleteContactToAdmin(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	if err := db.Delete(&model.ContactToAdmin{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete contact",
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// GetAllTempShopsWaiting retrieves all TempShop entries with status "waiting"
func GetAllTempShopsWaiting(db *gorm.DB, c *fiber.Ctx) error {
	var tempShops []model.TempShop

	// Fetch all TempShop entries where status is "Waiting"
	if err := db.Where("status = ?", "Waiting").
		Find(&tempShops).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve temp shops",
			"details": err.Error(),
		})
	}

	var tempShopResponses []fiber.Map
	for _, tempShop := range tempShops {
		if tempShop.ShopID == nil {
			continue // Skip if ShopID is nil
		}
		shopID := *tempShop.ShopID

		// Fetch social media info
		TempSocials, err := GetTempSocialsByShopID(db, shopID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to retrieve social media data",
				"details": err.Error(),
			})
		}

		// Fetch menu info
		tempMenus, err := getMenuForTempByShopID(db, shopID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to retrieve menu data",
				"details": err.Error(),
			})
		}

		// Fetch shop photos and menu photos separately
		photoData, err := getPhotoForTempByShopID(db, shopID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to retrieve photo data",
				"details": err.Error(),
			})
		}

		// Extract separate lists for shop photos and menu photos
		photosShop, _ := photoData["photos_shop"].([]fiber.Map)
		photosMenu, _ := photoData["photos_menu"].([]fiber.Map)

		// Fetch shop open time info
		tempTimes, err := GetTempTimeByShopID(db, shopID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to retrieve shop time data",
				"details": err.Error(),
			})
		}

		// Fetch permanent shop open date info
		shopOpenDates, err := GetShopOpenDateForTempByShopID(db, shopID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to retrieve shop open date data",
				"details": err.Error(),
			})
		}

		socials, _ := TempSocials["socials"].([]fiber.Map)

		// Extract time data
		addTime, _ := tempTimes["addTime"].([]fiber.Map)
		editTime, _ := tempTimes["editTime"].([]fiber.Map)
		deleteTime, _ := tempTimes["deleteTime"].([]fiber.Map)

		tempShopResponses = append(tempShopResponses, fiber.Map{
			"id":            tempShop.TempID,
			"name":          tempShop.Name,
			"description":	 tempShop.Description,
			"category_id":	 tempShop.ShopCategoryID,
			"shop_id":       shopID,
			"socials":       socials,
			"menus":         tempMenus,  // Include menus in the response
			"photos_shop":   photosShop, // Photos directly related to the shop
			"photos_menu":   photosMenu, // Photos linked to menus
			"addTime":       addTime,    // Include added times
			"editTime":      editTime,   // Include edited times
			"deleteTime":    deleteTime, // Include deleted times
			"time":          shopOpenDates, // Include shop open dates
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"temp_shops": tempShopResponses,
	})
}




func GetTempSocialsByShopID(db *gorm.DB, shopID uint) (fiber.Map, error) {
    var socials []model.SocialMedia

    // Subquery to get SocialIDs that exist in DeleteSocial
    subQuery := db.Table("delete_socials").Select("social_id")

    // Query for socials that are not in DeleteSocial
    if err := db.Where("shop_id = ? AND id NOT IN (?)", shopID, subQuery).
        Find(&socials).Error; err != nil {
        return nil, err
    }

    // Convert results to fiber.Map
    socialsResult := make([]fiber.Map, len(socials))
    for i, social := range socials {
        socialsResult[i] = fiber.Map{
            "id":        social.ID,
            "name":      social.Name,
            "platform":  social.Platform,
            "link":      social.Link,
            "shop_id":   social.ShopID,
            "is_public": social.IsPublic,
        }
    }

    return fiber.Map{
        "socials": socialsResult,
    }, nil
}
func getMenuForTempByShopID(db *gorm.DB, shopID uint) ([]fiber.Map, error) {
    var menus []model.ShopMenu

    // Subquery to get MenuIDs that exist in DeleteMenu
    subQuery := db.Table("delete_menus").Select("menu_id")

    // Query menus based on the conditions
    if err := db.Where("shop_id = ? AND id NOT IN (?)", shopID, subQuery).
        Find(&menus).Error; err != nil {
        return nil, err
    }

    // Convert results to fiber.Map
    var menuResults []fiber.Map

    for _, menu := range menus {
        // Get associated photos
        photos, err := getPhotoForTempByMenuID(db, menu.ID)
        if err != nil {
            return nil, err
        }

        menuResults = append(menuResults, fiber.Map{
            "id":                  menu.ID,
            "product_name":        menu.ProductName,
            "product_description": menu.ProductDescription,
            "price":               menu.Price,
            "shop_id":             menu.ShopID,
            "is_public":           menu.IsPublic,
            "photos":              photos,
        })
    }

    return menuResults, nil
}


func getPhotoForTempByMenuID(db *gorm.DB, menuID uint) ([]fiber.Map, error) {
	var photos []model.Photo

	// Subquery to get PhotoIDs that exist in DeletePhoto
	subQuery := db.Table("delete_photos").Select("photo_id")

	// Query photos based on the condition
	if err := db.Where("(is_public = ? AND id NOT IN (?)) OR is_public = ?", true, subQuery, false).
		Where("menu_id = ?", menuID).
		Find(&photos).Error; err != nil {
		return nil, err
	}

	// Convert results to fiber.Map
	result := make([]fiber.Map, len(photos))
	for i, photo := range photos {
		result[i] = fiber.Map{
			"id":        photo.ID,
			"path_file": photo.PathFile,
			"menu_id":   photo.MenuID,
			"is_public": photo.IsPublic,
		}
	}

	return result, nil
}

func getPhotoForTempByShopID(db *gorm.DB, shopID uint) (fiber.Map, error) {
	var photosShop []model.Photo
	var photosMenu []model.Photo

	// Subquery to get PhotoIDs that exist in DeletePhoto
	subQuery := db.Table("delete_photos").Select("photo_id")

	// Fetch photos that belong to the shop (excluding deleted ones)
	if err := db.Where("(is_public = ? AND id NOT IN (?)) OR is_public = ?", true, subQuery, false).
		Where("shop_id = ?", shopID).
		Find(&photosShop).Error; err != nil {
		return nil, err
	}

	// Get menu IDs that belong to the given shop
	var menuIDs []uint
	if err := db.Table("shop_menus").Where("shop_id = ?", shopID).Pluck("id", &menuIDs).Error; err != nil {
		return nil, err
	}

	// Fetch photos that belong to those menu IDs (excluding deleted ones)
	if len(menuIDs) > 0 {
		if err := db.Where("(is_public = ? AND id NOT IN (?)) OR is_public = ?", true, subQuery, false).
			Where("menu_id IN (?)", menuIDs).
			Find(&photosMenu).Error; err != nil {
			return nil, err
		}
	}

	// Convert results to fiber.Map
	convertPhotos := func(photos []model.Photo) []fiber.Map {
		result := make([]fiber.Map, len(photos))
		for i, photo := range photos {
			result[i] = fiber.Map{
				"id":        photo.ID,
				"path_file": photo.PathFile,
				"shop_id":   photo.ShopID,
				"menu_id":   photo.MenuID,
				"is_public": photo.IsPublic,
			}
		}
		return result
	}

	// Return structured response
	return fiber.Map{
		"photos_shop": convertPhotos(photosShop), // Photos with shop_id
		"photos_menu": convertPhotos(photosMenu), // Photos linked to menus
	}, nil
}


func GetTempTimeByShopID(db *gorm.DB, shopID uint) (fiber.Map, error) {
	var addTime []model.TempShopOpenDate
	var editTime []model.TempShopOpenDate
	var deleteTime []model.TempShopOpenDate

	// Query for added times
	if err := db.Where("shop_id = ? AND operation = ?", shopID, "add").Find(&addTime).Error; err != nil {
		return nil, err
	}

	// Query for edited times
	if err := db.Where("shop_id = ? AND operation = ?", shopID, "edit").Find(&editTime).Error; err != nil {
		return nil, err
	}

	// Query for deleted times
	if err := db.Where("shop_id = ? AND operation = ?", shopID, "delete").Find(&deleteTime).Error; err != nil {
		return nil, err
	}

	// Convert results to fiber.Map
	addTimeResult := make([]fiber.Map, len(addTime))
	for i, timeEntry := range addTime {
		addTimeResult[i] = fiber.Map{
			"id":                 timeEntry.ID,
			"start_time":         timeEntry.StartTime,
			"end_time":           timeEntry.EndTime,
			"shop_id":            timeEntry.ShopID,
			"market_open_date_id": timeEntry.MarketOpenDateID,
		}
	}

	editTimeResult := make([]fiber.Map, len(editTime))
	for i, timeEntry := range editTime {
		editTimeResult[i] = fiber.Map{
			"id":                 timeEntry.ID,
			"start_time":         timeEntry.StartTime,
			"end_time":           timeEntry.EndTime,
			"shop_id":            timeEntry.ShopID,
			"market_open_date_id": timeEntry.MarketOpenDateID,
		}
	}

	deleteTimeResult := make([]fiber.Map, len(deleteTime))
	for i, timeEntry := range deleteTime {
		deleteTimeResult[i] = fiber.Map{
			"id":                 timeEntry.ID,
			"start_time":         timeEntry.StartTime,
			"end_time":           timeEntry.EndTime,
			"shop_id":            timeEntry.ShopID,
			"market_open_date_id": timeEntry.MarketOpenDateID,
		}
	}

	return fiber.Map{
		"addTime":    addTimeResult,
		"editTime":   editTimeResult,
		"deleteTime": deleteTimeResult,
	}, nil
}

// GetShopOpenDateByShopID retrieves all ShopOpenDate entries for a given shop ID
func GetShopOpenDateForTempByShopID(db *gorm.DB, shopID uint) ([]fiber.Map, error) {
	var shopOpenDates []model.ShopOpenDate

	// Query for shop open dates
	if err := db.Where("shop_id = ?", shopID).Find(&shopOpenDates).Error; err != nil {
		return nil, err
	}

	// Convert result to fiber.Map
	shopOpenDateResults := make([]fiber.Map, len(shopOpenDates))
	for i, openDate := range shopOpenDates {
		shopOpenDateResults[i] = fiber.Map{
			"id":                 openDate.ID,
			"start_time":         openDate.StartTime,
			"end_time":           openDate.EndTime,
			"shop_id":            openDate.ShopID,
			"market_open_date_id": openDate.MarketOpenDateID,
		}
	}

	return shopOpenDateResults, nil
}
