package controller

import (
	"fmt"

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

func addToShopOpenDate(db *gorm.DB, time model.TempShopOpenDate) error {
	shopOpenDate := model.ShopOpenDate{
		StartTime:        time.StartTime,
		EndTime:          time.EndTime,
		ShopID:           time.ShopID,
		MarketOpenDateID: time.MarketOpenDateID,
	}

	if err := db.Create(&shopOpenDate).Error; err != nil {
		return fmt.Errorf("failed to create shop open date: %w", err)
	}

	return nil
}

func deleteToShopOpenDate(db *gorm.DB,time model.TempShopOpenDate) error {
		var record model.ShopOpenDate
		if err := db.Where("shop_id = ? AND market_open_date_id = ?", time.ShopID, time.MarketOpenDateID).First(&record).Error; err != nil {
			return err
		}
		if err := db.Delete(&record).Error; err != nil {
			return err
		}
		return nil
}

func editToShopOpenDate(db *gorm.DB,time model.TempShopOpenDate) error {
	 var record model.ShopOpenDate
	 if err := db.Where("shop_id = ? AND market_open_date_id = ?", time.ShopID, time.MarketOpenDateID).First(&record).Error; err != nil {
		 return err
	 }
	 record.StartTime = time.StartTime
	 record.EndTime = time.EndTime
	 if err := db.Save(&record).Error; err != nil {
		 return err
	 }
	 return nil
}