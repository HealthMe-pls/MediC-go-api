package controller

import (
	"fmt"
	"os"
	"github.com/HealthMe-pls/medic-go-api/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)


// ==================== MENU BIN ====================

// Get all deleted menus
func GetBinMenus(db *gorm.DB, c *fiber.Ctx) error {
	var menuBin []model.DeleteMenu
	if err := db.Find(&menuBin).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve deleted menus",
		})
	}
	return c.JSON(menuBin)
}

// Get a deleted menu by ID
func GetBinMenuByID(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var menu model.DeleteMenu
	if err := db.First(&menu, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Deleted menu not found",
		})
	}
	return c.JSON(menu)
}

// Get all deleted menus by TempID
func GetBinMenuByTempID(db *gorm.DB, c *fiber.Ctx) error {
	tempID := c.Params("temp_id")
	var menuBin []model.DeleteMenu
	if err := db.Where("temp_id = ?", tempID).Find(&menuBin).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve deleted menus",
		})
	}
	return c.JSON(menuBin)
}

// Add a deleted menu entry
func CreateBinMenu(db *gorm.DB, c *fiber.Ctx) error {
	var deleteMenu model.DeleteMenu
	if err := c.BodyParser(&deleteMenu); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := db.Create(&deleteMenu).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to store deleted menu",
		})
	}

	return c.JSON(deleteMenu)
}

// Delete a deleted menu entry by ID
func DeleteBinMenu(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	if err := db.Delete(&model.DeleteMenu{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete menu bin entry",
		})
	}
	return c.JSON(fiber.Map{"message": "Deleted menu bin entry removed"})
}

func DeleteBinMenuByTempID(db *gorm.DB, c *fiber.Ctx) error {
	tempID := c.Params("id")

	// Begin a transaction
	tx := db.Begin()

	// Fetch all delete menu entries for the given temp_id
	var deleteMenus []model.DeleteMenu
	if err := tx.Where("temp_id = ?", tempID).Find(&deleteMenus).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch delete menus",
		})
	}

	// Delete each menu using DeleteShopMenu
	for _, deleteMenu := range deleteMenus {
		if err := DeleteShopMenu(tx, deleteMenu.MenuID); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to delete menu ID %d", deleteMenu.MenuID),
			})
		}
	}

	// Delete the DeleteMenu entries after all ShopMenus are deleted
	if err := tx.Where("temp_id = ?", tempID).Delete(&model.DeleteMenu{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete menus by TempID",
		})
	}

	// Commit the transaction after all operations succeed
	tx.Commit()

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Deleted all menus for TempID"})
}




// ==================== PHOTO BIN ====================

func GetBinPhotos(db *gorm.DB, c *fiber.Ctx) error {
	var photoBin []model.DeletePhoto
	if err := db.Find(&photoBin).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve deleted photos",
		})
	}
	return c.JSON(photoBin)
}

func GetBinPhotoByID(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var photo model.DeletePhoto
	if err := db.First(&photo, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Deleted photo not found",
		})
	}
	return c.JSON(photo)
}

func GetBinPhotoByTempID(db *gorm.DB, c *fiber.Ctx) error {
	tempID := c.Params("temp_id")
	var photoBin []model.DeletePhoto
	if err := db.Where("temp_id = ?", tempID).Find(&photoBin).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve deleted photos",
		})
	}
	return c.JSON(photoBin)
}

func CreateBinPhoto(db *gorm.DB, c *fiber.Ctx) error {
	var deletePhoto model.DeletePhoto
	if err := c.BodyParser(&deletePhoto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := db.Create(&deletePhoto).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to store deleted photo",
		})
	}

	return c.JSON(deletePhoto)
}

func DeleteBinPhoto(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	if err := db.Delete(&model.DeletePhoto{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete photo bin entry",
		})
	}
	return c.JSON(fiber.Map{"message": "Deleted photo bin entry removed"})
}

// func DeleteBinPhotoByTempID(db *gorm.DB, c *fiber.Ctx) error {
// 	tempID := c.Params("temp_id")
// 	if err := db.Where("temp_id = ?", tempID).Delete(&model.DeletePhoto{}).Error; err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete photos by TempID"})
// 	}
// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Deleted all photos for TempID"})
// }
// DeleteBinPhotoByTempID deletes photos based on TempID from both the database and the uploads folder
func DeleteBinPhotoByTempID(db *gorm.DB, c *fiber.Ctx) error {
	tempID := c.Params("id")

	// Fetch all deletePhoto records matching the tempID
	var deletePhotos []model.DeletePhoto
	if err := db.Where("temp_id = ?", tempID).Find(&deletePhotos).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch delete photo records"})
	}

	// Extract photo IDs from DeletePhoto table
	var photoIDs []uint
	for _, deletePhoto := range deletePhotos {
		photoIDs = append(photoIDs, deletePhoto.PhotoID)
	}

	if len(photoIDs) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No photos found for given TempID"})
	}

	// Fetch photos from the Photo table based on extracted photo IDs
	var photos []model.Photo
	if err := db.Where("id IN ?", photoIDs).Find(&photos).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch photos from Photo table"})
	}

	// Delete files from the uploads folder
	for _, photo := range photos {
		filePath := fmt.Sprintf("./uploads/%s", photo.PathFile)
		if err := os.Remove(filePath); err != nil {
			fmt.Println("Error deleting file:", err) // Log error but continue deletion process
		}
	}

	// Delete photo records from the Photo table
	if err := db.Where("id IN ?", photoIDs).Delete(&model.Photo{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete photos from database"})
	}

	// Optionally delete entries from DeletePhoto table as well
	if err := db.Where("temp_id = ?", tempID).Delete(&model.DeletePhoto{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete delete-photo records"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Deleted all photos for TempID successfully"})
}

// ==================== SOCIAL BIN ====================

func GetBinSocials(db *gorm.DB, c *fiber.Ctx) error {
	var socialBin []model.DeleteSocial
	if err := db.Find(&socialBin).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve deleted social records",
		})
	}
	return c.JSON(socialBin)
}

func GetBinSocialByID(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var social model.DeleteSocial
	if err := db.First(&social, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Deleted social record not found",
		})
	}
	return c.JSON(social)
}

func GetBinSocialByTempID(db *gorm.DB, c *fiber.Ctx) error {
	tempID := c.Params("temp_id")
	var socialBin []model.DeleteSocial
	if err := db.Where("temp_id = ?", tempID).Find(&socialBin).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve deleted social records",
		})
	}
	return c.JSON(socialBin)
}

func CreateBinSocial(db *gorm.DB, c *fiber.Ctx) error {
	var deleteSocial model.DeleteSocial
	if err := c.BodyParser(&deleteSocial); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := db.Create(&deleteSocial).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to store deleted social record",
		})
	}

	return c.JSON(deleteSocial)
}

func DeleteBinSocial(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	if err := db.Delete(&model.DeleteSocial{}, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete social bin entry",
		})
	}
	return c.JSON(fiber.Map{"message": "Deleted social bin entry removed"})
}

// func DeleteBinSocialByTempID(db *gorm.DB, c *fiber.Ctx) error {
// 	tempID := c.Params("temp_id")
// 	if err := db.Where("temp_id = ?", tempID).Delete(&model.DeleteSocial{}).Error; err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete social data by TempID"})
// 	}
// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Deleted all social data for TempID"})
// }

func DeleteBinSocialByTempID(db *gorm.DB, c *fiber.Ctx) error {
	tempID := c.Params("id")

	// Fetch all the delete social entries for the given temp_id
	var deleteSocials []model.DeleteSocial
	if err := db.Where("temp_id = ?", tempID).Find(&deleteSocials).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch delete socials"})
	}

	// Extract the social media IDs to delete from SocialMedia table
	socialIDs := make([]uint, len(deleteSocials))
	for i, social := range deleteSocials {
		socialIDs[i] = social.SocialID
	}

	// Delete SocialMedia entries where IDs match the extracted social IDs
	if err := db.Where("id IN ?", socialIDs).Delete(&model.SocialMedia{}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete social media entries"})
	}

	// Delete the DeleteSocial entries
	// if err := db.Where("temp_id = ?", tempID).Delete(&model.DeleteSocial{}).Error; err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete social data by TempID"})
	// }

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Deleted all social data for TempID"})
}
