package controller

import (
	"github.com/HealthMe-pls/medic-go-api/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// --------------------- TempShop Controller --------------------- //

// GetTempShops ดึงข้อมูล TempShop ทั้งหมด
func GetTempShops(db *gorm.DB, c *fiber.Ctx) error {
	var tempShops []model.TempShop
	if err := db.Find(&tempShops).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve temp shops",
			"details": err.Error(),
		})
	}
	return c.JSON(tempShops)
}

// GetTempShopByID ดึงข้อมูล TempShop ตาม ID
func GetTempShopByID(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var tempShop model.TempShop
	if err := db.First(&tempShop, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Temp shop not found",
			"details": err.Error(),
		})
	}
	return c.JSON(tempShop)
}

// CreateTempShop สร้าง TempShop ใหม่
func CreateTempShop(db *gorm.DB, c *fiber.Ctx) error {
	tempShop := new(model.TempShop)
	if err := c.BodyParser(tempShop); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to parse request body",
			"details": err.Error(),
		})
	}
	if result := db.Create(&tempShop); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create temp shop",
			"details": result.Error.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(tempShop)
}

// UpdateTempShop อัปเดต TempShop ตาม ID
func UpdateTempShop(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var tempShop model.TempShop
	if err := db.First(&tempShop, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Temp shop not found",
			"details": err.Error(),
		})
	}
	if err := c.BodyParser(&tempShop); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to parse request body",
			"details": err.Error(),
		})
	}
	if result := db.Save(&tempShop); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update temp shop",
			"details": result.Error.Error(),
		})
	}
	return c.JSON(tempShop)
}

// DeleteTempShop ลบ TempShop ตาม ID
func DeleteTempShop(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	if result := db.Delete(&model.TempShop{}, id); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to delete temp shop",
			"details": result.Error.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// --------------------- TempShopOpenDate Controller --------------------- //

// GetTempShopOpenDates ดึงข้อมูล TempShopOpenDate ทั้งหมด
func GetTempShopOpenDates(db *gorm.DB, c *fiber.Ctx) error {
	var tempShopOpenDates []model.TempShopOpenDate
	if err := db.Find(&tempShopOpenDates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve temp shop open dates",
			"details": err.Error(),
		})
	}
	return c.JSON(tempShopOpenDates)
}

// GetTempShopOpenDateByID ดึง TempShopOpenDate ตาม ID
func GetTempShopOpenDateByID(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var tempShopOpenDate model.TempShopOpenDate
	if err := db.First(&tempShopOpenDate, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Temp shop open date not found",
			"details": err.Error(),
		})
	}
	return c.JSON(tempShopOpenDate)
}

// CreateTempShopOpenDate สร้าง TempShopOpenDate ใหม่
func CreateTempShopOpenDate(db *gorm.DB, c *fiber.Ctx) error {
	tempShopOpenDate := new(model.TempShopOpenDate)
	if err := c.BodyParser(tempShopOpenDate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to parse request body",
			"details": err.Error(),
		})
	}
	if result := db.Create(&tempShopOpenDate); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create temp shop open date",
			"details": result.Error.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(tempShopOpenDate)
}

// UpdateTempShopOpenDate อัปเดต TempShopOpenDate ตาม ID
func UpdateTempShopOpenDate(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var tempShopOpenDate model.TempShopOpenDate
	if err := db.First(&tempShopOpenDate, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Temp shop open date not found",
			"details": err.Error(),
		})
	}
	if err := c.BodyParser(&tempShopOpenDate); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to parse request body",
			"details": err.Error(),
		})
	}
	if result := db.Save(&tempShopOpenDate); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update temp shop open date",
			"details": result.Error.Error(),
		})
	}
	return c.JSON(tempShopOpenDate)
}

// DeleteTempShopOpenDate ลบ TempShopOpenDate ตาม ID
func DeleteTempShopOpenDate(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	if result := db.Delete(&model.TempShopOpenDate{}, id); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to delete temp shop open date",
			"details": result.Error.Error(),
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
