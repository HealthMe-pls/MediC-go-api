package controller

import (
	"github.com/HealthMe-pls/medic-go-api/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Handleapprove(db *gorm.DB,c *fiber.Ctx) error {
	id := c.Params("id")
	var tempShop model.TempShop
	if err := db.First(&tempShop, "id = ?", id).Error; err != nil {
        // If an error occurs (e.g., no entrepreneur found), return a 404
        return c.Status(fiber.StatusNotFound).SendString("TempShop not found")
    }
	tempID := tempShop.TempID
	Handletimeapprove(db,tempID)
	return c.JSON(tempShop)
}

func Handletimeapprove(db *gorm.DB,tempID uint) error {
	var tempShopOpenDate []model.TempShopOpenDate
	if err := db.Where("temp_id = ?", tempID).Find(&tempShopOpenDate).Error; err != nil {
		return nil
	}

	for _, time := range tempShopOpenDate {
		if time.Operation == "add" {
			addToShopOpenDate(db,time)
		}else if time.Operation == "delete" {
			deleteToShopOpenDate(db,time)
		}else if time.Operation == "edit" {
			editToShopOpenDate(db,time)
		}
	}
	return nil
}

func addToShopOpenDate(db *gorm.DB,time model.TempShopOpenDate) error {
	
}

func deleteToShopOpenDate(db *gorm.DB,time model.TempShopOpenDate) error {
	
}

func editToShopOpenDate(db *gorm.DB,time model.TempShopOpenDate) error {
	
}