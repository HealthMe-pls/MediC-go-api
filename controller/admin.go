package controller

import (
	"github.com/HealthMe-pls/medic-go-api/model"
	"github.com/HealthMe-pls/medic-go-api/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Get all Admins
func GetAdmins(db *gorm.DB,c *fiber.Ctx) error {	
	var admins []model.Admin
	db.Find(&admins)
	return c.JSON(admins)
}

// Get Admin by Username
func GetAdminByUsername(db *gorm.DB,c *fiber.Ctx) error {
	id := c.Params("id")
	var admin model.Admin
	// if err := db.First(&admin, "id = ?", id).Error; err != nil {
	// 	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
	// 		"error": "Admin not found",
	// 		"details": err.Error(),
	// 	})
	// 	// return c.Status(fiber.StatusNotFound).SendString("Admin not found")
	// }
	db.First(&admin, id)
	return c.JSON(admin)
}


// Create Admin
// func CreateAdmin(db *gorm.DB, c *fiber.Ctx) error {
//     // Parse the request body into the Admin struct
//     admin := new(model.Admin)
//     if err := c.BodyParser(admin); err != nil {
//         return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
//             "error": "Failed to parse request body",
// 			"details": err.Error(),
//         })
//     }

//     // Save the Admin to the database
//     if result := db.Create(&admin); result.Error != nil {
//         return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//             "error": "Failed to create admin",
			
// 		})
//     }


//     // Return the created Admin as a JSON response
//     return c.JSON(admin)
// }
func CreateAdmin(db *gorm.DB, c *fiber.Ctx) error {
    // Parse the request body into the Admin struct
    admin := new(model.Admin)
    if err := c.BodyParser(admin); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error":   "Failed to parse request body",
            "details": err.Error(),
        })
    }

    // Hash the password
    hashedPassword, err := utils.HashPassword(admin.Password)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to hash password",
        })
    }
    admin.Password = hashedPassword

    // Save the Admin to the database
    if result := db.Create(&admin); result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to create admin",
        })
    }

    // Return the created Admin as a JSON response
    return c.JSON(admin)
}
//login
func LoginAdmin(db *gorm.DB, c *fiber.Ctx) error {
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

	// Find the admin by username
	var admin model.Admin
	if err := db.Where("username = ?", loginRequest.Username).First(&admin).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}

	// Verify the password
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(loginRequest.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid username or password",
		})
	}

	// Login successful, return a response
	// Optionally, generate and return a JWT token
	return c.JSON(fiber.Map{
		"message": "Login successful",
		"admin": fiber.Map{
			"id":        admin.ID,
			"username":  admin.Username,
			"firstName": admin.FirstName,
			"lastName":  admin.LastName,
		},
	})
}
// Update Admin

func UpdateAdmin(db *gorm.DB,c *fiber.Ctx) error {
	id := c.Params("id")
	var admin model.Admin
	if err := db.First(&admin, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Admin not found")
	}
	// Parse updated data from request body
	if err := c.BodyParser(&admin); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Failed to parse request body")
	}
	if result := db.Save(&admin); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update admin")
	}
	return c.JSON(admin)
}


// Delete Admin
func DeleteAdmin(db *gorm.DB,c *fiber.Ctx) error {
	id := c.Params("id")
	if result := db.Delete(&model.Admin{}, "id = ?", id); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete admin")
	}
	return c.SendString("Admin successfully deleted")
}
