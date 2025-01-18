package controller

import (
	"github.com/HealthMe-pls/medic-go-api/model"
    "github.com/HealthMe-pls/medic-go-api/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
    "golang.org/x/crypto/bcrypt"
)

func GetEntrepreneur(db *gorm.DB,c *fiber.Ctx) error {
	var entrepreneur []model.Entrepreneur
	db.Find(&entrepreneur)
	return c.JSON(entrepreneur)
}

func GetEntrepreneurByID(db *gorm.DB, c *fiber.Ctx) error {
    id := c.Params("id")
    var entrepreneur model.Entrepreneur
    
    // Query the database for the entrepreneur with the provided id
    if err := db.First(&entrepreneur, "id = ?", id).Error; err != nil {
        // If an error occurs (e.g., no entrepreneur found), return a 404
        return c.Status(fiber.StatusNotFound).SendString("Entrepreneur not found")
    }

    // If successful, return the entrepreneur data as a JSON response
    return c.JSON(entrepreneur)
}
func CreateEntrepreneur(db *gorm.DB, c *fiber.Ctx) error {
    // Parse the request body into the Entrepreneur struct
    entrepreneur := new(model.Entrepreneur)
    if err := c.BodyParser(entrepreneur); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Failed to parse request body",
        })
    }

    // Check if the username already exists
    var existingEntrepreneur model.Entrepreneur
    err := db.Where("username = ?", entrepreneur.Username).First(&existingEntrepreneur).Error
    if err == nil {
        // Username already exists
        return c.Status(fiber.StatusConflict).JSON(fiber.Map{
            "error": "Username already exists",
        })
    } else if err != nil && err != gorm.ErrRecordNotFound {
        // Database error occurred while checking
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to check username",
        })
    }

    // Hash the password
    hashedPassword, err := utils.HashPassword(entrepreneur.Password)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to hash password",
        })
    }
    entrepreneur.Password = hashedPassword

    // Save the Entrepreneur to the database
    if result := db.Create(&entrepreneur); result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to create entrepreneur",
        })
    }

    // Return the created Entrepreneur as a JSON response
    return c.Status(fiber.StatusCreated).JSON(entrepreneur)
}
//login
func LoginEntrepreneur(db *gorm.DB, c *fiber.Ctx) error {
	// Parse the login request
	loginRequest := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	// Find the entrepreneur by username
	var entrepreneur model.Entrepreneur
	if err := db.Where("username = ?", loginRequest.Username).First(&entrepreneur).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}

	// Verify the password
	if err := bcrypt.CompareHashAndPassword([]byte(entrepreneur.Password), []byte(loginRequest.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}

	// Login successful, return a response
	// Optionally, generate and return a JWT token
	return c.JSON(fiber.Map{
		"message": "Login successful",
		"entrepreneur": fiber.Map{
			"id":          entrepreneur.ID,
			"username":    entrepreneur.Username,
			"firstName":   entrepreneur.FirstName,
			"lastName":    entrepreneur.LastName,
			"phoneNumber": entrepreneur.PhoneNumber,
		},
	})
}
// func CreateEntrepreneur(db *gorm.DB, c *fiber.Ctx) error {
// 	// Parse the request body into the Entrepreneur struct
// 	entrepreneur := new(model.Entrepreneur)
// 	if err := c.BodyParser(entrepreneur); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Failed to parse request body",
// 		})
// 	}

// 	// Check if the username already exists
//     var existingEntrepreneur model.Entrepreneur
// 	err := db.Where("username = ?", entrepreneur.Username).First(&existingEntrepreneur).Error
// 	if err == nil {
// 		// Username already exists
// 		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
// 			"error": "Username already exists",
// 		})
// 	} else if err != nil && err != gorm.ErrRecordNotFound {
// 		// Database error occurred while checking
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Failed to check username",
// 		})
// 	}

// 	// Save the Entrepreneur to the database
// 	if result := db.Create(&entrepreneur); result.Error != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": "Failed to create entrepreneur",
// 		})
// 	}

// 	// Return the created Entrepreneur as a JSON response
// 	return c.Status(fiber.StatusCreated).JSON(entrepreneur)
// }


func UpdateEntrepreneur(db *gorm.DB, c *fiber.Ctx) error {
    // Get the id parameter from the URL
    id := c.Params("id")
    var entrepreneur model.Entrepreneur

    // Find the entrepreneur by id
    if err := db.First(&entrepreneur, "id = ?", id).Error; err != nil {
        return c.Status(fiber.StatusNotFound).SendString("Entrepreneur not found")
    }

    // Parse the updated details from the request body
    if err := c.BodyParser(&entrepreneur); err != nil {
        return c.Status(fiber.StatusBadRequest).SendString("Failed to parse request body")
    }

    // Save the updated entrepreneur details to the database
    if result := db.Save(&entrepreneur); result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to update entrepreneur")
    }

    // Return the updated entrepreneur as a JSON response
    return c.JSON(entrepreneur)
}

// func DeleteEntrepreneur(db *gorm.DB, c *fiber.Ctx) error {
//     // Get the id parameter from the URL
//     id := c.Params("id")

//     // Delete the entrepreneur from the database by their id
//     if result := db.Delete(&model.Entrepreneur{}, "id = ?", id); result.Error != nil {
//         return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete entrepreneur")
//     }

//     // Return success message
//     return c.SendString("Entrepreneur successfully deleted")
// }

func DeleteEntrepreneurAndShops(db *gorm.DB, c *fiber.Ctx) error {
    // Get the entrepreneur's id from the URL parameter
    id := c.Params("id")

    // Delete all shops associated with the entrepreneur
    if err := db.Where("entrepreneur_username = ?", id).Delete(&model.Shop{}).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete shops")
    }

    // Now delete the entrepreneur
    if err := db.Where("id = ?", id).Delete(&model.Entrepreneur{}).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete entrepreneur")
    }

    return c.SendString("Entrepreneur and associated shops successfully deleted")
}

