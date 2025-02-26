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
