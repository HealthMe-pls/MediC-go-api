package controller

import (
	"fmt"
	"github.com/HealthMe-pls/medic-go-api/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetAllWorkshops retrieves all workshops
func GetWorkshops(db *gorm.DB, c *fiber.Ctx) error {
	var workshops []model.Workshop

	// Fetch all workshops
	if err := db.Find(&workshops).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch workshops",
			"details": err.Error(),
		})
	}

	// Construct the response
	var workshopResponses []fiber.Map
	for _, workshop := range workshops {
		// Fetch photos for each workshop
		photos, err := getPhotosByWorkshopID(db, workshop.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   "Failed to fetch photos for workshop",
				"details": err.Error(),
			})
		}

		// Add workshop details to the response
		workshopResponses = append(workshopResponses, fiber.Map{
			"id":          workshop.ID,
			"name":        workshop.Name,
			"description": workshop.Description,
			"price":       workshop.Price,
			"language":    workshop.Language,
			"instructor":  workshop.Instructor,
			"start_time":  workshop.StartTime,
			"end_time":    workshop.EndTime,
			"date":        workshop.Date,
			"photos":      photos,
		})
	}

	return c.JSON(workshopResponses)
}


// GetWorkshopByID retrieves a workshop by its ID, including photos
func GetWorkshopByID(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")

	// Fetch the workshop by ID
	var workshop model.Workshop
	if err := db.First(&workshop, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Workshop not found",
			"details": err.Error(),
		})
	}

	// Fetch photos for the workshop
	photos, err := getPhotosByWorkshopID(db, workshop.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch photos for workshop",
			"details": err.Error(),
		})
	}

	// Construct the workshop response
	workshopResponse := fiber.Map{
		"id":          workshop.ID,
		"name":        workshop.Name,
		"description": workshop.Description,
		"price":       workshop.Price,
		"language":    workshop.Language,
		"instructor":  workshop.Instructor,
		"start_time":  workshop.StartTime,
		"end_time":    workshop.EndTime,
		"date":        workshop.Date,
		"photos":      photos,
	}

	return c.JSON(workshopResponse)
}



// Helper function to fetch photos by workshop ID
func getPhotosByWorkshopID(db *gorm.DB, workshopID uint) ([]fiber.Map, error) {
	var photos []model.Photo

	// Query the database for photos associated with the given workshop ID
	if err := db.Where("workshop_id = ?", workshopID).Find(&photos).Error; err != nil {
		return nil, fmt.Errorf("could not fetch photos for workshop ID %d: %v", workshopID, err)
	}

	// Transform the photos into a response-friendly format
	var result []fiber.Map
	for _, photo := range photos {
		result = append(result, fiber.Map{
			"photo_id": photo.ID,
			"pathfile": photo.PathFile,
		})
	}

	return result, nil
}


// CreateWorkshop creates a new workshop and returns only its ID
func CreateWorkshop(db *gorm.DB, c *fiber.Ctx) error {
	var workshop model.Workshop
	if err := c.BodyParser(&workshop); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request payload",
			"details": err.Error(),
		})
	}

	if err := db.Create(&workshop).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create workshop",
			"details": err.Error(),
		})
	}

	// Return only the ID of the created workshop
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": workshop.ID,
	})
}

// UpdateWorkshop updates an existing workshop
func UpdateWorkshop(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	var workshop model.Workshop
	if err := db.First(&workshop, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Workshop not found",
			"details": err.Error(),
		})
	}
	if err := c.BodyParser(&workshop); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
			"details": err.Error(),
		})
	}
	if err := db.Save(&workshop).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update workshop",
			"details": err.Error(),
		})
	}
	return c.JSON(workshop)
}

// DeleteWorkshop deletes a workshop by its ID
func DeleteWorkshop(db *gorm.DB, c *fiber.Ctx) error {
	id := c.Params("id")
	if err := db.Delete(&model.Workshop{}, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete workshop",
			"details": err.Error(),
		})
	}
	return c.SendString("Workshop successfully deleted")
}
