package controller

import (
	"github.com/HealthMe-pls/medic-go-api/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// HandleApprove ใช้สำหรับ approve TempShop และเรียกใช้งานการ approve เวลา
func HandleApprove(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var tempShop model.TempShop
	if err := db.First(&tempShop, "id = ?", id).Error; err != nil {
		// ถ้าไม่พบ TempShop ให้ส่ง StatusNotFound
		return c.Status(fiber.StatusNotFound).SendString("TempShop not found")
	}
	tempID := tempShop.TempID

	// ส่ง fiber context ไปพร้อมกับ tempID ให้ HandleTimeApprove
	if err := HandleTimeApprove(db, c, tempID); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error handling time approve")
	}

	return c.JSON(tempShop)
}

// HandleTimeApprove ค้นหา TempShopOpenDate ทั้งหมดตาม tempID แล้วประมวลผลตาม Operation ที่ระบุ
func HandleTimeApprove(db *gorm.DB, c *fiber.Ctx, tempID uint) error {
	var tempShopOpenDates []model.TempShopOpenDate
	if err := db.Where("temp_id = ?", tempID).Find(&tempShopOpenDates).Error; err != nil {
		return err
	}

	for _, openDate := range tempShopOpenDates {
		if openDate.Operation == "add" {
			// บันทึก openDate ลงใน Context เพื่อส่งต่อให้ CreateShopOpenDate ใช้
			c.Locals("openDate", openDate)
			if err := addToShopOpenDate(db, c); err != nil {
				return err
			}
		}
	}
	return nil
}

// addToShopOpenDate เรียกใช้งาน CreateShopOpenDate โดยส่ง fiber context ไปด้วย
func addToShopOpenDate(db *gorm.DB, c *fiber.Ctx) error {
	return CreateShopOpenDate(db, c)
}


