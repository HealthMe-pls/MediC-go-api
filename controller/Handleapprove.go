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
	if err := db.First(&tempShop, "temp_id = ?", id).Error; err != nil {
        // If an error occurs (e.g., no entrepreneur found), return a 404
        return c.Status(fiber.StatusNotFound).SendString("TempShop not found")
    }
	
	tempID := tempShop.TempID
	// Change the status to "Approve"
	if err := ChangeStateHandleApprove(db, tempID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update status to Approve",
			"details": err.Error(),
		})
	}
	// Update Shop details from TempShop
	if err := UpdateShopFromTemp(db, tempID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update Shop from TempShop",
			"details": err.Error(),
		})
	}
	if err := DeleteBinMenuByTempID(db, c); err != nil {
		return err
	}

	if err := DeleteBinPhotoByTempID(db, c); err != nil {
		return err
	}

	if err := DeleteBinSocialByTempID(db, c); err != nil {
		return err
	}
	// Update IsPublic status to true for all related items by TempID
	if err := UpdateStatusToPublicByTempID(db, tempID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update IsPublic to true",
			"details": err.Error(),
		})
	}
	Handletimeapprove(db,tempID)
	// Fetch updated TempShop
	if err := db.First(&tempShop, "temp_id = ?", tempID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch updated TempShop",
			"details": fmt.Sprintf("Error fetching TempShop after approval for TempID %d: %v", tempID, err),
		})
	}

	// If everything goes well, return the TempShop
	return c.JSON(fiber.Map{
		"tempShop": tempShop,
		"error":    nil,
	})
}
func HandleNotApprove(db *gorm.DB, c *fiber.Ctx) error {
	// Extract the temp ID from the request
	tempID := c.Params("temp_id")

	// Check if tempID is valid
	var tempShop model.TempShop
	if err := db.First(&tempShop, "temp_id = ? AND status = ?", tempID, "Waiting").Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "TempShop not found or not in Waiting state",
		})
	}

	// Update the status to "NotApprove"
	if err := db.Model(&tempShop).Update("status", "NotApprove").Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update status",
			"details": err.Error(),
		})
	}

	// Return success response
	return c.JSON(fiber.Map{
		"message": "TempShop status updated to NotApprove",
		"temp_id": tempShop.TempID,
		"status":  "NotApprove",
	})
}

func ChangeStateHandleApprove(db *gorm.DB, tempID uint) error {
	// Update status to "Approve" where temp_id matches
	if err := db.Model(&model.TempShop{}).
		Where("temp_id = ?", tempID).
		Update("status", "Approve").Error; err != nil {
		return err
	}
	return nil
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